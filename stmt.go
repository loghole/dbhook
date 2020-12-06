package dbhook

import (
	"context"
	"database/sql/driver"
)

type Stmt struct {
	Stmt  driver.Stmt
	hooks Hook
	query string
}

func (stmt *Stmt) Close() error  { return stmt.Stmt.Close() }
func (stmt *Stmt) NumInput() int { return stmt.Stmt.NumInput() }

// Exec Deprecated.
// nolint:staticcheck // deprecated
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) { return stmt.Stmt.Exec(args) }

// Query Deprecated.
// nolint:staticcheck // deprecated
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return stmt.Stmt.Query(args)
}

// ExecContext must honor the context timeout and return when it is canceled.
// nolint:dupl // it's ok
func (stmt *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	var (
		err       error
		vals      = argsToValue(args)
		hookInput = &HookInput{
			Query:  stmt.query,
			Args:   vals,
			Caller: CallerStmtExec,
			Error:  nil,
		}
	)

	if stmt.hooks != nil {
		ctx, err = stmt.hooks.Before(ctx, hookInput)
		if err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	results, err := stmt.execContext(ctx, args)
	if err != nil {
		if stmt.hooks == nil {
			return nil, err //nolint:wrapcheck // need clear error
		}

		hookInput.Error = err

		if _, err := stmt.hooks.Error(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if stmt.hooks != nil {
		if _, err := stmt.hooks.After(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	return results, nil
}

// QueryContext must honor the context timeout and return when it is canceled.
// nolint:dupl // it's ok
func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	var (
		err       error
		vals      = argsToValue(args)
		hookInput = &HookInput{
			Query:  stmt.query,
			Args:   vals,
			Caller: CallerStmtExec,
			Error:  nil,
		}
	)

	if stmt.hooks != nil {
		ctx, err = stmt.hooks.Before(ctx, hookInput)
		if err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	rows, err := stmt.queryContext(ctx, args)
	if err != nil {
		if stmt.hooks == nil {
			return nil, err //nolint:wrapcheck // need clear error
		}

		hookInput.Error = err

		if _, err := stmt.hooks.Error(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if stmt.hooks != nil {
		if _, err := stmt.hooks.After(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	return rows, nil
}

func (stmt *Stmt) queryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s, ok := stmt.Stmt.(driver.StmtQueryContext); ok {
		return s.QueryContext(ctx, args)
	}

	return stmt.Query(argsToValue(args))
}

func (stmt *Stmt) execContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s, ok := stmt.Stmt.(driver.StmtExecContext); ok {
		return s.ExecContext(ctx, args)
	}

	return stmt.Exec(argsToValue(args))
}
