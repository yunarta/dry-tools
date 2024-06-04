package nexus_upload

import (
	"fmt"
	"github.com/spf13/cobra"
	dry_config "github.com/yunarta/dry-tools/dry-config"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// NexusUpload is our main structure.
type NexusUpload struct {
	service string
	recurse bool
	rootCmd *cobra.Command // rootCmd will hold our Cobra command.
}

// Execute is a wrapper around Cobra's Execute method.
func (n *NexusUpload) Execute() error {
	return n.rootCmd.Execute()
}

// Usage is a wrapper around Cobra's Usage method.
func (n *NexusUpload) Usage() error {
	return n.rootCmd.Usage()
}

// This Run method will be called by Cobra when the appropriate command is called.
func (n *NexusUpload) Run(cmd *cobra.Command, args []string) {
	err := n.run(cmd, args) // we delegate the running of the command to our private run method.
	if err != nil {
		log.Println("Error running command: ", err)
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

// This private run method does the actual work of running our command.
func (n *NexusUpload) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && cmd.Flags().NFlag() == 0 {
		_ = cmd.Usage()
	} else {
		// Loading configuration from file
		config, err := dry_config.LoadConfig()
		if err != nil {
			log.Println("Error loading config:", err)
			return err
		}

		// Resolving configuration for given service
		resolve, err := config.Resolve("nexus", n.service)
		if err != nil {
			log.Println("Error resolving config:", err)
			return err
		}

		// Extracting username and password from configuration
		username := resolve["username"]
		password := resolve["password"]

		// Source and destination paths from command arguments
		src := args[0]
		dest := args[1]

		log.Println("Finding files")
		// Finding the files in given directory (src)
		files, err := findFiles(src, n.recurse)

		if err != nil {
			log.Println("Error finding files:", err)
			return err
		}

		// Cut off trailing slash in dest if there is one
		endpoint, _ := strings.CutSuffix(dest, "/")
		for _, file := range files {
			localFilePath := file
			// Determine path to file relative to src
			relativeToSrc, _ := filepath.Rel(src, localFilePath)
			targetURLPath := endpoint + "/" + relativeToSrc

			log.Println("Uploading file", localFilePath, "to", targetURLPath)
			// Attempt to upload each file
			err = uploadFile(localFilePath, targetURLPath, username, password)
			if err != nil {
				log.Println("Error uploading file:", err)
				return err
			}
		}
	}

	return nil
}

// NewCmd creates a new instance of NexusUpload,
// initializes it with a new Cobra command, and returns a pointer to it.
func NewCmd() *NexusUpload {
	var cmd = NexusUpload{}

	var rootCmd = &cobra.Command{
		Use: "dry nexus upload [src] [dst]",
		Run: cmd.Run,
	}

	// Adding flags to our command
	rootCmd.Flags().StringVarP(&cmd.service, "service", "s", "nexus", "Service")
	rootCmd.Flags().BoolVarP(&cmd.recurse, "recurse", "r", false, "Recurse into directories")
	cmd.rootCmd = rootCmd

	return &cmd
}
