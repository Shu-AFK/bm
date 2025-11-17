# bm - bookmark manager CLI

Welcome to **bm**, a tiny command-line sidekick for stashing your favorite files, folders, and links. It keeps everything in a simple JSON file under your XDG data directory, so it plays nicely with dotfile syncing and cloud backups.

## What it does
- **Save quick shortcuts** with `bm add <name> <path/url>` and optional tags.
- **Browse your stash** with `bm list`, filter by tag, and sort the way you like.
- **Jump right in** using `bm open <name>` to launch files, folders, or URLs.
- **Tweak or tidy** with `bm edit` and `bm remove` when things change.

## Installation
- **Go users:** `go install github.com/Shu-AFK/bm@latest`
- **Binary packages:** I'll publish builds for popular platforms soon (it should soon be downloadable by your favourite package manager).

## Usage
Type `bm --help` to see the colorful help view, or append `--help` to any subcommand for details. Add a bookmark like this:

```bash
bm add docs ~/projects/docs/doc1.md
bm add homepage https://example.com --tags personal,reading
```

Then open it later with:

```bash
bm open docs
```

## Config & storage
Bookmarks live in `~/.local/share/bm/bookmarks.json` (or the equivalent XDG data path on your OS). The tool doesn't create extra config files, so it's easy to inspect or edit manually if you ever need to.

## Roadmap
- A shell script to make opening a path cd into it
- Prebuilt packages for more distros and package managers.

## Contributing
Feel free to file issues or PRs. If you spot a typo or have an idea for a smoother flow, I'd love to hear it.
