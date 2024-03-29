package launcher

import (
	"bytes"
	"sort"
	"strings"
	"time"

	"github.com/gjhenrique/yafl/sh"
	"github.com/gjhenrique/yafl/store"
)

type Mode struct {
	CacheTime        *int `toml:"cache"`
	Exec             string
	Prefix           string
	Key              string
	CallWithoutMatch bool `toml:"call_without_match"`
	HistoryEnabled   bool `toml:"history"`
}

const Delimiter byte = 31

func (m *Mode) ListEntries(historyStore *store.HistoryStore, cache *store.CacheStore) ([]*sh.Entry, error) {
	cmd := strings.Fields(m.Exec)

	duration := time.Duration(*m.CacheTime) * time.Second

	result, err := cache.FetchCache(m.Key, duration, func() ([]byte, error) {
		value, err := sh.SpawnSyncProcess(cmd, nil)
		if err != nil {
			return nil, err
		}
		return []byte(value), nil
	})
	if err != nil {
		return nil, err
	}

	entriesFromCmd := bytes.Split(result, []byte("\n"))

	entries := make([]*sh.Entry, len(entriesFromCmd))

	delimiter := []byte{Delimiter}
	for i, r := range entriesFromCmd {
		splittedEntry := bytes.Split(r, delimiter)

		text, id := r, r
		if len(r) == 0 {
			continue
		}

		if len(splittedEntry) > 1 {
			id = splittedEntry[0]
			text = bytes.Join(splittedEntry[1:], []byte(" "))
		}

		newText := text

		if m.Prefix != "" {
			prefix := []byte(m.Prefix)
			newText = make([]byte, len(prefix)+len(text))
			copy(newText, prefix)
			copy(newText[len(prefix):], text)
		}

		entries[i] = &sh.Entry{ModeKey: m.Key, Text: newText, Id: id}
	}

	if m.HistoryEnabled {
		historyEntries, err := historyStore.ListEntries(m.Key)
		if err != nil {
			return nil, err
		}

		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i] == nil {
				return true
			}
			posI, ok := historyEntries.FindPosition([]byte(entries[i].Id))
			if !ok {
				return false
			}

			if entries[j] == nil {
				return false
			}
			posJ, ok := historyEntries.FindPosition([]byte(entries[j].Id))
			if !ok {
				return true
			}

			return posI < posJ
		})

	}

	return entries, nil
}

func (m *Mode) Launch(input []byte) error {
	cmd := strings.Fields(m.Exec)

	err := sh.SpawnAsyncProcess(cmd, string(input))
	if err != nil {
		return err
	}

	return nil
}
