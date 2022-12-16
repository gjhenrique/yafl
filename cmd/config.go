package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

const APP_NAME = "yafl"

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

	if _, err := os.Stat(appFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(appFolder, 0755); err != nil {
			panic("Error when creating cache folder" + err.Error())
		}
	}

	return appFolder
}

func configFile() string {
	return filepath.Join(appFolder(), "config.toml")
}

func cacheFolder() string {
	// Only Linux related. Complete with other directories whenever it's supported
	systemCacheFolder := filepath.Join(os.Getenv("HOME"), ".cache")
	appCacheFolder := filepath.Join(systemCacheFolder, APP_NAME)

	if _, err := os.Stat(appCacheFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(appCacheFolder, 0755); err != nil {
			panic("Error when creating cache folder" + err.Error())
		}
	}

	return appCacheFolder
}

var cfgFile string
