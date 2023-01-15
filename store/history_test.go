package store

import (
	"testing"

	"github.com/gjhenrique/yafl/internal/test"
	"github.com/stretchr/testify/require"
)

func TestOrderIsMaintained(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := HistoryStore{Dir: workspace.CacheDir}

	err := store.IncrementEntry("key", []byte("a"))
	require.NoError(t, err)
	err = store.IncrementEntry("key", []byte("b"))
	require.NoError(t, err)
	err = store.IncrementEntry("key", []byte("c"))
	require.NoError(t, err)
	err = store.IncrementEntry("key", []byte("b"))
	require.NoError(t, err)

	entries, err := store.ListEntries("key")
	require.NoError(t, err)

	require.Len(t, entries, 3)
	require.Equal(t, string(entries[0]), "b")
	require.Equal(t, string(entries[1]), "a")
	require.Equal(t, string(entries[2]), "c")
}

func TestEmptyArrayWhenFileDoesNotExist(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := HistoryStore{Dir: workspace.CacheDir}

	entries, err := store.ListEntries("key")
	require.NoError(t, err)
	require.Len(t, entries, 0)
}
