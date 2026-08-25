# hypr-cfg

A terminal UI for browsing and editing [Hyprland](https://hypr.land) configuration files.

`hypr-cfg` parses your `hyprland.conf`, presents it as a navigable tree grouped by
category (monitors, keybindings, window rules, …), and lets you edit values in place.
Comments, blank lines, and file layout are preserved on save, and a backup copy of the
original file is written before any change touches disk.

## Features

- **Tree view** of the whole config, grouped into:
  - Variables (`$terminal = kitty`)
  - Monitors (`monitor=` / `monitorv2 {}`)
  - Env Variables (`env =`)
  - Autostart (`exec-once =` and other `exec*` keywords)
  - Bindings (`bind`, `bindm`, `bindel`, `bindl`, …)
  - Gestures (`gesture =`)
  - Window Rules, Layer Rules, Workspace Rules
  - Permissions (`permission =`)
  - Sources (`source =`)
  - Other — all `section {}` blocks (`general`, `decoration`, `input`, nested blocks like `shadow {}` and `touchpad {}`, `windowrule {}` blocks, …)
- **In-place editing** — select any leaf to edit its value in a form, then save back to the file.
- **Non-destructive saves** — original formatting and comments are kept; a `~hyprland.conf` backup is created next to the source file before writing.
- **Current syntax support** — targets the latest classic (hyprlang) config syntax as of Hyprland v0.54+, including `windowrule {}` blocks, `gesture`, `permission`, and workspace rules. Older configs (e.g. `windowrulev2 =` lines) still parse fine.

> **Note:** Hyprland v0.55+ also offers an optional Lua-based config (`hyprland.lua`).
> `hypr-cfg` works with the classic `.conf` format, which remains fully supported by Hyprland.

## Installation

Requires Go 1.21+.

```sh
git clone <this repo>
cd hypr-cfg
make build
```

The binary is produced at `bin/hypr-cfg`.

## Usage

```sh
./bin/hypr-cfg ~/.config/hypr/hyprland.conf
```

Or build and launch against your config in one step:

```sh
make run
```

### Controls

- **Arrow keys / mouse** — navigate the tree
- **Enter / click** — expand or collapse a group; open the edit form on a value
- In the edit form: **Save** writes the change back to the file, **Cancel** discards it

## Development

```sh
make unit-test   # run the test suite
make build       # test + build bin/hypr-cfg
```

Project layout:

| Path | Purpose |
|------|---------|
| `cmd/cli/` | CLI entry point |
| `pkg/parser.go` | Line-based config parser, tree builder, grouped model, file writer |
| `pkg/ui.go` | tview-based terminal UI |
| `tests/` | Unit tests and example config fixtures (current + legacy syntax) |

The parser is intentionally line-oriented: every line of the source file is kept as a
`RawLine` in a linked list, so edits only touch the line being changed and the file
round-trips byte-for-byte otherwise.

## License

[MIT](LICENSE) © MB Griaustinis Media
