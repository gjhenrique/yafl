package cmd

import (
	// "example/cmd"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
// assert that removes %F
func TestFunction(t *testing.T) {
	// t.Parallel()

	// Write temporary file
	// Remove temporary directory
	// Set directory
	// Set cache dir

	tempDir, err := ioutil.TempDir("", "yafl")
	defer os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	name := uniuri.New()
	server := new(gotmux.Server)
	session, err := server.NewSession(name)
	defer server.KillSession(name)
	if err != nil {
		t.Fatal(err)
	}

	panes, err := session.ListPanes()
	if len(panes) == 0 || err != nil {
		t.Fatal(err)
	}

	configFile := filepath.Join(tempDir, "config.toml")
	cacheDir := filepath.Join(tempDir, "cache")
	os.Mkdir(filepath.Join(tempDir, "apps"), 0755)
	os.Mkdir(cacheDir, 0755)
	exe := fmt.Sprintf("../yafl --config=%s --cache-dir=%s", configFile, tempDir)

	s := []string{
		"send-keys",
		"-t", fmt.Sprintf("%s:0", name),
		exe,
	}
	gotmux.RunCmd(s)

	time.Sleep(1 * time.Second)
	content, err := panes[0].Capture()
	fmt.Println(content)
}
