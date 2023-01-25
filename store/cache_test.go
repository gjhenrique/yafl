package store

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gjhenrique/yafl/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReturnCorrectValues(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key", 60*time.Second, cb)
	value, err := store.FetchCache("key", 60*time.Second, cb)
	require.NoError(t, err)
	require.Equal(t, string(value), "1")
}

func TestActionIsNotCalledMultipleTimes(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key", 60*time.Second, cb)
	store.FetchCache("key", 60*time.Second, cb)
	store.FetchCache("key", 60*time.Second, cb)
	require.Equal(t, i, 1)
}

func TestCacheIsNotPopulatedWhenItReturnsAnError(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), errors.New("Code Error")
	}

	store.FetchCache("key", 60*time.Second, cb)
	value, err := store.FetchCache("key", 60*time.Second, cb)
	require.Error(t, err, "Code Error")
	require.Equal(t, string(value), "2")
}

func TestErrorWhenDirectoryIsNotThere(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: filepath.Join(workspace.Dir, "not_a_folder")}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key", 60*time.Second, cb)
	store.FetchCache("key", 60*time.Second, cb)
	store.FetchCache("key", 60*time.Second, cb)
	require.Equal(t, i, 3)
}

func TestDifferentKeysNotInterferingWithEachOther(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key1", 60*time.Second, cb)
	store.FetchCache("key2", 60*time.Second, cb)
	require.Equal(t, i, 2)
}

func TestInvalidatesCacheProperly(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key", 0, cb)
	// Invokes it because there is no expiration time
	store.FetchCache("key", 10*time.Second, cb)
	store.FetchCache("key", 10*time.Second, cb)
	require.Equal(t, i, 2)
}

func TestInvokesAgainWhenCacheIsRemoved(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	store := CacheStore{Dir: workspace.CacheDir}
	i := 0
	cb := func() ([]byte, error) {
		i += 1
		return []byte(strconv.Itoa(i)), nil
	}

	store.FetchCache("key", 10*time.Second, cb)
	store.FetchCache("key", 10*time.Second, cb)
	err := store.Remove("key")
	require.NoError(t, err)

	store.FetchCache("key", 10*time.Second, cb)
	require.Equal(t, i, 2)
}
