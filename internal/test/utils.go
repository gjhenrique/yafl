package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func CheckTextInFile(t *testing.T, fileName string, matchText string) error {
	retry := 0

	for {
		content, err := os.ReadFile(fileName)
		require.NoError(t, err)
		if strings.Contains(string(content), matchText) {
			return nil
		}
		time.Sleep(10 * time.Millisecond)

		if retry > 20 {
			return fmt.Errorf("Text %s not found in %s", matchText, string(content))
		}

		retry += 1
	}
}

func (w *Workspace) RandomFile(t *testing.T) string {
	file, err := ioutil.TempFile(w.Dir, "random")
	require.NoError(t, err)
	return file.Name()
}

func (w *Workspace) RemoveWorkspace() {
	os.RemoveAll(w.Dir)
}
