package cli

import (
	"github.com/spf13/pflag"
)

type flagInfoString struct {
	Name        string
	Shorthand   string
	Description string
	Default     string
}

type flagInfoBool struct {
	Name        string
	Shorthand   string
	Description string
	Default     bool
}

var (
	FlagConfigPath = flagInfoString{
		Name:        "config-file",
		Shorthand:   "c",
		Description: "path to configuration file",
	}

	FlagLogLevel = flagInfoString{
		Name:        "log-level",
		Shorthand:   "l",
		Description: "logging level",
		Default:     "info",
	}

	FlagBackground = flagInfoBool{
		Name:        "daemon",
		Shorthand:   "",
		Description: "run process in background",
		Default:     false,
	}
	FlagDebug = flagInfoBool{
		Name:        "debug",
		Shorthand:   "d",
		Description: "switch debug mode",
		Default:     false,
	}
)

func stringFlag(f *pflag.FlagSet, valPtr *string, flagInfo flagInfoString) {
	f.StringVarP(valPtr,
		flagInfo.Name,
		flagInfo.Shorthand,
		flagInfo.Default,
		flagInfo.Description)
}

func boolFlag(f *pflag.FlagSet, valPtr *bool, flagInfo flagInfoBool) {
	f.BoolVarP(valPtr,
		flagInfo.Name,
		flagInfo.Shorthand,
		flagInfo.Default,
		flagInfo.Description)
}
