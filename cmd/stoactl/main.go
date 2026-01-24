package main

import (
	"os"

	"github.com/stoa-platform/stoactl/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
