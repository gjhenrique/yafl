package cmd

import (
	"os"

	"github.com/gjhenrique/lfzf/mode"
	"github.com/gjhenrique/lfzf/sh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lfzf",
	Short: "Launcher using fzf with modes",
	Run:   runRoot,
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
		panic(err)
	}

	mode, err := mode.FindModeByKey(modes, entry.ModeKey)
	if err != nil {
		panic(err)
	}

	err = mode.Launch(entry.Id)
	if err != nil {
		panic(err)
	}
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
