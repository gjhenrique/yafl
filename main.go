package main

import (
	"github.com/gjhenrique/lfzf/cmd"
)

func main() {
	// if len(os.Args[1:]) > 0 {
	// 	fmt.Println(modes[0])
	// } else {
	// 	fzf([]byte("a\nb"))
	// }

	// modes, _ := findModes(filepath.Join(appFolder(), "config.toml"))
	// mode := FindMode("f ola", modes)
	// mode.Launch("f ola ola ola")

	cmd.Execute()

	// appFolder()
	// entries, err := GetDesktopEntries()

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// names := applicationNames(entries)

	// appName, _ := fzf([]byte(names))
	// appName = strings.TrimSuffix(appName, "\n")

	// entry, _ := GetEntryFromName(appName)

	// launchApplication(entry.Exec)
}
