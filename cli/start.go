package cli

import (
	"os"
	"os/exec"
	"strings"

	"github.com/chenyu116/generator-mobile/logger"

	"github.com/chenyu116/generator-mobile/config"
	"github.com/chenyu116/generator-mobile/server"
	"github.com/spf13/cobra"
)

var (
	background bool
	debug      bool
	logLevel   string
	configPath string
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "start an instance",
	Long:    "",
	Example: "",
	RunE:    runStart,
}

func init() {
	flags := startCmd.Flags()
	boolFlag(flags, &background, FlagBackground)
	boolFlag(flags, &debug, FlagDebug)
	stringFlag(flags, &configPath, FlagConfigPath)
	stringFlag(flags, &logLevel, FlagLogLevel)
}

func runStart(cmd *cobra.Command, args []string) error {
	if background {
		args = make([]string, 0, len(os.Args))

		for _, arg := range os.Args {
			if strings.HasPrefix(arg, "--daemon") {
				continue
			}
			args = append(args, arg)
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Start()
	}

	logger.InitLogger(debug, logLevel)

	if configPath != "" {
		config.SetConfigPath(configPath)
	}

	config.ReadConfig()

	s := server.NewServer()

	s.Start()

	return nil
}
