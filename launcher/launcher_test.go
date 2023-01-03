package launcher

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gjhenrique/yafl/cache"
	"github.com/gjhenrique/yafl/internal/test"
	"github.com/gjhenrique/yafl/sh"
	"github.com/stretchr/testify/require"
)

// How to mock FZF? Inject it

func TestListEntriesWithRawMode(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	script := workspace.WriteScript(t, "echo -en \"abc\\ndef\"")
	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries, err := l.ListEntries("")
	require.NoError(t, err)

	require.Len(t, entries, 2)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: "abc", Text: "abc"})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: "def", Text: "def"})
}

func TestListEntriesWithDelimiter(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	// We need three inverted backlashes in echo
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en \"abc\\\\\\%sdef\\nghi\\\\\\%sjkl\"", sh.Delimiter, sh.Delimiter))
	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries, err := l.ListEntries("")
	require.NoError(t, err)

	require.Len(t, entries, 2)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: "abc", Text: "def"})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: "ghi", Text: "jkl"})
}

func TestErrorWithNoMode(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	var modes []*Mode
	l := createLauncher(modes, workspace.CacheDir, workspace)

	_, err := l.ListEntries("")
	require.Error(t, err)
}

func TestWithMultiplePrefix(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	script1 := workspace.WriteScript(t, "echo -en \"abc\"")
	script2 := workspace.WriteScript(t, "echo -en \"def\"")
	rootMode := mockMode(script1, "", "root")
	prefixMode := mockMode(script2, "f", "prefix")
	modes := append(append([]*Mode{}, rootMode...), prefixMode...)
	l := createLauncher(modes, workspace.CacheDir, workspace)

	entries, err := l.ListEntries("")
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "root", Id: "abc", Text: "abc"})

	prefixEntries, err := l.ListEntries("f def")
	require.NoError(t, err)
	require.Len(t, prefixEntries, 1)
	require.Equal(t, *prefixEntries[0], sh.Entry{ModeKey: "prefix", Id: "def", Text: "f def"})
}

// TODO: One test that checks prefix order

func mockMode(scriptName, prefix, key string) []*Mode {
	noCache := 0

	return []*Mode{
		{
			Cache:  &noCache,
			Exec:   fmt.Sprintf("bash %s", scriptName),
			Prefix: prefix,
			Key:    key,
		},
	}
}

func defaultSearcher(entries []*sh.Entry) (*sh.Entry, error) {
	entry := entries[0]
	if entry == nil {
		return nil, errors.New("Entry array can't be nil")
	}
	return entry, nil
}

func createLauncher(modes []*Mode, cacheDir string, workspace *test.Workspace) *Launcher {
	return &Launcher{
		cache:    &cache.CacheStore{Dir: workspace.Dir},
		modes:    modes,
		searcher: defaultSearcher,
	}
}
