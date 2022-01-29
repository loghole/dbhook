package dbhook

const testDataPath = "./tests/.snapshots" // nolint:deadcode,unused,varcheck // used in test

// CallerType is callers type list
type CallerType string

const (
	CallerStmt      CallerType = "stmt"
	CallerStmtExec  CallerType = "stmt_exec"
	CallerStmtQuery CallerType = "stmt_query"
	CallerExec      CallerType = "exec"
	CallerQuery     CallerType = "query"
	CallerBegin     CallerType = "begin"
	CallerCommit    CallerType = "commit"
	CallerRollback  CallerType = "rollback"
)
