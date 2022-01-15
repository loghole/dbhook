package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sqlite "github.com/mattn/go-sqlite3"

	"github.com/loghole/dbhook"
)

const durationName = "duration"

type Hook struct {
	log *log.Logger
}

func (h *Hook) Before(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.log.Printf("before %s: %s", input.Caller, input.Query)

	return context.WithValue(ctx, durationName, time.Now()), nil
}

func (h *Hook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.log.Printf("after %s: %s. duration: %v", input.Caller, input.Query, time.Since(ctx.Value(durationName).(time.Time)))

	return ctx, nil
}

func (h *Hook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.log.Printf("error %s: %v. duration: %v", input.Caller, input.Error, time.Since(ctx.Value(durationName).(time.Time)))

	return ctx, input.Error // if return nil, then the query will not return an error
}

func main() {
	// Init hooks
	hooks := dbhook.NewHooks(dbhook.WithHook(&Hook{log: log.Default()}))

	// Register the wrapper
	sql.Register("sqlite_with_hook", dbhook.Wrap(&sqlite.SQLiteDriver{}, hooks))

	// Connect to the registered wrapped driver
	db, _ := sql.Open("sqlite_with_hook", ":memory:")

	ctx := context.Background()
	// Queries
	db.ExecContext(ctx, "CREATE TABLE t (id INTEGER, text VARCHAR(3))")

	tx, _ := db.BeginTx(context.Background(), nil)
	tx.ExecContext(ctx, "INSERT into t (text) VALUES(?), (?)", "foo", "bar")
	tx.QueryContext(ctx, "SELECT id, text FROM t")
	tx.Commit()

	db.QueryContext(ctx, "SELEC id, text FROM t")
}

/*
  Output:
  before exec: CREATE TABLE t (id INTEGER, text VARCHAR(3))
  after exec: CREATE TABLE t (id INTEGER, text VARCHAR(3)). duration: 66.915µs
  before begin:
  after begin: . duration: 3.456µs
  before exec: INSERT into t (text) VALUES(?), (?)
  after exec: INSERT into t (text) VALUES(?), (?). duration: 13.756µs
  before query: SELECT id, text FROM t
  after query: SELECT id, text FROM t. duration: 6.382µs
  before commit:
  after commit: . duration: 24.786µs
  before query: SELEC id, text FROM t
  error query: near "SELEC": syntax error. duration: 10.57µs
*/
