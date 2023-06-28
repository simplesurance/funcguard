package funcguard

import (
	"go/ast"
	"go/types"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

type Analyzer struct {
	*analysis.Analyzer
	cfg   *Config
	rules map[string]string

	cmdLineFlags struct {
		parseFlags   bool
		writeCfgPath string
		configPath   string
	}

	doNothing bool
	lock      sync.Mutex

	logf LogFunc
}

type LogFunc func(format string, a ...any)

func NewAnalyzer(opts ...Option) (*Analyzer, error) {
	result := Analyzer{
		Analyzer: &analysis.Analyzer{
			Name: "funcguard",
			Doc:  "Report usages of prohibited functions",
			URL:  "https://github.com/simplesurance/funcguard",
		},
		logf: func(string, ...any) {},
	}

	for _, opt := range opts {
		opt(&result)
	}

	if !result.cmdLineFlags.parseFlags && (result.cfg == nil || len(result.cfg.Rules) == 0) {
		// When the linter is invoked via SingleChecker, the the
		// cmdLineFlags are not set yet, only when run() is invoked.
		// invoked. Therefore we can not use default config when
		// cmdLineFlags.configPath is empty here.
		result.cfg = &defaultConfig
		result.logf("Using default config")
	}

	if result.cfg != nil {
		var err error
		result.rules, err = cfgToRuleMap(result.cfg)
		if err != nil {
			return nil, err
		}
		result.cfg = nil // not needed anymore
	}

	result.Analyzer.Run = result.run

	return &result, nil
}

func (a *Analyzer) run(pass *analysis.Pass) (any, error) {
	// SingleChecker does not support to register flags and handle them before the Analyzer is run.
	// The only way to handle our own flags is in this run() method which is invoked multiple times.
	// To keep the code simple by still using singlechecker, we set
	// doNothing to true when an error occurred while processing our flags.
	// If this happens, the error is returned 1x, and then following run()
	// calls will return immediately.
	// This also causes that it's not possible to invoke the linter without
	// providing a package specifier, therefore invoking it only to write
	// the default config to a file without a package spec, is not
	// possible.
	// Refactor this after: https://github.com/golang/go/issues/53336
	a.lock.Lock() // TODO: is the lock needed? Is run called in parallel?
	if a.cmdLineFlags.parseFlags {
		a.cmdLineFlags.parseFlags = false
		err := a.parseCmdLineArgs()
		if err != nil {
			a.lock.Unlock()
			return nil, err
		}
	}
	a.lock.Unlock()

	if a.doNothing {
		return nil, nil
	}

	return a.analyze(pass)
}

func (a *Analyzer) parseCmdLineArgs() error {
	if a.cmdLineFlags.writeCfgPath != "" {
		a.doNothing = true
		if err := defaultConfig.writeToFile(a.cmdLineFlags.writeCfgPath); err != nil {
			return err
		}

		a.logf("Wrote default config to %s", a.cmdLineFlags.writeCfgPath)
		return nil
	}

	return a.setConfig()
}

func (a *Analyzer) setConfig() error {
	var cfg *Config

	if a.configPath != "" {
		var err error
		cfg, err = configFromFile(a.configPath)
		if err != nil {
			return err
		}
		a.logf("Loaded config from %s", a.configPath)

	} else {
		cfg = &defaultConfig
		a.logf("Using default config")
	}

	cfgMap, err := cfgToRuleMap(cfg)
	if err != nil {
		return err
	}

	a.rules = cfgMap
	return nil
}

func (a *Analyzer) analyze(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return true
			}

			call, fn := toFuncCall(n, pass.TypesInfo)
			if fn == nil {
				return true
			}

			allowed, errorMsg := a.isAllowed(fn.FullName())
			if !allowed {
				pass.Reportf(call.Pos(), errorMsg)
			}

			return false
		})
	}

	return nil, nil
}

func (a *Analyzer) isAllowed(fullFuncName string) (allowed bool, errorMsg string) {
	errorMsg, exists := a.rules[fullFuncName]
	return !exists, errorMsg
}

func toFuncCall(n ast.Node, typesInfo *types.Info) (ast.Node, *types.Func) {
	switch call := n.(type) {
	case *ast.CallExpr:
		fn, _ := typeutil.Callee(typesInfo, call).(*types.Func)
		return call, fn

	case *ast.AssignStmt:
		if ce, ok := call.Rhs[0].(*ast.CallExpr); ok {
			fn, _ := typeutil.Callee(typesInfo, ce).(*types.Func)
			return ce, fn
		}
	}

	return nil, nil
}
