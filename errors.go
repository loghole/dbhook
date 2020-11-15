package dbhook

import "errors"

var (
	ErrNonExecer                 = errors.New("ExecerContext created for a non Execer driver.Conn")
	ErrNonQueryer                = errors.New("QueryerContext created for a non Queryer driver.Conn")
	ErrNonSessionResetter        = errors.New("ResetSession created for a non SessionResetter driver.Conn")
	ErrNamedParametersNotSupport = errors.New("sql: driver does not support the use of Named Parameters")
)
