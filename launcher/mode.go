package launcher

import (
	"strings"
	"time"

	"github.com/gjhenrique/yafl/sh"
	"github.com/gjhenrique/yafl/store"
)

type Mode struct {
	Cache            *int
	Exec             string
	Prefix           string
	Key              string
	CallWithoutMatch bool `toml:"call_without_match"`
}

func (m *Mode) ListEntries(cache *store.CacheStore) ([]*sh.Entry, error) {
	cmd := strings.Fields(m.Exec)

	duration := time.Duration(*m.Cache) * time.Second

	result, err := cache.FetchCache(m.Key, duration, func() (string, error) {
		return sh.SpawnSyncProcess(cmd, nil)
	})
	if err != nil {
		return nil, err
	}

	entriesFromCmd := strings.Split(result, "\n")

	entries := make([]*sh.Entry, len(entriesFromCmd))

	for i, r := range entriesFromCmd {
		splittedEntry := strings.Split(r, sh.Delimiter)

		text := r
		id := text
		if len(splittedEntry) > 1 {
			id = splittedEntry[0]
			text = strings.Join(splittedEntry[1:], " ")
		}

		if m.Prefix != "" {
			text = m.Prefix + text
		}

		entries[i] = &sh.Entry{ModeKey: m.Key, Text: text, Id: id}
	}

	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (m *Mode) Launch(input string) error {
	cmd := strings.Fields(m.Exec)

	err := sh.SpawnAsyncProcess(cmd, input)
	if err != nil {
		return err
	}

	return nil
}
