package funcguard

import "flag"

type Option func(*Analyzer)

func WithCmdlineFlags() Option {
	return func(a *Analyzer) {
		a.parseCmdLineFlags = true

		flagSet := flag.NewFlagSet("", flag.ExitOnError)
		flagSet.StringVar(&a.writeCfgPath, "write-cfg", "", "Write the default config to the given path and exit (package argument must still be passed, but is ignored).")
		flagSet.StringVar(&a.configPath, "config", "", "Path to the configuration file. If not specified the default rules are used.")
		a.Analyzer.Flags = *flagSet
	}
}

func WithConfig(cfg *Config) Option {
	return func(a *Analyzer) {
		a.cfg = cfg
	}
}
