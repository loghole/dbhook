package dbhook

//nolint:lll // generate ok
//go:generate mockgen -destination mocks/sql.go -package mocks database/sql/driver Conn,ConnBeginTx,Driver,Stmt,Tx,ExecerContext,QueryerContext,Rows,Result
