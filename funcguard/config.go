package funcguard

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Function string `yaml:"function"`
	Error    string `yaml:"error"`
}

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

var defaultConfig = Config{
	Rules: []*Rule{
		{Function: "(*database/sql.DB).Begin", Error: "use context-aware method BeginTx instead of Begin"},
		{Function: "(*database/sql.DB).Exec", Error: "use context-aware method ExecContext instead of Exec"},
		{Function: "(*database/sql.DB).Ping", Error: "use context-aware method PingContext instead of Ping"},
		{Function: "(*database/sql.DB).Prepare", Error: "use context-aware method PrepareContext instead of Prepare"},
		{Function: "(*database/sql.DB).Query", Error: "use context-aware method QueryContext instead of Query"},
		{Function: "(*database/sql.DB).QueryRow", Error: "use context-aware method QueryRowContext instead of QueryRow"},
		{Function: "(*database/sql.Tx).Exec", Error: "use context-aware method ExecContext instead of Exec"},
		{Function: "(*database/sql.Tx).Prepare", Error: "use context-aware method PrepareContext instead of Prepare"},
		{Function: "(*database/sql.Tx).Query", Error: "use context-aware method QueryContext instead of Query"},
		{Function: "(*database/sql.Tx).QueryRow", Error: "use context-aware method QueryRowContext instead of QueryRow"},
		{Function: "(*database/sql.Tx).Stmt", Error: "use context-aware method StmtContext instead of Stmt"},

		{Function: "net/http.Get", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "net/http.Head", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "net/http.Post", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "net/http.PostForm", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "(*net/http.Client).Get", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "(*net/http.Client).Head", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "(*net/http.Client).Post", Error: "use context-aware http.NewRequestWithContext method instead"},
		{Function: "(*net/http.Client).PostForm", Error: "use context-aware http.NewRequestWithContext method instead"},
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
		if _, exist := result[rule.Function]; exist {
			return nil, fmt.Errorf("duplicate rule for function %s", rule.Function)
		}

		if rule.Error != "" {
			result[rule.Function] = rule.Error
			continue
		}

		result[rule.Function] = fmt.Sprintf("use of function %s is forbidden", rule.Function)
	}

	return result, nil
}
