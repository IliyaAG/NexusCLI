#!/bin/bash

set -e

RED='\033[0;31m'
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
            echo -e "${RED}Error: Unknown shell '$shell_name'.${RESET}"
            return 1
            ;;
    esac
    echo "$rc_file"
    return 0
}




display_shells

read -p "Enter your choice (number): " choice

if ! [[ "$choice" =~ ^[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid input. Please enter a number.${RESET}"
    exit 1
fi

num_shells=${#shells[@]}
if (( choice < 1 || choice > num_shells )); then
    echo -e "${RED}Error: Your choice '$choice' is out of range."
    echo -e "Please enter a number between 1 and $num_shells.${RESET}"
    exit 1
fi

selected_index=$((choice - 1))
selected_shell=${shells[$selected_index]}

echo -e "${YELLOW}-> You selected shell: '$selected_shell'.${RESET}"

rc_file=$(get_rc_file "$selected_shell")
if [ $? -ne 0 ]; then
    exit 1
fi

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
    echo -e "${YELLOW} -> Add go to PATH ${RESET}"
    echo 'export PATH=$PATH:/usr/local/go/bin' >> $rc_file
fi
