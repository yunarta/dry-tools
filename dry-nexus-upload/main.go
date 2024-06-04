package main

import (
	"github.com/yunarta/dry-tools/dry-nexus-upload/nexus_upload"
	"os"
)

func main() {
	var command = nexus_upload.NewCmd()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
