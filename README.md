<p align="center">
  <img src="assets/logo.png" width="250" style="display:block; margin:auto;" alt="bm logo"/>
</p>

<h1 align="center">bm - Bookmark Manager CLI</h1>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/Shu-AFK/bm"><img src="https://goreportcard.com/badge/github.com/Shu-AFK/bm" alt="Go Report Card"></a>
  <a href="https://github.com/Shu-AFK/bm/releases"><img src="https://img.shields.io/github/v/release/Shu-AFK/bm" alt="Latest Release"></a>
  <a href="https://github.com/Shu-AFK/bm/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Shu-AFK/bm" alt="License"></a>
</p>

---

**bm** is a small, focused command-line tool for storing shortcuts to files, folders, and URLs. Everything is kept in a single JSON file inside your XDG data directory, which makes it easy to sync, back up, or inspect manually.

## Features

* **Add shortcuts quickly** using `bm add <name> <path-or-url>` with optional tags.
* **List and filter** your bookmarks, including sorting by name, creation time, or target.
* **Open anything instantly** with `bm open <name>`, directories, files, and URLs all open with your platform’s default applications.
* **Install shell functions** using `bm install <shell> <show|apply>` (bash, zsh, fish, powershell). This enables `bm` to open directories with shell-native navigation.
* **Clean up or adjust entries** with `bm edit` and `bm remove`.

## Installation

### Go users

```bash
go install github.com/Shu-AFK/bm@latest
```

### Prebuilt binaries & Package Managers

Prebuilt binaries will be available for:

* **Homebrew** (macOS, Linux)
* **Scoop** (Windows)
* **AUR (bm-bin)** for Arch-based distros
* **deb / rpm / apk** packages for Linux

### Example installation commands (once packages are published)

**Homebrew:**

```bash
brew install Shu-AFK/tap/bm
```

**Scoop:**

```powershell
scoop bucket add bm https://github.com/Shu-AFK/scoop-bucket
scoop install bm
```

**AUR:**

```bash
yay -S bm-bin
```

**deb/rpm/apk:**
Download from the GitHub Releases page.

## Shell Completions

`bm` can generate completions for bash, zsh, fish, and powershell:

```bash
bm completion bash
bm completion zsh
bm completion fish
bm completion powershell
```

You can redirect these into your shell's completion directory if you want autocomplete.

## Usage

Run `bm --help` for an overview, or `bm <command> --help` for command-specific details.

Add bookmarks:

```bash
bm add docs ~/projects/docs/doc1.md
bm add homepage https://example.com --tags personal,reading
```

Open a saved shortcut:

```bash
bm open docs
```

List, filter, and sort:

```bash
bm list
bm list --tag reading
bm list --sort type
```

Install shell integration for directory opening:

```bash
bm install bash
bm install zsh
```

## Storage

Bookmarks are stored in:

```
~/.local/share/bm/bookmarks.json
```

(or your platform’s XDG data location). The format is simple and human-readable.

## Roadmap

* Extended metadata (notes, favorites).
* Prebuilt packages for additional distros.

## Contributing

Issues and pull requests are welcome. If something feels rough or unclear, suggestions are appreciated.
