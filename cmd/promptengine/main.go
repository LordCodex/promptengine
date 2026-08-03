package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LordCodex/promptengine/internal/app"
	"github.com/LordCodex/promptengine/internal/errors"
)

func main() {
	// Parse verbose flag from args manually for early bootstrap logging configuration
	verbose := false
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
			break
		}
	}

	ctx := context.Background()
	cliApp, err := app.Bootstrap(os.Stdout, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: Bootstrapping failed: %v\n", err)
		os.Exit(errors.ExitGeneralError)
	}

	// Skip first arg (binary name)
	if err := cliApp.Execute(ctx, os.Args[1:]); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			os.Exit(appErr.Code)
		}
		os.Exit(errors.ExitGeneralError)
	}
}
