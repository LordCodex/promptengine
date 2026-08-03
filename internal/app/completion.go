package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// BuildCompletion constructs the shell completion command builder
func (b *CommandBuilder) BuildCompletion() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `To load completions:

Bash:

  $ source <(promptengine completion bash)

  # To load completions for each session, add to ~/.bashrc:
  $ promptengine completion bash > /etc/bash_completion.d/promptengine

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can run the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, add to ~/.zshrc:
  $ promptengine completion zsh > "${fpath[1]}/_promptengine"

Fish:

  $ promptengine completion fish > ~/.config/fish/completions/promptengine.fish

PowerShell:

  PS> promptengine completion powershell | Out-String | Invoke-Expression
`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(b.Out)
			case "zsh":
				return cmd.Root().GenZshCompletion(b.Out)
			case "fish":
				return cmd.Root().GenFishCompletion(b.Out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(b.Out)
			default:
				return fmt.Errorf("unsupported shell Type %q", args[0])
			}
		},
	}
}

func init() {
	// Optional early environment validation
	if os.Getenv("PROMPTENGINE_NO_COLOR") != "" {
		// handle coloring defaults
	}
}
