package nexus_download

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func downloadFile(targetURLPath string, localFilePath string) error {
	// Create the file
	out, err := os.Create(localFilePath)
	if err != nil {
		log.Println("Error creating the local file:", err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(targetURLPath)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200 OK, got " + resp.Status)
	}

	// Save the data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Error writing to local file:", err)
		return err
	}

	return nil
}

func parseNexusUrl(source string) (string, string, string, error) {
	u, err := url.Parse(source)
	if err != nil {
		return "", "", "", err
	}

	host := u.Scheme + "://" + u.Hostname()
	if u.Port() != "" {
		host = host + ":" + u.Port()
	}

	pathSegments := strings.Split(u.Path, "/")

	repositoryName := pathSegments[2]
	repositoryPath := strings.Join(pathSegments[3:], "/")

	return host, repositoryName, repositoryPath, nil
}

type AssetChecksum struct {
	MD5 string `json:"md5"`
}

type Asset struct {
	DownloadURL string        `json:"downloadUrl"`
	Path        string        `json:"path"`
	ID          string        `json:"id"`
	Repository  string        `json:"repository"`
	Format      string        `json:"format"`
	Checksum    AssetChecksum `json:"checksum"`
}

type Download struct {
	DownloadURL string
	MD5         string
}

func findRemoteFiles(sourceURLPath string, username string, password string) ([]Download, bool, error) {
	nexusUrl, repository, path, err := parseNexusUrl(sourceURLPath)
	if err != nil {
		return nil, false, err
	}

	var urls []Download
	var nextPageToken string

	for {
		searchUrl := fmt.Sprintf("%s/service/rest/v1/search/assets?q=%s&repository=%s", nexusUrl, path, repository)

		if nextPageToken != "" {
			searchUrl += "&continuationToken=" + nextPageToken
		}

		req, _ := http.NewRequest("GET", searchUrl, nil)
		req.SetBasicAuth(username, password)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, false, err
		}

		defer resp.Body.Close()

		assets := struct {
			Items             []Asset `json:"items"`
			ContinuationToken string  `json:"continuationToken"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&assets)
		if err != nil {
			return nil, false, err
		}

		for _, asset := range assets.Items {
			download := Download{
				DownloadURL: asset.DownloadURL,
				MD5:         asset.Checksum.MD5,
			}

			if asset.DownloadURL == sourceURLPath {
				urls = []Download{download}
				return urls, true, nil
			} else {
				urls = append(urls, download)
			}
		}

		if assets.ContinuationToken == "" {
			break
		} else {
			nextPageToken = assets.ContinuationToken
		}
	}

	return urls, false, nil
}
