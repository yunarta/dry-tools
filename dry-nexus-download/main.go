package main

import (
	"github.com/yunarta/dry-tools/dry-nexus-download/nexus_download"
	"os"
)

func main() {
	var command = nexus_download.NewCmd()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
