package recovery

import (
	"fmt"
	"io"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// CommandPanicWrapper returns a Cobra run function wrapped with panic recovery.
func CommandPanicWrapper(out io.Writer, run func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(out, "Error: PromptEngine encountered an unexpected internal panic: %v\n", r)
				fmt.Fprintln(out, "Please submit a bug report with the following debug stack trace:")
				fmt.Fprintf(out, "%s\n", debug.Stack())
				err = fmt.Errorf("internal panic: %v", r)
			}
		}()
		return run(cmd, args)
	}
}
