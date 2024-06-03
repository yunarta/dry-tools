package nexus_upload

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFile(localFilePath string, targetURLPath string, username string, password string) error {
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(localFilePath))

	if err != nil {
		log.Println("Error creating multipart form file:", err)
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Println("Error copying to part:", err)
		return err
	}

	err = writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return err
	}

	req, err := http.NewRequest("PUT", targetURLPath, body)
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error doing HTTP Request:", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected status 200 OK or 201 Created, got " + resp.Status)
	}

	return nil
}

func findFiles(path string, recurse bool) ([]string, error) {
	var files []string

	if recurse {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("Filepath Walk error:", err)
				return err
			}

			if !info.IsDir() {
				files = append(files, path)
				log.Println("Found file:", path)
			}

			return nil
		})

		if err != nil {
			log.Println("Recurse filepath error:", err)
			return files, err
		}
	} else {
		infos, err := os.ReadDir(path)
		if err != nil {
			log.Println("Read dir error:", err)
			return files, err
		}

		for _, info := range infos {
			if !info.IsDir() {
				files = append(files, filepath.Join(path, info.Name()))
				log.Println("Found file:", filepath.Join(path, info.Name()))
			}
		}
	}

	log.Println("All files found:", files)
	return files, nil
}
