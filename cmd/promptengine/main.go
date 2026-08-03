package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LordCodex/promptengine/internal/app"
	"github.com/LordCodex/promptengine/internal/config"
	apperrors "github.com/LordCodex/promptengine/internal/errors"
)

func main() {
	args := os.Args[1:]
	flags := preparseFlags(args)

	cliApp, err := app.Bootstrap(app.BootstrapOptions{
		Out:   os.Stdout,
		Err:   os.Stderr,
		Flags: flags,
		Args:  args,
	})
	if err != nil {
		printError("fatal", err)
		os.Exit(apperrors.ExitCode(err))
	}

	if err := cliApp.Execute(context.Background(), args); err != nil {
		printError("error", err)
		os.Exit(apperrors.ExitCode(err))
	}
}

func printError(prefix string, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		fmt.Fprintf(os.Stderr, "%s: %s\ncategory: %s\nexit_code: %d\n", prefix, appErr.Error(), appErr.Category, appErr.ExitCode())
		if appErr.Recommendation != "" {
			fmt.Fprintf(os.Stderr, "recommendation: %s\n", appErr.Recommendation)
		}
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, err)
}

func preparseFlags(args []string) config.Flags {
	var flags config.Flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			value := true
			flags.Verbose = &value
		case "--debug":
			value := true
			flags.Debug = &value
		case "--json":
			value := true
			flags.JSON = &value
		case "--config":
			if i+1 < len(args) {
				flags.Config = args[i+1]
				i++
			}
		default:
			if strings.HasPrefix(args[i], "--config=") {
				flags.Config = args[i][len("--config="):]
			}
		}
	}
	return flags
}
