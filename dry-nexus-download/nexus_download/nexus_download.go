package nexus_download

import (
	"fmt"
	"github.com/spf13/cobra"
	dry_config "github.com/yunarta/dry-tools/dry-config"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// NexusDownload struct to hold service name and root command
type NexusDownload struct {
	username string
	password string
	service  string
	rootCmd  *cobra.Command
}

// Execute method to start the root command
func (n *NexusDownload) Execute() error {
	return n.rootCmd.Execute()
}

// Usage method to print out the usage of the command
func (n *NexusDownload) Usage() error {
	return n.rootCmd.Usage()
}

// Run method to execute the command
func (n *NexusDownload) Run(cmd *cobra.Command, args []string) {
	err := n.run(cmd, args)
	if err != nil {
		log.Println("Error running command: ", err)
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

// main run function
func (n *NexusDownload) run(cmd *cobra.Command, args []string) error {
	// check if command arguments and flags are empty
	if len(args) == 0 && cmd.Flags().NFlag() == 0 {
		_ = cmd.Usage()
	} else {
		// load the configuration file
		config, err := dry_config.LoadConfig()
		if err != nil {
			log.Println("Error loading config:", err)
			return fmt.Errorf("failed when loading config: %w", err)
		}

		var username, password string
		// if n.username and n.password is not null
		if n.username != "" && n.password != "" {
			username = n.username
			password = n.password
		} else {
			// Resolving configuration for given service
			resolve, err := config.Resolve("nexus", n.service)
			if err != nil {
				log.Println("Error resolving config:", err)
				return err
			}

			// Extracting username and password from configuration
			username = resolve["username"]
			password = resolve["password"]
		}

		src := args[0]
		dest := args[1]

		log.Println("Finding files")
		// find remote files to download
		downloads, single, err := findRemoteFiles(src, username, password)
		if err != nil {
			log.Println("Error finding files:", err)
			return fmt.Errorf("failed when trying to find files: %w", err)
		}

		// get the absolute path of destination
		path, err := filepath.Abs(dest)
		if err != nil {
			log.Println("Error getting absolute path:", err)
			return fmt.Errorf("failed when retrieving absolute path: %w", err)
		}

		targetFile := ""
		isFile := false

		pathInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			targetFile = path
			isFile = true
		} else if !pathInfo.IsDir() {
			targetFile = path
			isFile = true
		}

		// single and multiple file download cases
		if single {
			downloadURL := downloads[0].DownloadURL
			parsedDownloadURL, _ := url.Parse(downloadURL)
			fileName := filepath.Base(parsedDownloadURL.Path)

			if !isFile {
				targetFile = filepath.Join(path, fileName)
			}

			log.Println("Downloading file", downloadURL, "to", targetFile)
			err = downloadFile(downloadURL, targetFile)
			if err != nil {
				log.Println("Error downloading file:", err)
				return fmt.Errorf("failed when downloading single file: %w", err)
			}
		} else {
			if isFile {
				log.Println("Attempt of downloading multiple files into single path")
				return fmt.Errorf("attempt of downloading multiple files into single path")
			}

			sourceURL, _ := url.Parse(src)
			destDirectory, _ := strings.CutSuffix(path, "/")

			for _, download := range downloads {
				downloadURL, _ := url.Parse(download.DownloadURL)
				relativeFile, _ := filepath.Rel(sourceURL.Path, downloadURL.Path)
				targetFile = filepath.Join(destDirectory, filepath.Clean(relativeFile))

				targetParentDir := filepath.Dir(targetFile)
				_ = os.MkdirAll(targetParentDir, os.ModePerm)

				log.Println("Downloading file", downloadURL, "to", targetFile)
				err = downloadFile(download.DownloadURL, targetFile)
				if err != nil {
					log.Println("Error downloading file:", err)
					return fmt.Errorf("failed when downloading multiple files: %w", err)
				}
			}
		}
	}
	return nil
}

// NewCmd creates a new NexusDownload command
func NewCmd() *NexusDownload {
	var cmd = NexusDownload{}

	var rootCmd = &cobra.Command{
		Use: "dry nexus download [src] [dst]",
		Run: cmd.Run,
	}

	// setting service flag
	rootCmd.Flags().StringVarP(&cmd.service, "service", "s", "nexus", "Service")
	cmd.rootCmd = rootCmd

	return &cmd
}
