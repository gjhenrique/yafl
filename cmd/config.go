package cmd

import (
	"os"
	"path/filepath"
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
