# dbhook
[![GoDoc](https://pkg.go.dev/badge/github.com/loghole/dbhook)](https://pkg.go.dev/github.com/loghole/dbhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/loghole/dbhook)](https://goreportcard.com/report/github.com/loghole/dbhook)

This is a hook for any database/sql driver.  
DBhook allows to log requests, measure their duration, control the behavior of requests without changing the code base.  
This is the middelware for your database.

# Install
```sh
go get github.com/loghole/dbhook
```

# Usage
```go
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
```

# Real worl examples
- [Reconnect hook](https://github.com/loghole/database/blob/2688b139e9899d532ddfae97a21f427c8258c103/hooks/reconnect.go#L1-L42)
- [Error with self code](https://github.com/loghole/database/blob/2688b139e9899d532ddfae97a21f427c8258c103/hooks/simplerr.go#L1-L42)
- [Tracing hook](https://github.com/loghole/database/blob/76be80785f31df69d255da012ba728a65efa2785/hooks/tracing.go#L1-L72)
- [Metrics hook](https://github.com/loghole/database/blob/76be80785f31df69d255da012ba728a65efa2785/hooks/metrics.go#L1-L95)
