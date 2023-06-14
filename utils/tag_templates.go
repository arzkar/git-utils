/*
Copyright 2023 Arbaaz Laskar

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateTemplateVariables(dir, prevTag, newTag, message string) map[string]string {
	templateVariables := map[string]string{
		"repo_owner": getRepositoryOwner(dir),
		"repo_name":  getRepositoryName(dir),
		"prevTag":    prevTag,
		"newTag":     newTag,
	}
	return templateVariables
}

func getRepositoryOwner(dir string) string {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get remote URL:", err)
		os.Exit(1)
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")

	remote := parts[len(parts)-2]
	remote = strings.TrimSuffix(remote, ".git")

	if len(parts) == 2 {
		remoteParts := strings.Split(remote, ":")
		if len(remoteParts) != 2 {
			fmt.Println("Invalid remote URL format:", remote)
			os.Exit(1)
		}
		return strings.TrimPrefix(remoteParts[1], "git@")
	} else if len(parts) == 5 {
		parts = strings.Split(remote, ":")
		username := parts[0]
		return username

	} else {
		fmt.Println("Invalid remote URL:", url)
		os.Exit(1)
	}
	return ""
}

func getRepositoryName(dir string) string {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get remote URL:", err)
		os.Exit(1)
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")

	remote := parts[len(parts)-1]
	return strings.TrimSuffix(remote, ".git")
}
