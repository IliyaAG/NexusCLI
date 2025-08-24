#!/bin/bash
shells=("bash" "zsh" "ksh" "fish" "tcsh")

display_shells() {
    echo "Please select one of the following shells:"
    for i in "${!shells[@]}"; do
        echo "$((i+1)) ${shells[$i]}"
        echo "--------------------------------"
    done
}

display_shells
