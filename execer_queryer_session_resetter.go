package dbhook

type ExecerQueryerSessionResetter struct {
	*Conn
	*ExecerContext
	*QueryerContext
	*SessionResetter
}
