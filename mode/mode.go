package mode

import (
	"fmt"
	"os"
	"strings"

	sh "github.com/gjhenrique/lfzf/sh"
	toml "github.com/pelletier/go-toml/v2"
)

type Mode struct {
	Cache            int
	Exec             string
	Prefix           string
	Key              string
	CallWithoutMatch bool `toml:"call_without_match"`
}

func AppMode(modes []*Mode) *Mode {
	var appMode *Mode

	for _, m := range modes {
		if m.Key == "apps" {
			appMode = m
		}
	}

	return appMode
}

func AllModes(configFile string) ([]*Mode, error) {
	modes := make([]*Mode, 0)

	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return modes, err
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

	appMode := AppMode(modes)

	if appMode == nil {
		appMode = &Mode{
			Cache: 60,
			Exec:  fmt.Sprintf("%s apps", bin),
			Key:   "apps",
		}
		modes = append(modes, appMode)
	} else {
		if appMode.Exec != "" {
			appMode.Exec = fmt.Sprintf("%s apps", bin)
		}
	}

	return modes, nil
}

func (m *Mode) ListEntries() ([]sh.Entry, error) {
	cmd := strings.Fields(m.Exec)
	result, err := sh.SpawnSyncProcess(cmd, nil)

	entriesFromCmd := strings.Split(result, "\n")

	entries := make([]sh.Entry, len(entriesFromCmd))

	for i, r := range entriesFromCmd {
		splittedEntry := strings.Split(r, "\\x31")

		text := r
		id := text
		if len(splittedEntry) > 1 {
			id = splittedEntry[0]
			text = strings.Join(splittedEntry[1:], " ")
		}

		if m.Prefix != "" {
			text = m.Prefix + text
		}

		entries[i] = sh.Entry{ModeKey: m.Key, Text: text, Id: id}
	}

	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (m *Mode) Launch(input string) error {
	cmd := strings.Fields(m.Exec)

	err := sh.SpawnAsyncProcess(cmd, input)
	if err != nil {
		return err
	}

	return nil
}

func FindModeByKey(modes []*Mode, key string) (*Mode, error) {
	var mode *Mode

	for _, m := range modes {
		if m.Key == key {
			mode = m
		}
	}

	if mode == nil {
		return nil, fmt.Errorf("Mode %s not found", key)
	}

	return mode, nil
}

func FindModeByInput(modes []*Mode, input string) *Mode {
	var mode *Mode

	for _, m := range modes {
		if m.Prefix == "" {
			continue
		}

		if strings.HasPrefix(input, m.Prefix) {
			mode = m
		}
	}

	if mode == nil {
		mode = AppMode(modes)
	}

	return mode
}
