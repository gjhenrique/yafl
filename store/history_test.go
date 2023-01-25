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

	pos, ok := entries.FindPosition([]byte("b"))
	require.True(t, ok)
	require.Equal(t, pos, 0)

	pos, ok = entries.FindPosition([]byte("a"))
	require.True(t, ok)
	require.Equal(t, pos, 1)

	pos, ok = entries.FindPosition([]byte("c"))
	require.True(t, ok)
	require.Equal(t, pos, 2)

	_, ok = entries.FindPosition([]byte("d"))
	require.False(t, ok)
}

func TestEmptyArrayWhenFileDoesNotExist(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := HistoryStore{Dir: workspace.CacheDir}

	entries, err := store.ListEntries("key")
	require.NoError(t, err)

	_, ok := entries.FindPosition([]byte("b"))
	require.False(t, ok)
}
