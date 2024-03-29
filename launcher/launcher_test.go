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

func TestListEntriesWithRawMode(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	script := workspace.WriteScript(t, "echo -en \"abc\\ndef\"")
	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries, err := l.ListEntries([]byte(""))
	require.NoError(t, err)

	require.Len(t, entries, 2)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: []byte("abc"), Text: []byte("abc")})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: []byte("def"), Text: []byte("def")})
}

func TestListEntriesWithDelimiter(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	// We need three inverted backlashes in echo
	// \x1f writes the Hex 1f/Decimal 31 into the string
	script := workspace.WriteScript(t, "echo -en \"abc\x1fdef\\nghi\x1fjkl\"")
	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries, err := l.ListEntries([]byte(""))
	require.NoError(t, err)

	require.Len(t, entries, 2)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: []byte("abc"), Text: []byte("def")})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: []byte("ghi"), Text: []byte("jkl")})
}

func TestErrorWithNoMode(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	var modes []*Mode
	l := createLauncher(modes, workspace.CacheDir, workspace)

	_, err := l.ListEntries([]byte(""))
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

	entries, err := l.ListEntries([]byte(""))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "root", Id: []byte("abc"), Text: []byte("abc")})

	prefixEntries, err := l.ListEntries([]byte("f def"))
	require.NoError(t, err)
	require.Len(t, prefixEntries, 1)
	require.Equal(t, *prefixEntries[0], sh.Entry{ModeKey: "prefix", Id: []byte("def"), Text: []byte("f def")})
}

func TestLaunch(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.RandomFile(t)
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en $1 > %s", randomFile))

	l := createLauncher(mockMode(script, "", "test"), workspace.CacheDir, workspace)

	entries := []*sh.Entry{{ModeKey: "test", Id: []byte("abc"), Text: []byte("abc")}}
	err := l.ProcessEntries(entries)
	require.NoError(t, err)

	err = test.CheckTextInFile(t, randomFile, "abc")
	require.NoError(t, err)
}

func TestLaunchWithHistory(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.RandomFile(t)
	script := workspace.WriteScript(t, fmt.Sprintf("echo -en $1 > %s", randomFile))

	m := mockMode(script, "", "test")
	m[0].HistoryEnabled = true
	l := createLauncher(m, workspace.CacheDir, workspace)

	entries := []*sh.Entry{{ModeKey: "test", Id: []byte("b"), Text: []byte("b")}}
	err := l.ProcessEntries(entries)
	require.NoError(t, err)

	script2 := workspace.WriteScript(t, "echo -en \"a\\nb\\nc\\nd\"")
	m[0].Exec = fmt.Sprintf("bash %s", script2)

	entries, err = l.ListEntries([]byte(""))
	require.NoError(t, err)

	require.Len(t, entries, 4)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: []byte("b"), Text: []byte("b")})
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
		return nil, &sh.NoMatchError{Query: []byte("def")}
	}

	entries := []*sh.Entry{{ModeKey: "test", Id: []byte("abc"), Text: []byte("abc")}}
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
		return nil, &sh.NoMatchError{Query: []byte("def")}
	}

	entries := []*sh.Entry{{ModeKey: "test", Id: []byte("abc"), Text: []byte("abc")}}
	err := l.ProcessEntries(entries)
	require.EqualError(t, err, "Query def not found")
}

func TestHistoryListSort(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	script := workspace.WriteScript(t, "echo -en \"a\\nb\\nc\\nd\"")
	m := mockMode(script, "", "test")
	m[0].HistoryEnabled = true
	l := createLauncher(m, workspace.CacheDir, workspace)

	l.historyStore.IncrementEntry("test", []byte("b"))
	l.historyStore.IncrementEntry("test", []byte("b"))
	l.historyStore.IncrementEntry("test", []byte("c"))

	entries, err := l.ListEntries([]byte(""))
	require.NoError(t, err)

	require.Len(t, entries, 4)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: []byte("b"), Text: []byte("b")})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: []byte("c"), Text: []byte("c")})
	require.Equal(t, *entries[2], sh.Entry{ModeKey: "test", Id: []byte("a"), Text: []byte("a")})
	require.Equal(t, *entries[3], sh.Entry{ModeKey: "test", Id: []byte("d"), Text: []byte("d")})
}

func TestHistoryDisabledNoImpact(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	script := workspace.WriteScript(t, "echo -en \"a\\nb\\nc\\nd\"")
	m := mockMode(script, "", "test")
	l := createLauncher(m, workspace.CacheDir, workspace)

	l.historyStore.IncrementEntry("test", []byte("b"))
	l.historyStore.IncrementEntry("test", []byte("b"))
	l.historyStore.IncrementEntry("test", []byte("c"))

	entries, err := l.ListEntries([]byte(""))
	require.NoError(t, err)

	require.Len(t, entries, 4)
	require.Equal(t, *entries[0], sh.Entry{ModeKey: "test", Id: []byte("a"), Text: []byte("a")})
	require.Equal(t, *entries[1], sh.Entry{ModeKey: "test", Id: []byte("b"), Text: []byte("b")})
	require.Equal(t, *entries[2], sh.Entry{ModeKey: "test", Id: []byte("c"), Text: []byte("c")})
	require.Equal(t, *entries[3], sh.Entry{ModeKey: "test", Id: []byte("d"), Text: []byte("d")})

}

func mockMode(scriptName, prefix, key string) []*Mode {
	noCache := 0

	return []*Mode{
		{
			CacheTime: &noCache,
			Exec:      fmt.Sprintf("bash %s", scriptName),
			Prefix:    prefix,
			Key:       key,
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
		cacheStore:   &store.CacheStore{Dir: workspace.Dir},
		modes:        modes,
		searcher:     defaultSearcher,
		historyStore: &store.HistoryStore{Dir: workspace.Dir},
	}
}
