package main

import (
	"bytes"
	"errors"
	"fmt"
	toml "github.com/pelletier/go-toml/v2"

	"github.com/spf13/cobra"

	"io"
	"os"
	"os/exec"

	"strings"
	"syscall"

	"code.rocketnine.space/tslocum/desktop"
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

func fzf(input []byte) (string, error) {
	var result strings.Builder

	ownExe, err := os.Executable()
	bind := fmt.Sprintf("change:reload:sleep 0.1; %s search {q} || true", ownExe)
	cmd := exec.Command("fzf", "--ansi", "--sort", "--extended", "--no-multi", "--cycle", "--no-info", "--bind", bind)
	fmt.Println(cmd.Args)
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

func applicationName(entry *desktop.Entry) string {

	generic := ""

	if entry.GenericName != "" {
		generic = fmt.Sprintf("(%s)", entry.GenericName)
	}

	return fmt.Sprintf("%s %s", entry.Name, generic)

}

func launchApplication(execString string) error {
	execString2 := strings.Fields(execString)

	cmd := exec.Command(execString2[0], execString2[1:]...)

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
	Name   string
	Key    string
}

func findModes(configFile string) ([]*Mode, error) {
	modes := make([]*Mode, 0)

	dat, err := os.ReadFile(configFile)
	if err != nil {
		return modes, err
	}

	cfg := make(map[string]map[string]Mode)

	err = toml.Unmarshal(dat, &cfg)
	if err != nil {
		return modes, err
	}
	fmt.Println(err)

	for k := range cfg["modes"] {
		mode := cfg["modes"][k]
		mode.Key = k
		modes = append(modes, &mode)
	}

	return modes, nil
}

func main() {
	if len(os.Args[1:]) > 0 {
		fmt.Println("d\ne")
	} else {
		fzf([]byte("a\nb"))
	}
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
