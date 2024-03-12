package sh

import (
	"bufio"
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
	Id      []byte
	Text    []byte
}

const Delimiter = "\\x31"

func SpawnAsyncProcess(command []string, options string) error {
	args := command[1:]

	// Removing double quotes for now. It's having problems when execing the binary
	for i, arg := range args {
		args[i] = strings.Trim(arg, "\"")
	}

	if len(options) > 0 {
		args = append(args, string(options))
	}

	cmd := exec.Command(command[0], args...)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:     true,
		Foreground: false,
	}

	err := cmd.Start()

	return err
}

func SpawnSyncProcess(command []string, input []byte) ([]byte, error) {
	var b bytes.Buffer

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = bufio.NewWriter(&b)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	io.Copy(stdin, bytes.NewReader(input))
	if err = stdin.Close(); err != nil {
		return nil, err
	}

	err = cmd.Run()

	return b.Bytes(), err
}

func FormatEntries(entries []*Entry) string {
	var sb strings.Builder

	for _, e := range entries {
		if e == nil {
			continue
		}

		s := fmt.Sprintf("%s\034%s\034%s\n", e.ModeKey, e.Id, e.Text)
		sb.WriteString(s)
	}

	return sb.String()
}

type NoMatchError struct {
	Query []byte
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

	splittedResult := bytes.Split(result, []byte("\n"))
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
	separatedEntry := bytes.Split(splittedResult[1], []byte("\034"))
	if len(separatedEntry) != 3 {
		return nil, fmt.Errorf("Result %s not compatible with delimiters", result)
	}
	return &Entry{ModeKey: string(separatedEntry[0]), Id: separatedEntry[1], Text: separatedEntry[2]}, nil
}
