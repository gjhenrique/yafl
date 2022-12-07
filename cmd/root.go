package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"
	"runtime"
)

const APP_NAME = "lfzf"

func configFolder() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	} else if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config")
	}
}

func appFolder() string {
	configFolder := configFolder()
	appFolder := filepath.Join(configFolder, APP_NAME)

	err := os.MkdirAll(appFolder, 0755)

	if err != nil {
		panic("Error when creating database folder" + err.Error())
	}

	return appFolder
}

func configFile() string {
	return filepath.Join(appFolder(), "config.toml")
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lfzf",
	Short: "Launcher using fzf with modes",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	fmt.Println("Root cmd")
	// Get the apps mode
	// Execute the apps mode and pipe the return to fzf
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(appCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/lfzf/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
