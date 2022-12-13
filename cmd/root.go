package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gjhenrique/lfzf/mode"
	"github.com/gjhenrique/lfzf/sh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lfzf",
	Short: "Launcher using fzf with modes",
	Run:   runRoot,
}

func launchCommand(mode *mode.Mode, input string) {
	err := mode.Launch(input)
	if err != nil {
		panic(err)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	modes, err := mode.AllModes(configFile())
	if err != nil {
		panic(err)
	}

	selectedMode := mode.FindModeByInput(modes, "")
	entries, err := selectedMode.ListEntries()
	if err != nil {
		panic(err)
	}

	entry, err := sh.Fzf(entries)
	if err != nil {
		if _, ok := err.(*sh.SkippedInputError); ok {
			fmt.Println("User skipped")
			os.Exit(0)
		}

		if noMatchErr, ok := err.(*sh.NoMatchError); ok {
			mode := mode.FindModeByInput(modes, noMatchErr.Query)
			if mode.CallWithoutMatch {
				fmt.Println(noMatchErr.Query)
				query := strings.TrimPrefix(noMatchErr.Query, mode.Prefix)
				launchCommand(mode, query)
				os.Exit(0)
			}
		}
		panic(err)
	}

	mode, err := mode.FindModeByKey(modes, entry.ModeKey)
	if err != nil {
		panic(err)
	}

	launchCommand(mode, entry.Id)
}

func Execute() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(searchCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/lfzf/config.toml)")
}
