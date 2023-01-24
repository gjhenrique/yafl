# yafl

A fzf launcher inspired by [sway-fzf-launcher] philosophy and Alfred[alfred] workflows

## Inspiration

I use a launcher not only to select a desktop application, but to additionaly provide some other unrelated tasks:
- choosing a Firefox bookmarks based on the title
- search words in different websites like Github or DuckDuckGo
- pick an emoji based on some tags

I was using [rofi] with `combi-mode` with the bang hack. But some things were missing:

- It's running in XWayland mode. So, the clipboard doesn't work when copying text from Emacs (another XWayland application). Kinda annoying to copy a code snipped in the text editor and not being able to paste it in the launcher.
- It (rarely) freezes with sway. `rofi` grabs the keyboard and mouse, so I need to go to TTY3 and kill the process manually. It doesn't happen often, though.
- The search could be more powerful

rofi indeed rocks, but I saw the [sway-fzf-launcher] I .
It delegates the user input and search to [fzf] and has only two responsibilities:
1. List desktop applications and binaires
1. `exec` the selected entry

Additionally, it's a simple cli app, so it executes in any terminal. 
The good part is there is no XWayland involved (yay clipboard).

<!-- Portable -->
<!-- Extensible -->
<!-- Cache will be there after one minute -->

<!-- Fzf for searching -->
<!-- Shortcuts from the terminal, not a custom application -->

## Why not extend it?

sway-fzf-launcher is clever, but has some downsides.
It's written in awk/bash and supporting modes would be painful to support (at least for my level of experience).
I was also trying to find a reason to write something in Go, so I created `yafl`.

Being a Go binary, `yafl` is more portable than bash scripts (hello Windows, but not yet. Right now, it's only Linux).

## But what's a mode?

`yafl` activates a mode whenever the `prefix` option matches the search term.
When that happens, `yafl` invokes the command specified by `exec` option.
This script may print some entries delimited by the end-of-line character (`\n`) in stdout.
When `fzf` returns the entry, we execute. This "API" is somewhat compatible with [rofi script mode][rofi-script]from rofi. Pretty ingenous.

``` toml
# Contents from $HOME/.config/yafl/config.toml
[modes.bookmark]
cache = 60
exec = "/bin/bookmarks"
prefix = "f"
name = "Search with bookmark"
# It stores the history
history = true
```

<!-- Now type f and <space> and `yafl` invokes the script. It's heavily based on [rofi script mode][rofi-script], i.e., -->
<!-- apps returns entries in a line-delimiter. -->

Additionally, if you have duplicated entries, you may delimit the key and value by `\x1f`, like `key\x1fDisplayValue`.
When you select the entry, `yafl` invokes the script passing the `key`.

More examples are in this directory [examples]

<!-- `yafl` uses this option because it needs to check which mode will invoke -->

To check which mode _should be activated_ based on the input, it uses the `fzf` `change:reload` option, so `fzf` invokes `yafl` again **for every keystroke pressed**.

<!-- ## How does it work? -->

<!-- When you invoke `yafl` without any options, it: -->
<!-- 1. Reads all `.desktop` files -->
<!-- 1. Sends them to fzf -->
<!-- 1. You can use `fzf` to select the application you want -->
<!-- 1. If a prefix matches the mode, `yafl` invokes the scripts and populates the list with it instead -->
<!-- 1. After the entry is selected, the mode script is invoked with the key of the entry or the entry itself. -->


## Yikes. Isn't it slow?

Yes!

Invoking the script every time brings some overhead.
To solve that, by default, `yafl` caches the entries of the mode for 60 seconds after it's invoked, so the cache instead of executing to list the entries.
I couldn't notice any input lags, but YMMV.

Cache invalidation needs to be manually called with:

``` shell
yafl cache clean mode_key
```

## Roadmap
- [ ] History support
- [ ] Real integration tests with tmux
- [ ] CI to run specs
- [ ] Handle errors properly instead of "panic"-ing
- [ ] Supports binaries

## History Support

[fzf]: https://github.com/junegunn/fzf
[sway-fzf-launcher]: https://github.com/Biont/sway-launcher-desktop
[alfred]: https://www.alfredapp.com/
[rofi]: https://github.com/davatorium/rofi
[rofi-script]: https://davatorium.github.io/rofi/current/rofi-script.5/
[rofi-bang]: https://github.com/davatorium/rofi/blob/8155b2c476694e55452b14cbd15058d85df095db/doc/rofi.1.markdown?plain=1#L909-L911
