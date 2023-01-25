package launcher

import (
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

func ParseModesFromConfig(configFile string) ([]*Mode, error) {
	modes := make([]*Mode, 0)

	fileData, err := os.ReadFile(configFile)
	if err != nil {
		fileData = []byte("")
		err = nil
	}

	cfg := make(map[string]map[string]Mode)

	err = toml.Unmarshal(fileData, &cfg)
	if err != nil {
		return modes, err
	}

	for k := range cfg["modes"] {
		mode := cfg["modes"][k]
		mode.Key = k

		if mode.Cache == nil {
			mode.Cache = &defaultCacheTime
		}

		// Transforming f into f<space>
		// When there is a space, we don't touch it
		if mode.Prefix != "" {
			if !strings.HasSuffix(mode.Prefix, " ") {
				mode.Prefix = mode.Prefix + " "
			}
		}

		modes = append(modes, &mode)
	}

	bin, err := os.Executable()
	if err != nil {
		return modes, err
	}

	app := appMode(modes)

	if app == nil {
		app = &Mode{
			Cache:          &defaultCacheTime,
			Exec:           fmt.Sprintf("%s apps", bin),
			Key:            "apps",
			HistoryEnabled: true,
		}
		modes = append(modes, app)
	} else {
		if app.Exec != "" {
			app.Exec = fmt.Sprintf("%s apps", bin)
		}
	}

	return modes, nil
}
