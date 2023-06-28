package testdata

import "database/sql"

// not initializing it is fine, the code is not executed only analyzed
var db *sql.DB

func basic() {
	db.Exec("SELECT * FROM users") // want "use context-aware method ExecContext instead of Exec"
	tx, _ := db.Begin()            // want "use context-aware method BeginTx instead of Begin"
	tx.Exec("")                    // want "use context-aware method ExecContext instead of Exec"
	tx.Prepare("")                 // want "use context-aware method PrepareContext instead of Prepare"
	tx.Query("")                   // want "use context-aware method QueryContext instead of Query"
	tx.QueryRow("")                // want "use context-aware method QueryRowContext instead of QueryRow"
	tx.Stmt(nil)                   // want "use context-aware method StmtContext instead of Stmt"
	db.Ping()                      // want "use context-aware method PingContext instead of Ping"
	db.Prepare("")                 // want "use context-aware method PrepareContext instead of Prepare"
	db.Query("")                   // want "use context-aware method QueryContext instead of Query"
	db.QueryRow("")                // want "use context-aware method QueryRowContext instead of QueryRow"
}

func withAssignment() {
	_, _ = db.Exec("SELECT * FROM users") // want "use context-aware method ExecContext instead of Exec"
}

func callInDefer() {
	defer func() {
		_, _ = db.Exec("SELECT * FROM users") // want "use context-aware method ExecContext instead of Exec"

	}()
}

func dbInStruct() {
	s := struct {
		db *sql.DB
	}{}
	s.db.Exec("SELECT * FROM users") // want "use context-aware method ExecContext instead of Exec"
}
