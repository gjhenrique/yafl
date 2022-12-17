package cmd

import (
	"os"
	"strings"

	"github.com/gjhenrique/yafl/cache"
	"github.com/gjhenrique/yafl/launcher"
	"github.com/gjhenrique/yafl/sh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yafl",
	Short: "Launcher using fzf with modes",
	Run:   runRoot,
}

func launchCommand(m *launcher.Mode, input string) {
	err := m.Launch(input)
	if err != nil {
		panic(err)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	modes, err := launcher.AllModes(defaultConfigFile())
	if err != nil {
		panic(err)
	}

	selectedMode := launcher.FindModeByInput(modes, "")
	c := cache.CacheStore{Dir: cacheFolder()}
	entries, err := selectedMode.ListEntries(c)
	if err != nil {
		panic(err)
	}

	entry, err := sh.Fzf(entries)
	if err != nil {
		if _, ok := err.(*sh.SkippedInputError); ok {
			os.Exit(1)
		}

		if noMatchErr, ok := err.(*sh.NoMatchError); ok {
			m := launcher.FindModeByInput(modes, noMatchErr.Query)
			if m.CallWithoutMatch {
				query := strings.TrimPrefix(noMatchErr.Query, m.Prefix)
				launchCommand(m, query)
				os.Exit(0)
			}
		}
		panic(err)
	}

	m, err := launcher.FindModeByKey(modes, entry.ModeKey)
	if err != nil {
		panic(err)
	}

	launchCommand(m, entry.Id)
}

func Execute() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(cacheCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/yafl/config.toml)")
}
