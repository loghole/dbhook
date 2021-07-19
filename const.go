package dbhook

const testDataPath = "./test/.snapshots" // nolint:deadcode,unused,varcheck // used in test

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
