# yafl

yafl (yet another fzflauncher) uses [sway-fzf-launcher] ideas and [rofi] script mode.

![screen](https://s9.gifyu.com/images/recording3.gif)

## Why not rofi or sway-fzf-launcher?

Running rofi on Wayland can be annoying.
Rarely it crashes the entire desktop because it steals the keyboard and mouse.
Unfortunately, pasting into it from other XWayland apps also doesn't work.

[sway-fzf-launcher] is marvellous because it delegates the display part to a terminal emulator and search logic to fzf, but it has some downsides.
First, it's written in awk/bash and supporting modes would be painful to support (at least for my level of experience).

Above all, I used an excuse to write something in Go, so I created `yafl`.
Being a Go binary, `yafl` is easier to test and more portable than bash scripts (hello, Windows. But not yet. Right now, it's only Linux).

## What's a mode?

`yafl` returns the entries from a different mode whenever the `prefix` option matches the search term.
When that happens, `yafl` invokes the command specified by `exec` option.
The script should print the entries separated by `\n` in `stdout`.
`fzf` then displays the entries and you can select some of them. `yafl` will call the **same** cli, but passing the selected entry as an argument now. This "API" is somewhat compatible with [script mode][rofi-script]from rofi. Pretty clever stuff from them 👏!

Additionally, you may delimit the key and value by `\x1f`, like `echo -n key\x1fvalue`.
Then, when you select the entry from fzf, `yafl` invokes the script passing the `key` as the first argument.

More examples are in this `examples` directory:
- [yafl_bookmark]: Lists all the bookmarks of a Firefox installation and opens a new tab whenever you find them
- [yafl_search]: Search by multiple search engines
- [yafl_moji]: Wraps [rofimoji] and to list and select emojis

## How does it activate the modes?

To check which mode _should be activated_ based on the input, it uses the fzf `change:reload` option, i.e. `fzf` invokes `yafl` again **for every keystroke you press**.

Yikes. That's slow.

Invoking the mode script every time brings some overhead.
To solve that, by default, `yafl` caches the mode entries for 60 seconds after one invokes it.
I couldn't notice any input lags, but YMMV.

Cache invalidation shouldn't be an issue, but you can manually call with:

``` shell
yafl cache clean mode_key
```

## Running on sway

``` shell
for_window [app_id="^yafl"] floating enable, sticky enable, resize set 700 px 500 px, border pixel 10
set $menu exec $term --app-id=yafl -e /home/guilherme/Projects/mine/yafl/yafl
bindsym $mod+d exec $menu
```

## Modes API

Create a file in `$HOME/.config/yaml/config.toml` with the following options

``` toml
# bookmark is the mode key
[modes.bookmark]
# Cache time before yafl executes the script again
cache = 60
# Cli that returns entries searchable by fzf and executes the entry
# \x1f byte is the delimiter between the key and value, but it's not necessary
exec = "bookmarks"
# Input starting with "f <SPC>" activates the mode and shows only the entries of this mode
prefix = "f"
# Ranks the selected entries based on previous selections. Defaults to false
history = true
# Calls the exec cli with the input even if there are no matches. Defaults to false
call_without_match = true
```

[fzf]: https://github.com/junegunn/fzf
[sway-fzf-launcher]: https://github.com/Biont/sway-launcher-desktop
[rofi]: https://github.com/davatorium/rofi
[rofi-script]: https://davatorium.github.io/rofi/current/rofi-script.5/
[yafl_bookmark]: ./examples/yafl_bookmark
[yafl_moji]: ./examples/yafl_moji
[yafl_search]: ./examples/yafl_search
