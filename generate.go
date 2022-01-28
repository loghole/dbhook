package dbhook

//nolint:lll // generate ok
//go:generate mockgen -destination mocks/sql.go -package mocks database/sql/driver Conn,ConnBeginTx,Driver,Stmt,Tx,ExecerContext,QueryerContext,Rows,Result

//go:generate mockgen -destination mocks/stmt.go -package mocks -source=tests/interfaces/stmt.go

//go:generate mockgen -destination mocks/conn.go -package mocks -source=tests/interfaces/conn.go
