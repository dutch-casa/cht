#!/usr/bin/env bash
set -euo pipefail

REPO="dutch-casa/cht"
INSTALL_DIR="/usr/local/bin"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

info() { printf '\033[1;34m==> %s\033[0m\n' "$*"; }
err()  { printf '\033[1;31merror: %s\033[0m\n' "$*" >&2; exit 1; }

# --- Detect platform ---

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
UARCH="$(uname -m)"

case "$UARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
    arm64)   GOARCH="arm64" ; UARCH="aarch64" ;;
    *)       err "Unsupported architecture: $UARCH" ;;
esac

# bat uses x86_64/aarch64 naming, fzf uses amd64/arm64.
case "$UARCH" in
    aarch64) BAT_ARCH="aarch64" ; FZF_ARCH="arm64" ;;
    x86_64)  BAT_ARCH="x86_64"  ; FZF_ARCH="amd64" ;;
esac

# --- Download helper ---

fetch() {
    local url="$1" dest="$2"
    if ! curl -fsSL -o "$dest" "$url"; then
        err "Download failed: $url"
    fi
}

# --- cht ---

if ! command -v cht >/dev/null 2>&1; then
    info "Installing cht..."
    fetch "https://github.com/${REPO}/releases/latest/download/cht-${OS}-${GOARCH}" "$TMP/cht"
    chmod +x "$TMP/cht"
    sudo mv "$TMP/cht" "$INSTALL_DIR/cht"
else
    info "cht already installed, skipping."
fi

# --- bat ---

if ! command -v bat >/dev/null 2>&1; then
    info "Installing bat..."
    BAT_VERSION="$(curl -fsSL "https://api.github.com/repos/sharkdp/bat/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
    case "$OS" in
        darwin) BAT_TARGET="${BAT_ARCH}-apple-darwin" ;;
        linux)  BAT_TARGET="${BAT_ARCH}-unknown-linux-musl" ;;
    esac
    BAT_ASSET="bat-${BAT_VERSION}-${BAT_TARGET}"
    fetch "https://github.com/sharkdp/bat/releases/latest/download/${BAT_ASSET}.tar.gz" "$TMP/bat.tar.gz"
    tar xzf "$TMP/bat.tar.gz" -C "$TMP"
    sudo mv "$TMP/${BAT_ASSET}/bat" "$INSTALL_DIR/bat"
else
    info "bat already installed, skipping."
fi

# --- fzf ---

if ! command -v fzf >/dev/null 2>&1; then
    info "Installing fzf..."
    FZF_VERSION="$(curl -fsSL "https://api.github.com/repos/junegunn/fzf/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
    FZF_ASSET="fzf-${FZF_VERSION#v}-${OS}_${FZF_ARCH}"
    fetch "https://github.com/junegunn/fzf/releases/latest/download/${FZF_ASSET}.tar.gz" "$TMP/fzf.tar.gz"
    tar xzf "$TMP/fzf.tar.gz" -C "$TMP"
    sudo mv "$TMP/fzf" "$INSTALL_DIR/fzf"
else
    info "fzf already installed, skipping."
fi

# --- Zsh completions ---

if [[ "$SHELL" == */zsh ]]; then
    comp_dir="${HOME}/.zsh/completions"
    mkdir -p "$comp_dir"
    "$INSTALL_DIR/cht" --completions > "$comp_dir/_cht"

    if ! grep -q '\.zsh/completions' ~/.zshrc 2>/dev/null; then
        info "Adding completions to .zshrc..."
        printf '\nfpath=(~/.zsh/completions $fpath)\nautoload -Uz compinit && compinit\n' >> ~/.zshrc
    fi
fi

# --- Cache ---

info "Building topic cache..."
"$INSTALL_DIR/cht" --refresh

info "Done. Run 'cht' to start."
