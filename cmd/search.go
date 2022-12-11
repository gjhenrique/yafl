package cmd

import (
	"fmt"
	"strings"

	"github.com/gjhenrique/lfzf/mode"

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

	modes, err := mode.AllModes(configFile())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(modes)
	selectedMode := mode.FindMode(query, modes)
	err = selectedMode.ListEntries()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
