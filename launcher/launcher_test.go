package launcher

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gjhenrique/yafl/internal/test"
	"github.com/gjhenrique/yafl/sh"
	"github.com/gjhenrique/yafl/store"
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
	prefixMode := mockMode(script2, "f ", "prefix")
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

func TestLaunch(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.RandomFile(t)
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en $1 > %s", randomFile))

	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries := []*sh.Entry{{ModeKey: "test", Id: "abc", Text: "abc"}}
	err := l.ProcessEntries(entries)
	require.NoError(t, err)

	err = test.CheckTextInFile(t, randomFile, "abc")
	require.NoError(t, err)
}

func TestLaunchWithNoMatch(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.RandomFile(t)
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en $1 > %s", randomFile))

	modes := mockMode(script, "", "test")
	modes[0].CallWithoutMatch = true
	l := createLauncher(modes, workspace.CacheDir, workspace)
	l.searcher = func([]*sh.Entry) (*sh.Entry, error) {
		return nil, &sh.NoMatchError{Query: "def"}
	}

	entries := []*sh.Entry{{ModeKey: "test", Id: "abc", Text: "abc"}}
	err := l.ProcessEntries(entries)
	require.NoError(t, err)

	err = test.CheckTextInFile(t, randomFile, "def")
	require.NoError(t, err)
}

func TestErrorWithNoMatch(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.RandomFile(t)
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en $1 > %s", randomFile))

	modes := mockMode(script, "", "test")
	modes[0].CallWithoutMatch = false
	l := createLauncher(modes, workspace.CacheDir, workspace)
	l.searcher = func([]*sh.Entry) (*sh.Entry, error) {
		return nil, &sh.NoMatchError{Query: "def"}
	}

	entries := []*sh.Entry{{ModeKey: "test", Id: "abc", Text: "abc"}}
	err := l.ProcessEntries(entries)
	require.EqualError(t, err, "Query def not found")
}

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
	if len(entries) == 0 {
		return nil, errors.New("Entry array can't be nil")
	}
	return entries[0], nil
}

func createLauncher(modes []*Mode, cacheDir string, workspace *test.Workspace) *Launcher {
	return &Launcher{
		cache:    &store.CacheStore{Dir: workspace.Dir},
		modes:    modes,
		searcher: defaultSearcher,
	}
}
