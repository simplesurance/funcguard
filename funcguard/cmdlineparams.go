package funcguard

import "flag"

type cmdlineParams struct {
	flagSet      *flag.FlagSet
	writeCfgPath string
	cfgPath      string
}

func newCmdlineParams() *cmdlineParams {
	result := cmdlineParams{}
	result.flagSet = flag.NewFlagSet("", flag.ExitOnError)
	result.flagSet.StringVar(&result.writeCfgPath, "write-cfg", "", "Write the default config to the given path and exit (package argument must still be passed, but is ignored).")
	result.flagSet.StringVar(&result.cfgPath, "config", "", "Path to the configuration file. If not specified the default rules are used.")
	return &result
}
