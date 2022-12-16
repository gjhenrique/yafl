package cmd

import (
	// "fmt"
	"fmt"
	"strings"
	"time"

	"github.com/gjhenrique/lfzf/cache"
	"github.com/gjhenrique/lfzf/mode"
	"github.com/gjhenrique/lfzf/sh"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Introspect and interact with the cache of the launcher",
	Run:   searchCache,
}

func searchCache(cmd *cobra.Command, args []string) {
	modes, err := mode.AllModes(configFile())
	if err != nil {
		panic(err)
	}

	c := cache.CacheStore{Dir: cacheFolder()}

	a, err := c.FetchCache("apps", 1*time.Hour, func() (string, error) {
		appMode := mode.AppMode(modes)

		cmd := strings.Fields(appMode.Exec)
		return sh.SpawnSyncProcess(cmd, nil)
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(a)
}
