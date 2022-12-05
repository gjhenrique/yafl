package main

import (
	"bytes"
	"code.rocketnine.space/tslocum/desktop"
	"errors"
	"fmt"
	// "github.com/gjhenrique/lfzf/cmd"
	toml "github.com/pelletier/go-toml/v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type DesktopEntry struct {
	Name string
	Exec string
}

func GetDesktopEntries() ([]*desktop.Entry, error) {
	allEntries := make([]*desktop.Entry, 0, 100)

	dirs := desktop.DataDirs()
	entries, err := desktop.Scan(dirs)

	if err != nil {
		return nil, err
	}

	for _, dir := range entries {
		for _, entry := range dir {
			allEntries = append(allEntries, entry)
		}
	}

	return allEntries, nil
}

func commandWithOutput(command []string, input []byte) (string, error) {
	var result strings.Builder

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = &result
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	io.Copy(stdin, bytes.NewReader(input))
	if err = stdin.Close(); err != nil {
		return "", err
	}

	err = cmd.Run()

	return result.String(), err
}

func fzf(input []byte) (string, error) {
	ownExe, err := os.Executable()
	if err != nil {
		return "", err
	}
	bind := fmt.Sprintf("change:reload:sleep 0.05; %s search {q} || true", ownExe)
	command := []string{"fzf", "--ansi", "--sort", "--extended", "--no-multi", "--cycle", "--no-info", "--bind", bind}

	return commandWithOutput(command, input)
}

func applicationName(entry *desktop.Entry) string {
	generic := ""

	if entry.GenericName != "" {
		generic = fmt.Sprintf("(%s)", entry.GenericName)
	}

	return fmt.Sprintf("%s %s", entry.Name, generic)

}

func launchApplication(execString string) error {
	execList := strings.Fields(execString)

	cmd := exec.Command(execList[0], execList[1:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err := cmd.Start()

	return err
}

func applicationNames(entries []*desktop.Entry) string {
	names := make([]string, len(entries))

	for i, entry := range entries {
		names[i] = applicationName(entry)
	}

	return strings.Join(names, "\n")

}

func GetEntryFromName(chosenApp string) (*desktop.Entry, error) {
	entries, err := GetDesktopEntries()

	if err != nil {
		return nil, err
	}

	var entry *desktop.Entry

	for _, e := range entries {
		if applicationName(e) == chosenApp {
			entry = e
		}
	}

	if entry == nil {
		return nil, errors.New("Didn't find any app name")
	}

	return entry, nil
}

type Mode struct {
	Cache  int
	Exec   string
	Prefix string
	Key    string
}

func findModes(configFile string) ([]*Mode, error) {
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

	var appMode *Mode
	for _, m := range modes {
		if m.Key == "apps" {
			appMode = m
		}
	}

	bin, err := os.Executable()
	if err != nil {
		return modes, err
	}

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

const APP_NAME = "lfzf"

func appFolder() string {
	configFolder := configFolder()
	appFolder := filepath.Join(configFolder, APP_NAME)

	err := os.MkdirAll(appFolder, 0755)

	if err != nil {
		panic("Error when creating database folder" + err.Error())
	}

	return appFolder
}

func (m *Mode) List() error {
	cmd := strings.Split(m.Exec, " ")
	result, err := commandWithOutput(cmd, nil)
	if err != nil {
		return err
	}

	os.Stdout.Write([]byte(result))

	return nil
}

func (m *Mode) Launch(input string) error {
	cmd := strings.Split(m.Exec, " ")

	input = strings.TrimPrefix(input, m.Prefix)
	cmd = append(cmd, input)
	_, err := commandWithOutput(cmd, []byte(input))
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

func main() {
	// if len(os.Args[1:]) > 0 {
	// 	fmt.Println(modes[0])
	// } else {
	// 	fzf([]byte("a\nb"))
	// }

	// modes, _ := findModes(filepath.Join(appFolder(), "config.toml"))
	// mode := FindMode("f ola", modes)
	// mode.Launch("f ola ola ola")

	// cmd.Execute()

	// appFolder()
	// entries, err := GetDesktopEntries()

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// names := applicationNames(entries)

	// appName, _ := fzf([]byte(names))
	// appName = strings.TrimSuffix(appName, "\n")

	// entry, _ := GetEntryFromName(appName)

	// launchApplication(entry.Exec)
}
