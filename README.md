# funcguard

funcguard is a configurable Golang linter that reports the usage of specific
functions. \
Functions are identified by their full path. Usage of those functions is
also identified if the package is imported with a different identifier, or their
types are embedded into other structs. \
By default, it reports the use of functions without a context
parameter from the `database/sql` and `net/http` packages where equivalent
context-aware functions exist.

## Installation

```sh
go install github.com/simplesurance/funcguard/cmd/funcguard@latest
```

## Execution

```sh
funcguard ./...
```

## Configuration

### Create the Default Configuration File

```sh
funcguard -write-cfg funcguard.yml ./...
```

(The command must be run in a directory containing `.go` files, despite if files
will be analyzed.)

### Default Configuration

```yaml
rules:
    - function-path: (*database/sql.DB).Begin
      error-msg: use context-aware method BeginTx instead of Begin
    - function-path: (*database/sql.DB).Exec
      error-msg: use context-aware method ExecContext instead of Exec
    - function-path: (*database/sql.DB).Ping
      error-msg: use context-aware method PingContext instead of Ping
    - function-path: (*database/sql.DB).Prepare
      error-msg: use context-aware method PrepareContext instead of Prepare
    - function-path: (*database/sql.DB).Query
      error-msg: use context-aware method QueryContext instead of Query
    - function-path: (*database/sql.DB).QueryRow
      error-msg: use context-aware method QueryRowContext instead of QueryRow
    - function-path: (*database/sql.Tx).Exec
      error-msg: use context-aware method ExecContext instead of Exec
    - function-path: (*database/sql.Tx).Prepare
      error-msg: use context-aware method PrepareContext instead of Prepare
    - function-path: (*database/sql.Tx).Query
      error-msg: use context-aware method QueryContext instead of Query
    - function-path: (*database/sql.Tx).QueryRow
      error-msg: use context-aware method QueryRowContext instead of QueryRow
    - function-path: (*database/sql.Tx).Stmt
      error-msg: use context-aware method StmtContext instead of Stmt
    - function-path: net/http.Get
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: net/http.Head
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: net/http.Post
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: net/http.PostForm
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: (*net/http.Client).Get
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: (*net/http.Client).Head
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: (*net/http.Client).Post
      error-msg: use context-aware http.NewRequestWithContext method instead
    - function-path: (*net/http.Client).PostForm
      error-msg: use context-aware http.NewRequestWithContext method instead
```

### Execution with Custom Rules from Configuration Files

```sh
funcguard -config funcguard.yml ./...
```

A custom configuration to only forbid the use of `fmt.Println` could be:

```yaml
rules:
    - function-path: fmt.Println
      error-msg: fmt.Println is forbidden
```
