package sh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Entry struct {
	ModeKey string
	Id      string
	Text    string
}

func SpawnAsyncProcess(command string) error {
	execList := strings.Fields(command)

	cmd := exec.Command(execList[0], execList[1:]...)

	err := cmd.Start()

	return err
}

func SpawnSyncProcess(command []string, input []byte) (string, error) {
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

func FormatEntries(entries []Entry) string {
	var sb strings.Builder

	for _, e := range entries {
		s := fmt.Sprintf("%s\034%s\034%s\n", e.ModeKey, e.Id, e.Text)
		sb.WriteString(s)
	}

	return sb.String()
}

func Fzf(entries []Entry) (*Entry, error) {
	ownExe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	bind := fmt.Sprintf("change:reload:%s search {q} || true", ownExe)
	command := []string{
		"fzf",
		"--nth=..3",
		"--with-nth=3",
		"--delimiter=\034",
		"--no-sort",
		"--ansi",
		"--extended",
		"--no-multi",
		"--cycle",
		"--no-info",
		"--bind",
		bind,
	}

	result, err := SpawnSyncProcess(command, []byte(FormatEntries(entries)))
	if err != nil {
		return nil, err
	}

	// TODO: Make \034 a delimiter
	separatedEntry := strings.Split(result, "\034")
	if len(separatedEntry) != 3 {
		return nil, fmt.Errorf("Result %s not compatible with delimiters", result)
	}
	return &Entry{ModeKey: separatedEntry[0], Id: separatedEntry[1], Text: separatedEntry[2]}, nil
}
