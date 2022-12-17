package cmd

import (
	"fmt"
	"strings"

	"github.com/gjhenrique/yafl/cache"
	"github.com/gjhenrique/yafl/mode"
	"github.com/gjhenrique/yafl/sh"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Invoke search with input from fzf prompt",
	Run:   search,
}

func search(cmd *cobra.Command, args []string) {
	var query string

	if len(args) > 0 {
		query = strings.Join(args, " ")
	}

	modes, err := mode.AllModes(defaultConfigFile())
	if err != nil {
		panic(err)
	}

	selectedMode := mode.FindModeByInput(modes, query)
	c := cache.CacheStore{Dir: cacheFolder()}
	entries, err := selectedMode.ListEntries(c)
	if err != nil {
		panic(err)
	}

	s := sh.FormatEntries(entries)
	fmt.Print(s)
}
