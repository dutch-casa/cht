# cht

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)]()
[![Release](https://img.shields.io/github/v/release/dutch.casa/cht)](https://github.com/dutch-casa/cht/releases/latest)

A fast CLI for browsing [cheat.sh](https://cheat.sh) cheat sheets with syntax highlighting and fuzzy search.

---

## Demo

```
$ cht curl
# curl
# Transfer data from or to a server using various protocols.

# Download a file from a URL and save it with a specific name
curl -o filename.ext http://example.com/file.txt

# Follow redirects if the URL has moved
curl -L http://example.com
...
```

## Install

One command. No dependencies required — the installer handles everything.

```bash
curl -fsSL https://raw.githubusercontent.com/dutch.casa/cht/main/install.sh | bash
```

This will:

- Download the `cht` binary for your platform
- Install `bat` and `fzf` if not already present
- Set up zsh tab completions
- Build the local topic cache

### Requirements

- macOS (arm64 / amd64) or Linux (arm64 / amd64)
- `curl` and `bash` (preinstalled on both platforms)

## Usage

```
cht                     # fuzzy search all topics (default)
cht <topic>             # fetch a specific cheat sheet
cht python/lambda       # fetch a subtopic
cht --help              # print usage
cht --refresh           # rebuild the topic cache
cht --completions       # print zsh completion script
```

### Tab Completion

If the installer set up completions, you can tab-complete topic names:

```
$ cht cu<TAB>
curl  cups  cut  ...
```

Works with [fzf-tab](https://github.com/Aloxaf/fzf-tab) for fuzzy matching in the completion menu.

### Fuzzy Mode

Running `cht` with no arguments opens an interactive fuzzy finder over all available topics. The preview pane shows the cheat sheet contents as you navigate.

## How It Works

- Fetches plain text from `cheat.sh` and pipes it through `bat` for syntax highlighting and paging
- Topic list is cached locally at `~/.cache/cht/topics` and refreshed on demand
- Cache staleness warning appears after 30 days but never blocks

## Building from Source

Requires Go 1.25+.

```bash
git clone https://github.com/dutch-casa/cht.git
cd cht
make build
```

### Cross-compile release binaries

```bash
make release
```

Produces binaries in `dist/` for darwin/arm64, darwin/amd64, linux/arm64, and linux/amd64.

### Creating a release

```bash
git tag v1.0.0
git push --tags
make release
gh release create v1.0.0 dist/* --title "v1.0.0"
```

## Uninstall

```bash
sudo rm /usr/local/bin/cht
rm -rf ~/.cache/cht
rm -f ~/.zsh/completions/_cht
```

## License

[MIT](LICENSE)
