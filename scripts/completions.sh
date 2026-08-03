#!/bin/bash
set -e

# PromptEngine Shell Completion Installer Helper
# Automatically detects active shell and wires completions.

SHELL_NAME=$(basename "$SHELL")

if ! command -v promptengine >/dev/null 2>&1; then
  echo "Error: promptengine CLI is not installed on PATH."
  exit 1
fi

echo "Installing shell completions for $SHELL_NAME..."

case "$SHELL_NAME" in
  bash)
    TARGET_DIR="$HOME/.local/share/bash-completion/completions"
    mkdir -p "$TARGET_DIR"
    promptengine completion bash > "$TARGET_DIR/promptengine"
    echo "✓ Bash completion script saved to $TARGET_DIR/promptengine"
    echo "Ensure 'bash-completion' is loaded in your shell."
    ;;
  zsh)
    TARGET_DIR="${fpath[1]}"
    if [ -z "$TARGET_DIR" ] || [ ! -w "$TARGET_DIR" ]; then
      TARGET_DIR="$HOME/.zsh/completions"
      mkdir -p "$TARGET_DIR"
      echo "fpath=( $TARGET_DIR \$fpath )" >> "$HOME/.zshrc"
    fi
    promptengine completion zsh > "$TARGET_DIR/_promptengine"
    echo "✓ Zsh completion script saved to $TARGET_DIR/_promptengine"
    echo "Restart shell or run: autoload -U compinit && compinit"
    ;;
  fish)
    TARGET_DIR="$HOME/.config/fish/completions"
    mkdir -p "$TARGET_DIR"
    promptengine completion fish > "$TARGET_DIR/promptengine.fish"
    echo "✓ Fish completion script saved to $TARGET_DIR/promptengine.fish"
    ;;
  *)
    echo "Shell '$SHELL_NAME' is not supported by automated installer."
    echo "Run 'promptengine completion --help' for manual setup instructions."
    exit 1
    ;;
esac
