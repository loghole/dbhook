package interfaces

import "database/sql/driver"

type AdvancedConn interface {
	driver.Conn
	driver.ConnPrepareContext
	driver.SessionResetter
	driver.QueryerContext
	driver.ExecerContext
	driver.ConnBeginTx
	driver.NamedValueChecker
}

type QueryerConn interface {
	driver.Conn
	driver.SessionResetter
	driver.Queryer
	driver.Execer
}

type QConn interface {
	driver.Conn
	driver.QueryerContext
}

type EConn interface {
	driver.Conn
	driver.ExecerContext
}
