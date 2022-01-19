package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"

	"github.com/loghole/dbhook"
)

type ctxKey string

const (
	key ctxKey = "test"

	ctxValue = "some value"
)

type beforeHook struct {
	t *testing.T
}

func (h *beforeHook) Before(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	assert.Equal(h.t, ctx.Value(key), ctxValue)

	return ctx, nil
}

func (h *beforeHook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	assert.Equal(h.t, ctx.Value(key), ctxValue)

	return ctx, nil
}

func (h *beforeHook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	assert.Equal(h.t, ctx.Value(key), ctxValue)

	return ctx, nil
}

func TestQuery(t *testing.T) {
	var (
		driverName = "test-query"
		ctrl       = gomock.NewController(t)
		ctx        = context.WithValue(context.Background(), key, ctxValue)
	)

	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	var (
		drv   = MakeDefaultDriver(ctrl, "mock-query")
		hook  = &beforeHook{t: t}
		hooks = dbhook.NewHooks(
			dbhook.WithHooksBefore(hook),
			dbhook.WithHooksAfter(hook),
			dbhook.WithHooksError(hook),
		)
	)

	sql.Register(driverName, dbhook.Wrap(drv, hooks))

	db, err := sql.Open(driverName, "test-query")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, "SELECT some FROM table")
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}

	var (
		id    string
		rows2 = db.QueryRowContext(ctx, "SELECT some FROM table")
	)

	if err := rows2.Scan(&id); err != nil {
		t.Fatal(err)
	}
}

func TestExec(t *testing.T) {
	var (
		driverName = "test-exec"
		ctrl       = gomock.NewController(t)
		ctx        = context.WithValue(context.Background(), key, ctxValue)
	)

	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	var (
		drv   = MakeDefaultDriver(ctrl, "mock-exec")
		hook  = &beforeHook{t: t}
		hooks = dbhook.NewHooks(
			dbhook.WithHooksBefore(hook),
			dbhook.WithHooksAfter(hook),
			dbhook.WithHooksError(hook),
		)
	)

	sql.Register(driverName, dbhook.Wrap(drv, hooks))

	db, err := sql.Open(driverName, "test-exec")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO table(value) VALUES($1)`, "val"); err != nil {
		t.Fatal(err)
	}
}
