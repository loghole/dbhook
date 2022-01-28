package interfaces

import "database/sql/driver"

type StmtWithContext interface {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext
}
