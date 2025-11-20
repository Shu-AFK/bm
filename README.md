<p align="center">
  <img src="assets/logo.png" width="250" alt="bm logo"/>
</p>

<h1 align="center">bm — Bookmark Manager CLI</h1>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/Shu-AFK/bm"><img src="https://goreportcard.com/badge/github.com/Shu-AFK/bm" alt="Go Report Card"></a>
  <a href="https://github.com/Shu-AFK/bm/releases"><img src="https://img.shields.io/github/v/release/Shu-AFK/bm" alt="Latest Release"></a>
  <a href="https://github.com/Shu-AFK/bm/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Shu-AFK/bm" alt="License"></a>
</p>

---

**bm** is your command-line sidekick for teleporting to files, folders, and URLs in one hop.  
It's fast, minimal, and keeps all bookmarks in a single human-readable JSON file under your XDG data directory.

---

# Screenshots
## Built-in help, commands, and clean terminal UI
<img src="assets/main.png" width="650"/>

## Add bookmarks instantly
<img src="assets/add.png" width="650"/>

## List bookmarks with clean formatting and colors
<img src="assets/list.png" width="650"/>

## Open anything, files, directories, and URLs
<img src="assets/open.png" width="900"/>

---

# Fuzzy search and the interactive selector

When you run:

```bash
bm open <name>
```

bm uses fuzzy matching to find the closest bookmark, even if you only type part of the name or make a typo.

If several results are similar, bm automatically shows an interactive terminal selector so you can choose the exact one you want without retyping.

This makes navigation extremely fast while still being precise.

---

# Features

- Add bookmarks quickly with `bm add <name> <path-or-url>`.
- Open directories, files, or URLs with `bm open <name>`.
- Fuzzy matching for fast navigation.
- Interactive selector when fuzzy matches are close.
- List and filter bookmarks with clean table output.
- Optional shell integration for directory navigation (`cd`).
- Everything stored in a single JSON file.

---

# Installation

## From source (Go)

```bash
go install github.com/Shu-AFK/bm@latest
```

## Package Managers

### Homebrew (macOS & Linux)
```bash
brew install Shu-AFK/tap/bm
```

### Scoop (Windows)
```powershell
scoop bucket add bm https://github.com/Shu-AFK/scoop-bucket
scoop install bm
```

### AUR (Arch Linux)
```bash
yay -S bm-bin
```

### Debian / RPM / Arch packages
Download from the GitHub Releases page.

---

# Usage examples

### Add bookmarks
```bash
bm add docs ~/projects/docs/doc1.md
bm add homepage https://example.com --tags personal,reading
```

### Open bookmarks
```bash
bm open docs
bm open home
```

### List, filter, and sort
```bash
bm list
bm list --tag reading
bm list --sort type
```

### Install shell integration (either show to get the function or apply to put it directly into your profile)
```bash
bm install bash show
bm install zsh apply
bm install fish apply
```

### Shell completions
```bash
bm completion bash
bm completion zsh
bm completion fish
bm completion powershell
```

---

# Storage

Bookmarks are stored at:

```
~/.local/share/bm/bookmarks.json
```

---

# Roadmap

- Extra metadata (notes, favorites)
- More sort/filter options
- Additional package manager support

---

# Contributing

Contributions, issues, and feature suggestions are welcome.
