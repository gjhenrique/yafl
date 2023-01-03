package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type Workspace struct {
	Id       string
	Dir      string
	CacheDir string
}

func SetupWorkspace(t *testing.T) *Workspace {
	tempDir, err := ioutil.TempDir("", "yafl")
	require.NoError(t, err)

	cacheDir := filepath.Join(tempDir, "cache")
	err = os.Mkdir(cacheDir, 0755)
	require.NoError(t, err)

	return &Workspace{
		Dir:      tempDir,
		CacheDir: cacheDir,
	}
}

func (w *Workspace) WriteScript(t *testing.T, script string) string {
	file, err := ioutil.TempFile(w.Dir, "script")
	os.Chmod(file.Name(), 0755)
	os.WriteFile(file.Name(), []byte(script), 0755)
	require.NoError(t, err)
	return file.Name()
}

// Create workspace method that creates a temporary file populating its stuff

func (w *Workspace) RemoveWorkspace() {
	os.RemoveAll(w.Dir)
}
