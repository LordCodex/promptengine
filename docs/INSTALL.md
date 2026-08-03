# Installation Guide

This document details how to install the PromptEngine CLI on various platforms.

---

## 1. Install Script (macOS / Linux)

The easiest way to install PromptEngine is using our shell script:

```bash
curl -fsSL https://raw.githubusercontent.com/LordCodex/promptengine/main/scripts/install.sh | bash
```

To configure custom install location, set the `PROMPTENGINE_INSTALL_DIR` env:

```bash
curl -fsSL https://raw.githubusercontent.com/LordCodex/promptengine/main/scripts/install.sh | PROMPTENGINE_INSTALL_DIR=$HOME/bin bash
```

---

## 2. Install Script (Windows PowerShell)

For Windows, run the following PowerShell command in an elevated session:

```powershell
iwr -useb https://raw.githubusercontent.com/LordCodex/promptengine/main/scripts/install.ps1 | iex
```

---

## 3. Go Install (All Platforms)

If you have Go 1.21+ installed, build and install from source directly:

```bash
go install github.com/LordCodex/promptengine/cmd/promptengine@latest
```

---

## Unsupported: Homebrew

No Homebrew tap is included in the v1.0 release. Use the macOS/Linux installer script or `go install` instead.

The supported installation methods are listed above.

---

## 4. Shell Autocompletions Setup

To configure shell completion support, run the helper script:

```bash
curl -fsSL https://raw.githubusercontent.com/LordCodex/promptengine/main/scripts/completions.sh | bash
```

Or reference manual shell generation commands:

```bash
# Bash
promptengine completion bash > /etc/bash_completion.d/promptengine

# Zsh
promptengine completion zsh > "${fpath[1]}/_promptengine"

# Fish
promptengine completion fish > ~/.config/fish/completions/promptengine.fish
```
