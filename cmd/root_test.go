package cmd

import (
	// "example/cmd"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	// "time"

	"github.com/dchest/uniuri"
	gotmux "github.com/jubnzv/go-tmux"
)

// Test applications showing up

// Make new session.

// Clicking

// call_without_match
// assert icon
// assert Name appears correctly
// assert that application desktop was called
// test that cache (mock time somehow)
// assert that prefix is done correctly
// assert that prefix with space is done correctly
func TestFunction(t *testing.T) {
	// t.Parallel()

	// Write temporary file
	// Remove temporary directory
	// Set directory
	// Set cache dir

	name := uniuri.New()
	server := new(gotmux.Server)
	session, err := server.NewSession(name)
	defer server.KillSession(name)
	if err != nil {
		t.Fatal(err)
	}

	// time.Sleep(5 * time.Second)
	panes, err := session.ListPanes()
	if len(panes) == 0 || err != nil {
		t.Fatal(err)
	}

	s := []string{
		"send-keys",
		"-t", fmt.Sprintf("%s:0", name),
		"../yafl", "Enter",
	}
	gotmux.RunCmd(s)

	time.Sleep(1 * time.Second)
	content, err := panes[0].Capture()
	fmt.Println(content)
}

func TestTemp(t *testing.T) {
	file, err := ioutil.TempDir("", "yafl")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(file)
}
