package sh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Entry struct {
	ModeKey string
	Id      string
	Text    string
}

const Delimiter = "\\x31"

func SpawnAsyncProcess(command []string, options string) error {
	args := command[1:]

	if options != "" {
		args = append(command[1:], options)
	}

	cmd := exec.Command(command[0], args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:     true,
		Foreground: false,
	}

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

func FormatEntries(entries []*Entry) string {
	var sb strings.Builder

	for _, e := range entries {
		s := fmt.Sprintf("%s\034%s\034%s\n", e.ModeKey, e.Id, e.Text)
		sb.WriteString(s)
	}

	return sb.String()
}

type NoMatchError struct {
	Query string
}

func (e *NoMatchError) Error() string {
	return fmt.Sprintf("Query %s not found", e.Query)
}

type SkippedInputError struct{}

func (e *SkippedInputError) Error() string {
	return "User selected nothing"
}

func Fzf(entries []*Entry) (*Entry, error) {
	ownExe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	bind := fmt.Sprintf("change:reload:%s search {q} || true", ownExe)
	command := []string{
		"fzf",
		"--exact",
		"--no-info",
		"--nth=..3",
		"--with-nth=3",
		"--delimiter=\034",
		"--color=16,gutter:-1",
		"--margin=1,2",
		"--no-sort",
		"--ansi",
		"--header=",
		"--extended",
		"--no-multi",
		"--cycle",
		"--no-info",
		"--print-query",
		"--bind",
		bind,
	}

	result, err := SpawnSyncProcess(command, []byte(FormatEntries(entries)))

	splittedResult := strings.Split(result, "\n")
	query := splittedResult[0]

	if err != nil {
		// https://www.mankier.com/1/fzf#Exit_Status
		if exitError, ok := err.(*exec.ExitError); ok {
			// TODO: Return no match error
			// TODO: Treat 130 case - User pressed key
			if exitError.ExitCode() == 1 {
				// TODO: Handle this

				return nil, &NoMatchError{Query: query}
			} else if exitError.ExitCode() == 130 {
				return nil, &SkippedInputError{}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// TODO: Make \034 a variable
	// Use same delimiter from rofi
	separatedEntry := strings.Split(splittedResult[1], "\034")
	if len(separatedEntry) != 3 {
		return nil, fmt.Errorf("Result %s not compatible with delimiters", result)
	}
	return &Entry{ModeKey: separatedEntry[0], Id: separatedEntry[1], Text: separatedEntry[2]}, nil
}
