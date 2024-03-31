package cmd

import (
	"os"

	"github.com/gjhenrique/yafl/sh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yafl",
	Short: "Launcher using fzf with modes",
	Run:   runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	l := newLauncher()

	entries, err := l.ListEntries([]byte(""))
	if err != nil {
		displayError(err)
	}

	err = l.ProcessEntries(entries)
	if err != nil {
		if _, ok := err.(*sh.SkippedInputError); ok {
			os.Exit(1)
		}

		displayError(err)
	}
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
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "cache directory (default is $HOME/.cache)")
}
