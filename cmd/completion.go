package cmd

import (
    "os"

    "github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate completion script",
    Long: `To load completions:

Bash:
  $ source <(nexuscli completion bash)
  $ nexuscli completion bash > /etc/bash_completion.d/nexuscli

Zsh:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  $ nexuscli completion zsh > "${fpath[1]}/_nexuscli"

Fish:
  $ nexuscli completion fish | source
  $ nexuscli completion fish > ~/.config/fish/completions/nexuscli.fish

PowerShell:
  PS> nexuscli completion powershell | Out-String | Invoke-Expression
  PS> nexuscli completion powershell > nexuscli.ps1
`,
    DisableFlagsInUseLine: true,
    ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
    Args:                  cobra.ExactValidArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        switch args[0] {
        case "bash":
            _ = rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            _ = rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            _ = rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            _ = rootCmd.GenPowerShellCompletion(os.Stdout)
        }
    },
}

func init() {
    rootCmd.AddCommand(completionCmd)
}
