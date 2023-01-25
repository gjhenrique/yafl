package launcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gjhenrique/yafl/sh"
	"github.com/gjhenrique/yafl/store"
)

var (
	defaultCacheTime = 60
)

type Launcher struct {
	cacheStore   *store.CacheStore
	historyStore *store.HistoryStore
	modes        []*Mode
	searcher     func([]*sh.Entry) (*sh.Entry, error)
}

func NewLauncher(modes []*Mode, cacheFolder string, searcher func([]*sh.Entry) (*sh.Entry, error)) (*Launcher, error) {
	c := store.CacheStore{
		Dir: cacheFolder,
	}

	h := store.HistoryStore{
		Dir: cacheFolder,
	}

	return &Launcher{
		cacheStore:   &c,
		modes:        modes,
		searcher:     searcher,
		historyStore: &h,
	}, nil
}

func (l *Launcher) ListEntries(input string) ([]*sh.Entry, error) {
	selectedMode := l.findModeByInput(input)
	if selectedMode == nil {
		return nil, errors.New("Couldn't find any mode that matches this entry")
	}
	return selectedMode.ListEntries(l.historyStore, l.cacheStore)
}

func (l *Launcher) ProcessEntries(entries []*sh.Entry) error {
	var m *Mode

	entry, err := l.searcher(entries)

	if err != nil {
		if noMatchErr, ok := err.(*sh.NoMatchError); ok {
			m := l.findModeByInput(noMatchErr.Query)
			if m.CallWithoutMatch {
				query := strings.TrimPrefix(noMatchErr.Query, m.Prefix)
				entry = &sh.Entry{ModeKey: m.Key, Id: query, Text: noMatchErr.Query}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	if m, err = l.findModeByKey(entry.ModeKey); err != nil {
		return err
	}

	if err = m.Launch(entry.Id); err != nil {
		return err
	}

	l.historyStore.IncrementEntry(m.Key, []byte(entry.Id))

	return nil
}

func (l *Launcher) findModeByKey(key string) (*Mode, error) {
	var mode *Mode

	for _, m := range l.modes {
		if m.Key == key {
			mode = m
		}
	}

	if mode == nil {
		return nil, fmt.Errorf("Mode %s not found", key)
	}

	return mode, nil
}

func (l *Launcher) findModeByInput(input string) *Mode {
	var mode *Mode

	for _, m := range l.modes {
		if m.Prefix == "" {
			continue
		}

		if strings.HasPrefix(input, m.Prefix) {
			mode = m
		}
	}

	if mode == nil {
		for _, m := range l.modes {
			// Get the first empty prefix
			// Add a multi-mode in the future
			if m.Prefix == "" {
				mode = m
			}
		}
	}

	return mode
}

func appMode(modes []*Mode) *Mode {
	var appMode *Mode

	for _, m := range modes {
		if m.Key == "apps" {
			appMode = m
		}
	}

	return appMode
}
