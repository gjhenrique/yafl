package sh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func CommandDeattached(command string) error {
	execList := strings.Fields(command)

	cmd := exec.Command(execList[0], execList[1:]...)

	err := cmd.Start()

	return err
}

func CommandWithOutput(command []string, input []byte) (string, error) {
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

func Fzf(input []byte) (string, error) {
	ownExe, err := os.Executable()
	if err != nil {
		return "", err
	}
	bind := fmt.Sprintf("change:reload:sleep 0.05; %s search {q} || true", ownExe)
	command := []string{"fzf", "--ansi", "--sort", "--extended", "--no-multi", "--cycle", "--no-info", "--bind", bind}

	return CommandWithOutput(command, input)
}
