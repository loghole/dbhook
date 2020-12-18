package dbhook

import (
	"context"
	"database/sql/driver"
)

type Conn struct {
	Conn        driver.Conn
	ConnBeginTx driver.ConnBeginTx
	hooks       Hook
}

func (conn *Conn) Prepare(query string) (driver.Stmt, error) { return conn.Conn.Prepare(query) }
func (conn *Conn) Close() error                              { return conn.Conn.Close() }

// Begin Deprecated.
// nolint:staticcheck // deprecated
func (conn *Conn) Begin() (driver.Tx, error) { return conn.Conn.Begin() }

func (conn *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var (
		err       error
		stmt      driver.Stmt
		hookInput = &HookInput{
			Query:  query,
			Caller: CallerStmt,
			Args:   nil,
			Error:  nil,
		}
	)

	if conn.hooks != nil {
		ctx, err = conn.hooks.Before(ctx, hookInput)
		if err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if c, ok := conn.Conn.(driver.ConnPrepareContext); ok {
		stmt, err = c.PrepareContext(ctx, query)
	} else {
		stmt, err = conn.Prepare(query)
	}

	if err != nil {
		if conn.hooks == nil {
			return nil, err //nolint:wrapcheck // need clear error
		}

		hookInput.Error = err

		if _, err := conn.hooks.Error(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if conn.hooks != nil {
		if _, err := conn.hooks.After(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	return &Stmt{Stmt: stmt, hooks: conn.hooks, query: query}, nil
}

func (conn *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var (
		initCtx   = ctx
		err       error
		tx        driver.Tx
		hookInput = &HookInput{
			Caller: CallerBegin,
			Query:  "",
			Args:   nil,
			Error:  nil,
		}
	)

	if conn.hooks != nil {
		ctx, err = conn.hooks.Before(ctx, hookInput)
		if err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if c, ok := conn.Conn.(driver.ConnBeginTx); ok {
		tx, err = c.BeginTx(ctx, opts)
	} else {
		tx, err = conn.Begin()
	}

	if err != nil {
		if conn.hooks == nil {
			return nil, err //nolint:wrapcheck // need clear error
		}

		hookInput.Error = err

		if _, err := conn.hooks.Error(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	if conn.hooks != nil {
		if _, err := conn.hooks.After(ctx, hookInput); err != nil {
			return nil, err //nolint:wrapcheck // need clear error
		}
	}

	return &Tx{Tx: tx, hooks: conn.hooks, ctx: initCtx}, nil
}
