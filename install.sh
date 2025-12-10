#!/usr/bin/env bash
#
# NTM Install Script
# https://github.com/Dicklesworthstone/ntm
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Dicklesworthstone/ntm/main/install.sh | bash
#
# Options:
#   --version=TAG   Install specific version (default: latest)
#   --dir=PATH      Install to custom directory (default: /usr/local/bin or ~/.local/bin)
#   --no-shell      Skip shell integration prompt
#

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Defaults
REPO="Dicklesworthstone/ntm"
VERSION=""
INSTALL_DIR=""
NO_SHELL=false

# Parse arguments
for arg in "$@"; do
  case $arg in
    --version=*)
      VERSION="${arg#*=}"
      ;;
    --dir=*)
      INSTALL_DIR="${arg#*=}"
      ;;
    --no-shell)
      NO_SHELL=true
      ;;
    --help|-h)
      echo "NTM Install Script"
      echo ""
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --version=TAG   Install specific version (default: latest)"
      echo "  --dir=PATH      Install to custom directory"
      echo "  --no-shell      Skip shell integration prompt"
      echo "  --help          Show this help"
      exit 0
      ;;
    *)
      echo "Unknown option: $arg"
      exit 1
      ;;
  esac
done

log() {
  echo -e "${GREEN}==>${NC} $1"
}

warn() {
  echo -e "${YELLOW}Warning:${NC} $1"
}

error() {
  echo -e "${RED}Error:${NC} $1" >&2
  exit 1
}

# Detect OS and architecture
detect_platform() {
  local os arch

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m)

  case "$arch" in
    x86_64|amd64)
      arch="amd64"
      ;;
    arm64|aarch64)
      arch="arm64"
      ;;
    *)
      error "Unsupported architecture: $arch"
      ;;
  esac

  case "$os" in
    linux)
      os="linux"
      ;;
    darwin)
      os="darwin"
      ;;
    *)
      error "Unsupported OS: $os"
      ;;
  esac

  echo "${os}-${arch}"
}

# Get the latest release version from GitHub
get_latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

# Determine install directory
get_install_dir() {
  if [[ -n "$INSTALL_DIR" ]]; then
    echo "$INSTALL_DIR"
    return
  fi

  # Prefer /usr/local/bin if writable, otherwise ~/.local/bin
  if [[ -w /usr/local/bin ]]; then
    echo "/usr/local/bin"
  else
    echo "${HOME}/.local/bin"
  fi
}

# Check if a command exists
has_cmd() {
  command -v "$1" >/dev/null 2>&1
}

# Download and install
install_ntm() {
  local platform version install_dir binary_name download_url tmp_dir

  platform=$(detect_platform)
  install_dir=$(get_install_dir)
  binary_name="ntm"

  # Get version
  if [[ -z "$VERSION" ]]; then
    log "Fetching latest version..."
    version=$(get_latest_version)
    if [[ -z "$version" ]]; then
      error "Could not determine latest version"
    fi
  else
    version="$VERSION"
  fi

  log "Installing ntm ${version} for ${platform}"

  # Construct download URL
  download_url="https://github.com/${REPO}/releases/download/${version}/ntm-${platform}"

  # Create temp directory
  tmp_dir=$(mktemp -d)
  trap "rm -rf $tmp_dir" EXIT

  # Download binary
  log "Downloading from ${download_url}..."
  if has_cmd curl; then
    curl -fsSL "$download_url" -o "${tmp_dir}/${binary_name}"
  elif has_cmd wget; then
    wget -q "$download_url" -O "${tmp_dir}/${binary_name}"
  else
    error "Neither curl nor wget found. Please install one."
  fi

  # Make executable
  chmod +x "${tmp_dir}/${binary_name}"

  # Create install directory if needed
  mkdir -p "$install_dir"

  # Install
  log "Installing to ${install_dir}/${binary_name}..."
  if [[ -w "$install_dir" ]]; then
    mv "${tmp_dir}/${binary_name}" "${install_dir}/${binary_name}"
  else
    sudo mv "${tmp_dir}/${binary_name}" "${install_dir}/${binary_name}"
  fi

  # Verify installation
  if "${install_dir}/${binary_name}" version >/dev/null 2>&1; then
    log "Successfully installed ntm ${version}"
  else
    error "Installation verification failed"
  fi

  # Check PATH
  if ! echo "$PATH" | grep -q "$install_dir"; then
    warn "${install_dir} is not in your PATH"
    echo ""
    echo "Add to your shell rc file:"
    echo "  export PATH=\"\$PATH:${install_dir}\""
    echo ""
  fi

  # Shell integration
  if [[ "$NO_SHELL" != true ]]; then
    setup_shell_integration "$install_dir"
  fi

  echo ""
  echo -e "${GREEN}${BOLD}✓ Installation complete!${NC}"
  echo ""
  echo "Quick start:"
  echo "  ntm spawn myproject --cc=2 --cod=2   # Create session with agents"
  echo "  ntm attach myproject                  # Attach to session"
  echo "  ntm palette                           # Open command palette"
  echo ""
  echo "Run 'ntm' for full help."
}

# Setup shell integration
setup_shell_integration() {
  local install_dir="$1"
  local shell_name rc_file init_cmd

  # Detect shell
  shell_name=$(basename "$SHELL")

  case "$shell_name" in
    zsh)
      rc_file="${HOME}/.zshrc"
      init_cmd='eval "$(ntm init zsh)"'
      ;;
    bash)
      rc_file="${HOME}/.bashrc"
      init_cmd='eval "$(ntm init bash)"'
      ;;
    fish)
      rc_file="${HOME}/.config/fish/config.fish"
      init_cmd='ntm init fish | source'
      ;;
    *)
      return
      ;;
  esac

  # Check if already configured
  if [[ -f "$rc_file" ]] && grep -q "ntm init" "$rc_file"; then
    log "Shell integration already configured in ${rc_file}"
    return
  fi

  echo ""
  echo -e "${CYAN}Shell Integration${NC}"
  echo ""
  echo "Add this to ${rc_file}:"
  echo -e "  ${BOLD}${init_cmd}${NC}"
  echo ""

  if [[ -t 0 ]]; then
    printf "Add it now? [y/N]: "
    read -r answer
    case "$answer" in
      y|Y|yes|YES)
        echo "" >> "$rc_file"
        echo "# NTM - Named Tmux Manager" >> "$rc_file"
        echo "$init_cmd" >> "$rc_file"
        log "Added to ${rc_file}"
        echo ""
        echo "Run 'source ${rc_file}' or restart your shell to activate."
        ;;
      *)
        echo "Skipped. Add it manually when ready."
        ;;
    esac
  fi
}

# Run installation
install_ntm
