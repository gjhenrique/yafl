package launcher

import (
	"strings"
	"testing"

	"github.com/gjhenrique/yafl/internal/test"
	"github.com/stretchr/testify/require"
)

func TestNonExistentFile(t *testing.T) {
	modes, err := ParseModesFromConfig("/tmp/nonsense.txt")
	require.NoError(t, err)
	require.Len(t, modes, 1)
	assertAppMode(t, modes[0])
}

func TestBlankFile(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.WriteScript(t, "")

	modes, err := ParseModesFromConfig(randomFile)
	require.NoError(t, err)
	require.Len(t, modes, 1)
	assertAppMode(t, modes[0])
}

func TestCompleteMode(t *testing.T) {
	var mode = `
[modes.bookmark]
exec="ls -lah"
prefix="s"
cache=10
call_without_match=true`

	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.WriteScript(t, mode)

	modes, err := ParseModesFromConfig(randomFile)
	require.NoError(t, err)
	require.Len(t, modes, 2)
	assertAppMode(t, modes[1])
	require.Equal(t, modes[0].Key, "bookmark")
	require.Equal(t, modes[0].Prefix, "s ")
	require.Equal(t, modes[0].Key, "bookmark")
	require.Equal(t, modes[0].Exec, "ls -lah")
	require.Equal(t, *modes[0].CacheTime, 10)
	require.True(t, modes[0].CallWithoutMatch)
}

func TestModeWithDefaultValues(t *testing.T) {
	var mode = `
[modes.bookmark]
exec="ls -lah"`

	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.WriteScript(t, mode)

	modes, err := ParseModesFromConfig(randomFile)
	require.NoError(t, err)
	require.Len(t, modes, 2)
	assertAppMode(t, modes[1])
	require.Equal(t, modes[0].Key, "bookmark")
	require.Equal(t, modes[0].Prefix, "")
	require.Equal(t, modes[0].Key, "bookmark")
	require.Equal(t, modes[0].Exec, "ls -lah")
	require.Equal(t, *modes[0].CacheTime, 60)
	require.False(t, modes[0].CallWithoutMatch)
	require.False(t, modes[0].HistoryEnabled)
}

func TestDoesNotOverrideExecButCacheMode(t *testing.T) {
	var mode = `
[modes.apps]
exec="ls -lah"
cache=30`

	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	randomFile := workspace.WriteScript(t, mode)

	modes, err := ParseModesFromConfig(randomFile)
	require.NoError(t, err)
	require.Len(t, modes, 1)
	require.True(t, strings.HasSuffix(modes[0].Exec, "test apps"))
	require.Equal(t, *modes[0].CacheTime, 30)
}

func assertAppMode(t *testing.T, mode *Mode) {
	require.Equal(t, mode.Key, "apps")
	require.Equal(t, *mode.CacheTime, *&defaultCacheTime)
	require.True(t, strings.HasSuffix(mode.Exec, " apps"))
	require.Equal(t, mode.Prefix, "")
	require.False(t, mode.CallWithoutMatch)
	require.True(t, mode.HistoryEnabled)
}
