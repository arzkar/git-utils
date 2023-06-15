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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	repoOwner      = "arzkar"
	repoName       = "git-utils"
	releasesAPI    = "https://api.github.com/repos/%s/%s/releases/latest"
	currentVersion = "v0.5.0"
	configFile     = "config.json"
	updateInterval = 24 * time.Hour
)

type release struct {
	TagName string `json:"tag_name"`
}

func init() {
	UpdateChecker()
}

func UpdateChecker() {
	// Read the cached version and publication time
	cachedConfig, err := ReadConfigFile()
	if err != nil {
		fmt.Println("Failed to read cached config:", err)
	}

	// Check if the cached version is up-to-date
	if time.Since(cachedConfig.LastUpdated) < updateInterval {
		return
	}

	latestVersion := ""
	url := fmt.Sprintf(releasesAPI, repoOwner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to check for new version:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var rel release
	err = json.Unmarshal(body, &rel)
	if err != nil {
		fmt.Println("Failed to parse response body:", err)
		return
	}

	latestVersion = rel.TagName

	if compareVersions(latestVersion, currentVersion) > 0 {
		fmt.Printf(color.RedString("A newer version (%s) of the CLI is available. Please update to the latest version.")+color.GreenString("\nhttps://github.com/arzkar/git-utils#installation\n"), latestVersion)
	}

	// Update the latest version and publication time in the config file
	err = UpdateConfig(func(config *Config) {
		config.Version = latestVersion
		config.LastUpdated = time.Now()
	})
	if err != nil {
		fmt.Println("Failed to update the config:", err)
	}
}

func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	return strings.Compare(v1, v2)
}

func ParseTemplate(template string, variables map[string]string) string {
	for key, value := range variables {
		template = strings.ReplaceAll(template, "{"+key+"}", value)
	}
	return template
}

func ReadConfigFile() (Config, error) {
	config := Config{}
	filePath := GetConfigFilePath()

	data, err := os.ReadFile(filePath)
	if err != nil {
		// If the file doesn't exist, create the default config file
		if os.IsNotExist(err) {
			err = createDefaultConfigFile()
			if err != nil {
				return config, err
			}
			// Read the newly created config file
			data, err = os.ReadFile(filePath)
			if err != nil {
				return config, err
			}
		} else {
			return config, err
		}
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func createDefaultConfigFile() error {
	config := Config{
		Tags: struct {
			Messages map[string]string `json:"messages"`
		}{
			Messages: make(map[string]string),
		},
	}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	filePath := GetConfigFilePath()
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UpdateConfig(updateFunc func(config *Config)) error {
	filePath := GetConfigFilePath()

	// Read the existing config file
	config, err := ReadConfigFile()
	if err != nil {
		return err
	}

	// Call the update function to modify the config
	updateFunc(&config)

	// Write the updated config back to the file
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
