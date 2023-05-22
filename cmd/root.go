package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var RootCmd = &cobra.Command{
	Use:          "",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		session := internal.GetSession()
		cmdName := cmd.Name()
		if cmdName != "" && cmdName != "authorize" && cmdName != "logout" && cmdName != "end" && cmdName != "help" && !session.IsAuthenticated {
			return fmt.Errorf("Authorization required.\n")
		}
		session.LastActivityTime = time.Now()
		return nil
	},
}

func init() {
	// remove extra help cruft
	RootCmd.SetHelpTemplate(`
Available Commands:
{{- range $index, $command := .Commands}}
	{{.Name}}{{"\t"}}{{.Short}}{{end}}
`)
}
