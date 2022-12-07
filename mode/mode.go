package mode

import (
	"fmt"
	"os"
	"strings"

	sh "github.com/gjhenrique/lfzf/sh"
	toml "github.com/pelletier/go-toml/v2"
)

type Mode struct {
	Cache  int
	Exec   string
	Prefix string
	Key    string
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
		modes = append(modes, &mode)
	}

	bin, err := os.Executable()
	if err != nil {
		return modes, err
	}

	appMode := AppMode(modes)

	if appMode == nil {
		appMode = &Mode{
			Cache: 30,
			Exec:  fmt.Sprintf("%s apps search", bin),
			Key:   "apps",
		}
		modes = append(modes, appMode)
	} else {
		if appMode.Exec != "" {
			appMode.Exec = fmt.Sprintf("%s apps search", bin)
		}
	}

	return modes, nil
}

func (m *Mode) ListEntries() error {
	cmd := strings.Fields(m.Exec)
	result, err := sh.CommandWithOutput(cmd, nil)
	if err != nil {
		return err
	}

	os.Stdout.Write([]byte(result))

	return nil
}

func (m *Mode) Launch(input string) error {
	cmd := strings.Fields(m.Exec)

	input = strings.TrimPrefix(input, m.Prefix)
	cmd = append(cmd, input)

	err := sh.CommandDeattached(strings.Join(cmd, " "))
	if err != nil {
		return err
	}

	return nil
}

func FindMode(input string, modes []*Mode) *Mode {
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
		for _, m := range modes {
			if m.Key == "apps" {
				mode = m
			}
		}
	}

	return mode
}
