#!/bin/bash

set -e

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RESET='\033[0m'
GO_VERSION='go1.22.4'
shells=("bash" "zsh" "fish")

display_shells() {
    echo "Choose your favorite Shell number:"
    for i in "${!shells[@]}"; do
        echo "$((i+1)) ${shells[$i]}"
    done
}

get_rc_file() {
    local shell_name=$1
    local rc_file=""

    case "$shell_name" in
        "bash") rc_file="$HOME/.bashrc" ;;
        "zsh") rc_file="$HOME/.zshrc" ;;
        "fish") rc_file="$HOME/.config/fish/config.fish" ;;
        *)
            echo "Error: Unknown shell '$shell_name'."
            return 1
            ;;
    esac
    echo "$rc_file"
    return 0
}

echo -e "${YELLOW} -> Check if GO is installed ${RESET}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN} $GO_VERSION is installed on your system ${RESET}"
else
    echo -e "${YELLOW} -> Installing GO ${RESET}"
    echo -e "${YELLOW} -> Downloading GO source code ${RESET}"
    wget https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz
    echo -e "${YELLOW} -> Extracting files ${RESET}"
    sudo tar -C /usr/local -xzf "${GO_VERSION}.linux-amd64.tar.gz"


fi
