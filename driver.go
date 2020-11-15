package dbhook

import (
	"database/sql/driver"
	"fmt"
)

type Driver struct {
	driver.Driver
	hooks Hook
}

func (drv *Driver) Open(name string) (driver.Conn, error) {
	conn, err := drv.Driver.Open(name)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	wrapped := &Conn{
		Conn:        conn,
		hooks:       drv.hooks,
		ConnBeginTx: nil,
	}

	switch {
	case isExecerAndQueryerAndSessionResetter(conn):
		return &ExecerQueryerSessionResetter{
			Conn:            wrapped,
			ExecerContext:   &ExecerContext{Conn: wrapped},
			QueryerContext:  &QueryerContext{Conn: wrapped},
			SessionResetter: &SessionResetter{Conn: wrapped},
		}, nil
	case isExecerAndQueryer(conn):
		return &ExecerQueryer{
			Conn:           wrapped,
			ExecerContext:  &ExecerContext{Conn: wrapped},
			QueryerContext: &QueryerContext{Conn: wrapped},
		}, nil
	case isExecer(conn):
		return &ExecerContext{Conn: wrapped}, nil
	case isQueryer(conn):
		return &QueryerContext{Conn: wrapped}, nil
	}

	return wrapped, nil
}

func isExecerAndQueryer(conn driver.Conn) bool {
	return isExecer(conn) && isQueryer(conn)
}

func isExecerAndQueryerAndSessionResetter(conn driver.Conn) bool {
	return isSessionResetter(conn) && isExecer(conn) && isQueryer(conn)
}

func isExecer(conn driver.Conn) bool {
	_, ok := conn.(driver.ExecerContext)

	return ok
}

func isQueryer(conn driver.Conn) bool {
	_, ok := conn.(driver.QueryerContext)

	return ok
}

func isSessionResetter(conn driver.Conn) bool {
	_, ok := conn.(driver.SessionResetter)

	return ok
}
