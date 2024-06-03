package nexus_download

import (
	"fmt"
	"testing"
)

func TestUtil_findRemoteFiles(t *testing.T) {
	files, err := findRemoteFiles("https://nexus.mobilesolutionworks.com/repository/artifacts/ephemeral-bamboo-agent/framework-bundles/atlassian-plugins-osgi-bridge-7.5.3.jar", "gh", "gh")
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", files)
}

func TestUtil_findRemoteFiles2(t *testing.T) {
	files, err := findRemoteFiles("https://nexus.mobilesolutionworks.com/repository/artifacts/ephemeral-bamboo-agent/", "gh", "gh")
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", files)
}
