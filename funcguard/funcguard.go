package funcguard

import (
	"go/ast"
	"go/types"
	"log"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

type Analyzer struct {
	*analysis.Analyzer
	rules map[string]string

	parseCmdLineArgsExecuted bool
	doNothing                bool
	lock                     sync.Mutex
	cmdlineParams            *cmdlineParams
}

func NewAnalyzer() *Analyzer {
	params := newCmdlineParams()

	result := Analyzer{
		Analyzer: &analysis.Analyzer{
			Name:  "funcguard",
			Doc:   "Report usages of blocked functions",
			URL:   "https://github.com/simplesurance/funcguard",
			Flags: *params.flagSet,
		},
		cmdlineParams: params,
	}
	result.Analyzer.Run = result.run
	return &result
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
	a.lock.Lock()
	if !a.parseCmdLineArgsExecuted {
		err := a.parseCmdLineArgs()
		a.parseCmdLineArgsExecuted = true
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
	if a.cmdlineParams.writeCfgPath != "" {
		a.doNothing = true
		if err := defaultConfig.writeToFile(a.cmdlineParams.writeCfgPath); err != nil {
			return err
		}

		log.Printf("Wrote default config to %s", a.cmdlineParams.writeCfgPath)
		return nil
	}

	if err := a.setConfig(); err != nil {
		a.doNothing = true
		return err
	}

	return nil
}

func (a *Analyzer) setConfig() error {
	var cfg *Config
	if a.cmdlineParams.cfgPath == "" {
		cfg = &defaultConfig
		log.Printf("Using default config")
	} else {
		var err error
		cfg, err = configFromFile(a.cmdlineParams.cfgPath)
		if err != nil {
			return err
		}
		log.Printf("Loaded config from %s", a.cmdlineParams.cfgPath)
	}

	cfgMap, err := cfgToRuleMap(cfg)
	if err != nil {
		return err
	}

	a.rules = cfgMap
	return nil
}

func (a *Analyzer) analyze(pass *analysis.Pass) (any, error) {
	if !filesImportDatabaseSQL(pass.Pkg) {
		return nil, nil
	}

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

func filesImportDatabaseSQL(pkg *types.Package) bool {
	for _, importStmt := range pkg.Imports() {
		if importStmt.Path() == "database/sql" {
			return true
		}
	}

	return false
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
