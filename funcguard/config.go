package funcguard

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	FunctionPath string `yaml:"function-path"`
	ErrorMsg     string `yaml:"error-msg"`
}

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

var defaultConfig = Config{
	Rules: []*Rule{
		{FunctionPath: "(*database/sql.DB).Begin", ErrorMsg: "use context-aware method BeginTx instead of Begin"},
		{FunctionPath: "(*database/sql.DB).Exec", ErrorMsg: "use context-aware method ExecContext instead of Exec"},
		{FunctionPath: "(*database/sql.DB).Ping", ErrorMsg: "use context-aware method PingContext instead of Ping"},
		{FunctionPath: "(*database/sql.DB).Prepare", ErrorMsg: "use context-aware method PrepareContext instead of Prepare"},
		{FunctionPath: "(*database/sql.DB).Query", ErrorMsg: "use context-aware method QueryContext instead of Query"},
		{FunctionPath: "(*database/sql.DB).QueryRow", ErrorMsg: "use context-aware method QueryRowContext instead of QueryRow"},
		{FunctionPath: "(*database/sql.Tx).Exec", ErrorMsg: "use context-aware method ExecContext instead of Exec"},
		{FunctionPath: "(*database/sql.Tx).Prepare", ErrorMsg: "use context-aware method PrepareContext instead of Prepare"},
		{FunctionPath: "(*database/sql.Tx).Query", ErrorMsg: "use context-aware method QueryContext instead of Query"},
		{FunctionPath: "(*database/sql.Tx).QueryRow", ErrorMsg: "use context-aware method QueryRowContext instead of QueryRow"},
		{FunctionPath: "(*database/sql.Tx).Stmt", ErrorMsg: "use context-aware method StmtContext instead of Stmt"},

		{FunctionPath: "net/http.Get", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "net/http.Head", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "net/http.Post", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "net/http.PostForm", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "(*net/http.Client).Get", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "(*net/http.Client).Head", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "(*net/http.Client).Post", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
		{FunctionPath: "(*net/http.Client).PostForm", ErrorMsg: "use context-aware http.NewRequestWithContext method instead"},
	},
}

func (c *Config) write(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}

func (c *Config) writeToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	err = c.write(f)
	if err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

func readConfig(r io.Reader) (*Config, error) {
	var result Config
	return &result, yaml.NewDecoder(r).Decode(&result)
}

func configFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readConfig(f)
}

func cfgToRuleMap(cfg *Config) (map[string]string, error) {
	result := make(map[string]string, len(cfg.Rules))
	for _, rule := range cfg.Rules {
		if _, exist := result[rule.FunctionPath]; exist {
			return nil, fmt.Errorf("duplicate rule for function %s", rule.FunctionPath)
		}

		if rule.ErrorMsg != "" {
			result[rule.FunctionPath] = rule.ErrorMsg
			continue
		}

		result[rule.FunctionPath] = fmt.Sprintf("use of function %s is forbidden", rule.FunctionPath)
	}

	return result, nil
}
