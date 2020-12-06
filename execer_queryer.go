package dbhook

type ExecerQueryer struct {
	*Conn
	*ExecerContext
	*QueryerContext
}
