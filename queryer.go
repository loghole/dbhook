package dbhook //nolint:dupl // it's ok

import (
	"context"
	"database/sql/driver"
	"fmt"
)

type QueryerContext struct {
	*Conn
}

func (conn *QueryerContext) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	var (
		err       error
		vals      = argsToValue(args)
		hookInput = &HookInput{
			Query:  query,
			Args:   vals,
			Caller: CallerQuery,
			Error:  nil,
		}
	)

	if conn.hooks != nil {
		ctx, err = conn.hooks.Before(ctx, hookInput)
		if err != nil {
			return nil, err
		}
	}

	results, err := conn.queryContext(ctx, query, args)
	if err != nil {
		if conn.hooks == nil {
			return nil, err
		}

		hookInput.Error = err

		if _, err := conn.hooks.Error(ctx, hookInput); err != nil {
			return nil, err
		}
	}

	if conn.hooks != nil {
		if _, err := conn.hooks.After(ctx, hookInput); err != nil {
			return nil, err // nolint:wrapcheck // need clear error
		}
	}

	return results, nil
}

func (conn *QueryerContext) queryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	switch c := conn.Conn.Conn.(type) {
	case driver.QueryerContext:
		return c.QueryContext(ctx, query, args) // nolint:wrapcheck // need clear error
	case driver.Queryer: // nolint:staticcheck // deprecated
		dargs, err := namedValueToValue(args)
		if err != nil {
			return nil, fmt.Errorf("can't contert named value to value: %w", err)
		}

		return c.Query(query, dargs) // nolint:wrapcheck // need clear error
	default:
		// This should not happen
		return nil, ErrNonQueryer
	}
}
