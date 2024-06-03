package main

import (
	"os"
)

func main() {
	var command = nexus_upload.NewCmd()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
