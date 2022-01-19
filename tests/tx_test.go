package tests

import (
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"

	"github.com/loghole/dbhook"
)

func TestTx(t *testing.T) {
	var (
		driverName = "test-tx"
		ctrl       = gomock.NewController(t)
		ctx        = context.WithValue(context.Background(), key, ctxValue)
	)

	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	var (
		drv   = MakeTxDriver(ctrl, "mock-tx")
		hook  = &beforeHook{t: t}
		hooks = dbhook.NewHooks(
			dbhook.WithHooksBefore(hook),
			dbhook.WithHooksAfter(hook),
			dbhook.WithHooksError(hook),
		)
	)

	sql.Register(driverName, dbhook.Wrap(drv, hooks))

	tmpDB, err := sql.Open(driverName, "test-tx")
	if err != nil {
		t.Fatal(err)
	}

	var (
		db = sqlx.NewDb(tmpDB, "postgres")
	)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(ctx, `SELECT test FROM test WHERE id=$1`, 12); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

type testID struct {
	ID string `db:"id"`
}

func TestSQLxTx(t *testing.T) {
	var (
		driverName = "test-sqlx-tx"
		ctrl       = gomock.NewController(t)
		ctx        = context.WithValue(context.Background(), key, ctxValue)
	)

	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	var (
		drv   = MakeTxDriver(ctrl, "mock-sqlx-tx")
		hook  = &beforeHook{t: t}
		hooks = dbhook.NewHooks(
			dbhook.WithHooksBefore(hook),
			dbhook.WithHooksAfter(hook),
			dbhook.WithHooksError(hook),
		)
	)

	sql.Register(driverName, dbhook.Wrap(drv, hooks))

	tmpDB, err := sql.Open(driverName, "test")
	if err != nil {
		t.Fatal(err)
	}

	var (
		db    = sqlx.NewDb(tmpDB, "postgres")
		tt    = TestToken{Token: strconv.FormatInt(time.Now().Unix(), 30)}
		ttGet = testID{}
	)

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.NamedExecContext(ctx, `INSERT INTO tokens(token) VALUES(:token)`, tt); err != nil {
		t.Fatal(err)
	}

	if err := tx.GetContext(ctx, &ttGet, `SELECT test FROM test LIMIT 1`, 123); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
