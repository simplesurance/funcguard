# funcguard

funcguard is a configurable Golang linter that reports the usage of specific
functions.
By default, it reports the use of functions without a context
parameter from the `database/sql` and `net/http` packages where equivalent
context-aware functions exist.

## Installation

```sh
go install github.com/simplesurance/funcguard@latest
```

## Configuration

### Create the Default Configuration File

```sh
funcguard -write-cfg funcguard.yml ./...
```

(The command must be run in a directory containing `.go` files, despite if files
will be analyzed.)

### Run with Custom Rules from Configuration Files

```sh
funcguard -config funcguard.yml ./...
```

## Default Configuration

```yaml
rules:
    - function: (*database/sql.DB).Begin
      error: use context-aware method BeginTx instead
    - function: (*database/sql.DB).Exec
      error: use context-aware method ExecContext instead of Exec
    - function: (*database/sql.DB).Ping
      error: use context-aware method PingContext instead of Ping
    - function: (*database/sql.DB).Prepare
      error: use context-aware method PrepareContext instead of Prepare
    - function: (*database/sql.DB).Query
      error: use context-aware method QueryContext instead of Query
    - function: (*database/sql.DB).QueryRow
      error: use context-aware method QueryRowContext instead of QueryRow
    - function: (*database/sql.Tx).Exec
      error: use context-aware method ExecContext instead of Exec
    - function: (*database/sql.Tx).Prepare
      error: use context-aware method PrepareContext instead of Prepare
    - function: (*database/sql.Tx).Query
      error: use context-aware method QueryContext instead of Query
    - function: (*database/sql.Tx).QueryRow
      error: use context-aware method QueryRowContext instead of QueryRow
    - function: (*database/sql.Tx).Stmt
      error: use context-aware method StmtContext instead of Stmt
    - function: net/http.Get
      error: use context-aware http.NewRequestWithContext method instead
    - function: net/http.Head
      error: use context-aware http.NewRequestWithContext method instead
    - function: net/http.Post
      error: use context-aware http.NewRequestWithContext method instead
    - function: net/http.PostForm
      error: use context-aware http.NewRequestWithContext method instead
    - function: (*net/http.Client).Get
      error: use context-aware http.NewRequestWithContext method instead
    - function: (*net/http.Client).Head
      error: use context-aware http.NewRequestWithContext method instead
    - function: (*net/http.Client).Post
      error: use context-aware http.NewRequestWithContext method instead
    - function: (*net/http.Client).PostForm
      error: use context-aware http.NewRequestWithContext method instead
```
