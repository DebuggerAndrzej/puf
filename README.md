<div align="center" width="100%">
    <img src="https://github.com/DebuggerAndrzej/puf/assets/118397780/2f492e05-2613-48ed-a5d4-551c60dd93f8" width="200">
</div>
<h2 align="center">PUF - Partially Unzip File</h2>
A simple go tui tool to help you unzip only files that are useful. Works well with compressed files within compressed files.

# Installation
```
 go install github.com/DebuggerAndrzej/puf@latest
```
Requirements:
- go 1.22 or newer

> [!TIP]
> default installation path for go is ~/go/bin so in order to have tli command available this path has to be added to shell user paths.

# Flags

- `f` - path to archive you want to extract from
- `d` - path to destination folder (optional)
> [!NOTE]
> all paths can be relative or absolute
- `r` - regex choosing which files to show (optional)
> [!TIP]
> if you want to search only for files with .go extension you can use `-r go\$` for fish and `-r go$` for bash.

# Shortcuts

**Selection:**
- `tab` or `space` - (de)select
- `a` - (de)select all

**Basic movements:**
- `up` or `k` -  up
- `down` or `j` -  down
- `left` or `h` - left (page)
- `right` or `l` - right (page)

**Quit:**
- `q` or `ctrl+c` - quit
