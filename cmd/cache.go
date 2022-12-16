package cmd

import (
	// "fmt"
	"strings"

	"github.com/gjhenrique/lfzf/cache"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Introspect and interact with the launcher cache",
}

var cleanCache = &cobra.Command{
	Use:   "clean",
	Short: "Remove cache",
	Run:   removeCache,
}

func removeCache(cmd *cobra.Command, args []string) {
	modeKey := strings.Join(args, " ")
	c := cache.CacheStore{Dir: cacheFolder()}
	c.Remove(modeKey)
}

func init() {
	cacheCmd.AddCommand(cleanCache)
}
