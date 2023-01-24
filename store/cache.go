package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	enc              = binary.BigEndian
	lenUnixMilli     = 8
	lenPayloadLength = 4
)

type CacheStore struct {
	Dir string
}

func populateCache(f *os.File, duration time.Duration, action func() (string, error)) (string, error) {
	t := time.Now().Add(duration)

	result, err := action()
	if err != nil {
		return result, err
	}

	tBytes := make([]byte, lenUnixMilli)
	enc.PutUint64(tBytes, uint64(t.UnixMilli()))
	f.WriteAt(tBytes, 0)

	lenBytes := make([]byte, lenPayloadLength)
	enc.PutUint32(lenBytes, uint32(len([]byte(result))))
	f.WriteAt(lenBytes, int64(lenUnixMilli))

	f.WriteAt([]byte(result), int64(lenUnixMilli+lenPayloadLength))

	return result, nil
}

func (c *CacheStore) FetchCache(key string, duration time.Duration, action func() (string, error)) (string, error) {
	fileName := filepath.Join(c.Dir, fmt.Sprintf("%s_cache", key))

	f, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, _ = os.Create(fileName)
		} else {
			// TODO: Log this action
			return action()
		}
	}

	b := make([]byte, lenUnixMilli)
	_, err = f.ReadAt(b, 0)
	if err != nil {
		return populateCache(f, duration, action)
	}

	tMilli := enc.Uint64(b)
	cacheT := time.Unix(int64(tMilli/1000), 0)
	if cacheT.Before(time.Now()) {
		return populateCache(f, duration, action)
	}

	lenBytes := make([]byte, lenPayloadLength)
	_, err = f.ReadAt(lenBytes, int64(lenUnixMilli))
	if err != nil {
		return populateCache(f, duration, action)
	}
	len := enc.Uint32(lenBytes)
	strBytes := make([]byte, len)
	_, err = f.ReadAt(strBytes, int64(lenUnixMilli+lenPayloadLength))

	return string(strBytes), nil
}

func (c *CacheStore) Remove(key string) error {
	path := filepath.Join(c.Dir, fmt.Sprintf("%s_cache", key))
	// Vulnerable to file traversal attack. It doesn't matter because we don't expose this tool to the Internet anyway
	return os.Remove(path)
}
