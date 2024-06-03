package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Check if there are enough arguments
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments.")
		os.Exit(1)
	}

	// Join the arguments to form the command (skip the first element, which is the program name)
	commandName := "dry-" + strings.Join(os.Args[1:3], "-")
	args := os.Args[3:]

	// Find the executable in the system path
	executable, err := exec.LookPath(commandName)
	if err != nil {
		fmt.Printf("Executable not found: %s\n", commandName)
		os.Exit(1)
	}

	// Execute the command with the provided arguments
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
