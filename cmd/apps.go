package cmd

import (
	"fmt"
	"strings"

	"github.com/gjhenrique/yafl/apps"
	"github.com/gjhenrique/yafl/sh"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var appCmd = &cobra.Command{
	Use:   "apps",
	Short: "Launch applications",
	Run:   runApps,
}

func runApps(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		entry, err := apps.GetEntryFromName(strings.Join(args, " "))
		if err != nil {
			panic(err)
		}

		err = sh.SpawnAsyncProcess(strings.Fields(entry.Exec), "")
		if err != nil {
			panic(err)
		}
	} else {
		entries, error := apps.FormattedApplicationNames()
		if error != nil {
			panic(error)
		}
		fmt.Print(entries)
	}
}
