package dbhook

import (
	"database/sql/driver"
	"fmt"
)

// Driver is the struct implemented driver.Driver.
type Driver struct {
	driver.Driver
	hooks Hook
}

// Open returns a new connection to the database.
// The name is a string in a driver-specific format.
//
// Open may return a cached connection (one previously
// closed), but doing so is unnecessary; the sql package
// maintains a pool of idle connections for efficient re-use.
//
// The returned connection is only used by one goroutine at a
// time.
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
