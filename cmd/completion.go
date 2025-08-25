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

  # To load completions for each session, execute once:
  # Linux:
  $ nexuscli completion bash > /etc/bash_completion.d/nexuscli
  # macOS:
  $ nexuscli completion bash > /usr/local/etc/bash_completion.d/nexuscli

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ nexuscli completion zsh > "${fpath[1]}/_nexuscli"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ nexuscli completion fish | source

  # To load completions for each session, execute once:
  $ nexuscli completion fish > ~/.config/fish/completions/nexuscli.fish

PowerShell:

  PS> nexuscli completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> nexuscli completion powershell > nexuscli.ps1
  # and source this file from your PowerShell profile.
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
