#!/bin/bash

set -e

print_coloric() {
    local color=$1
    shift
    local text="$*"

    case $color in
        red)    echo -e "\033[0;31mi -> $text\033[0m";;
        green)  echo -e "\033[0;32m -> $text\033[0m";;
        yellow) echo -e "\033[0;33m -> $text\033[0m";;
        *) echo -e "-> $text";;
    esac
}

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
            print_coloric red "Error: Unknown shell '$shell_name'."
            return 1
            ;;
    esac
    echo "$rc_file"
    return 0
}


display_shells

read -p "Enter your choice (number): " choice

if ! [[ "$choice" =~ ^[0-9]+$ ]]; then
    print_coloric red "Error: Invalid input. Please enter a number."
    exit 1
fi

num_shells=${#shells[@]}
if (( choice < 1 || choice > num_shells )); then
    print_coloric red "Error: Your choice '$choice' is out of range."
    print_coloric red "Please enter a number between 1 and $num_shells."
    exit 1
fi

selected_index=$((choice - 1))
selected_shell=${shells[$selected_index]}

print_coloric yellow "You selected shell: '$selected_shell'."

rc_file=$(get_rc_file "$selected_shell")
if [ $? -ne 0 ]; then
    exit 1
fi

print_coloric yellow "Check if GO is installed"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_coloric green "$GO_VERSION is installed on your system"
else
    print_coloric yellow "Installing GO"
    print_coloric yellow "Downloading GO source code"
    wget https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz
    print_coloric yellow "Extracting files"
    sudo tar -C /usr/local -xzf "${GO_VERSION}.linux-amd64.tar.gz"
    print_coloric yellow "Add go to PATH"
    echo 'export PATH=$PATH:/usr/local/go/bin' >> $rc_file
    source $rc_file
fi
print_coloric yellow "Building nexuscli"
sudo go build -o /usr/local/bin/nexuscli .
print_coloric yellow "Add nexuscli auto completion"
/usr/local/bin/nexuscli completion $selected_shell
