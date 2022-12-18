package launcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/gjhenrique/yafl/cache"
	"github.com/gjhenrique/yafl/sh"
	toml "github.com/pelletier/go-toml/v2"
)

type Launcher struct {
	cache      *cache.CacheStore
	configFile string
	modes      []*Mode
}

func NewLauncher(configFile string, cacheFolder string) (*Launcher, error) {
	c := cache.CacheStore{
		Dir: cacheFolder,
	}
	modes, err := allModes(configFile)
	if err != nil {
		return nil, err
	}

	return &Launcher{
		cache:      &c,
		configFile: configFile,
		modes:      modes,
	}, nil
}

func (l *Launcher) Fzf(entries []*sh.Entry) (*sh.Entry, *Mode, error) {
	entry, err := sh.Fzf(entries)

	if err != nil {
		if noMatchErr, ok := err.(*sh.NoMatchError); ok {
			m := l.findModeByInput(noMatchErr.Query)
			if m.CallWithoutMatch {
				query := strings.TrimPrefix(noMatchErr.Query, m.Prefix)
				entry := &sh.Entry{ModeKey: m.Key, Id: query, Text: noMatchErr.Query}
				return entry, m, nil
			}
		}

		return nil, nil, err
	}

	m, err := l.findModeByKey(entry.ModeKey)
	return entry, m, nil
}

func (l *Launcher) ListEntries(input string) ([]*sh.Entry, error) {
	selectedMode := l.findModeByInput(input)
	return selectedMode.ListEntries(l.cache)
}

func (l *Launcher) findModeByKey(key string) (*Mode, error) {
	var mode *Mode

	for _, m := range l.modes {
		if m.Key == key {
			mode = m
		}
	}

	if mode == nil {
		return nil, fmt.Errorf("Mode %s not found", key)
	}

	return mode, nil
}

func (l *Launcher) findModeByInput(input string) *Mode {
	var mode *Mode

	for _, m := range l.modes {
		if m.Prefix == "" {
			continue
		}

		if strings.HasPrefix(input, m.Prefix) {
			mode = m
		}
	}

	if mode == nil {
		mode = appMode(l.modes)
	}

	return mode
}

func appMode(modes []*Mode) *Mode {
	var appMode *Mode

	for _, m := range modes {
		if m.Key == "apps" {
			appMode = m
		}
	}

	return appMode
}

func allModes(configFile string) ([]*Mode, error) {
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
			Cache: 60,
			Exec:  fmt.Sprintf("%s apps", bin),
			Key:   "apps",
		}
		modes = append(modes, app)
	} else {
		if app.Exec != "" {
			app.Exec = fmt.Sprintf("%s apps", bin)
		}
	}

	return modes, nil
}
