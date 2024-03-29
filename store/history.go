package store

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type HistoryStore struct {
	Dir string
}

type historyEntry struct {
	key      []byte
	position int
}

type HistoryEntries struct {
	entries map[string]historyEntry
}

func formatEntry(entryKey []byte, n int) []byte {
	e := append(entryKey, "\t"...)
	return append(e, strconv.Itoa(n)...)

}

func (h HistoryEntries) FindPosition(key []byte) (int, bool) {
	val, ok := h.entries[string(key)]
	if ok {
		return val.position, true
	}

	return 0, false
}

func (h *HistoryStore) ListEntries(modeKey string) (HistoryEntries, error) {
	fileContent, err := ioutil.ReadFile(h.fileName(modeKey))
	fields := bytes.Split(fileContent, []byte("\n"))
	s := HistoryEntries{
		entries: make(map[string]historyEntry, 0),
	}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		} else {
			return s, err
		}
	}

	var pos int
	for _, field := range fields {
		arr := bytes.Split(field, []byte("\t"))
		first := arr[0]
		if len(first) > 0 {
			entry := historyEntry{key: first, position: pos}
			s.entries[string(entry.key)] = entry
			pos += 1
		}
	}

	return s, nil
}

func (h *HistoryStore) fileName(modeKey string) string {
	return filepath.Join(h.Dir, fmt.Sprintf("%s_history", modeKey))
}

func (h *HistoryStore) IncrementEntry(modeKey string, entryKey []byte) error {
	fileContent, err := ioutil.ReadFile(h.fileName(modeKey))

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fileContent = make([]byte, 0)
		} else {
			return err
		}
	}

	fields := bytes.Split(fileContent, []byte("\n"))
	var found bool

	for i, field := range fields {
		arr := bytes.Split(field, []byte("\t"))
		first := arr[0]
		if bytes.Compare(first, []byte(entryKey)) == 0 {
			found = true

			last := arr[len(arr)-1]
			n, err := strconv.Atoi(string(last))

			if err != nil {
				n = 0
			}

			fields[i] = formatEntry(entryKey, n+1)
		}
	}

	if !found {
		entry := formatEntry(entryKey, 1)
		fields = append(fields, entry)
	}

	sort.SliceStable(fields, func(i, j int) bool {
		arr1 := bytes.Split(fields[i], []byte("\t"))
		a, _ := strconv.Atoi(string(arr1[len(arr1)-1]))

		arr2 := bytes.Split(fields[j], []byte("\t"))
		b, _ := strconv.Atoi(string(arr2[len(arr2)-1]))

		return a > b
	})

	finalFields := bytes.Join(fields, []byte("\n"))

	return ioutil.WriteFile(h.fileName(modeKey), finalFields, 0644)
}
