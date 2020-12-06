package test

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

type TestToken struct {
	Token string `db:"token"`
}

func TestSQLx(t *testing.T) {
	var (
		driverName = "test-sqlx"
		ctrl       = gomock.NewController(t)
		ctx        = context.WithValue(context.Background(), key, ctxValue)
	)

	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	var (
		drv   = MakeDefaultDriver(ctrl, "mock-sqlx")
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
		db = sqlx.NewDb(tmpDB, "postgres")
		tt = TestToken{Token: strconv.FormatInt(time.Now().Unix(), 30)}
	)

	if _, err := db.NamedExecContext(ctx, `INSERT INTO tokens(token) VALUES(:token)`, tt); err != nil {
		t.Fatal(err)
	}
}
