package apps

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"code.rocketnine.space/tslocum/desktop"
)

type DesktopEntry struct {
	Name string
	Exec string
}

func getDesktopEntries() ([]*desktop.Entry, error) {
	allEntries := make([]*desktop.Entry, 0, 100)

	dirs := desktop.DataDirs()
	entries, err := desktop.Scan(dirs)

	if err != nil {
		return nil, err
	}

	for _, dir := range entries {
		for _, entry := range dir {
			// Getting the regex from rofi
			r := regexp.MustCompile(`%[fFuU]`)
			entry.Exec = r.ReplaceAllString(entry.Exec, "")

			allEntries = append(allEntries, entry)
		}
	}

	return allEntries, nil
}

func applicationName(entry *desktop.Entry) string {
	generic := ""

	if entry.GenericName != "" {
		generic = fmt.Sprintf("(%s)", entry.GenericName)
	}

	return fmt.Sprintf("%s %s", entry.Name, generic)

}

func applicationNames(entries []*desktop.Entry) string {
	names := make([]string, len(entries))

	for i, entry := range entries {
		names[i] = applicationName(entry)
	}

	return strings.Join(names, "\n")

}

func FormattedApplicationNames() (string, error) {
	entries, err := getDesktopEntries()
	if err != nil {
		return "", err
	}

	return applicationNames(entries), nil
}

func GetEntryFromName(chosenApp string) (*desktop.Entry, error) {
	entries, err := getDesktopEntries()

	if err != nil {
		return nil, err
	}

	var entry *desktop.Entry

	for _, e := range entries {
		if applicationName(e) == chosenApp {
			entry = e
		}
	}

	if entry == nil {
		return nil, errors.New("Didn't find any app name")
	}

	return entry, nil
}
