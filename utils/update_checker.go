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
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	repoOwner      = "arzkar"
	repoName       = "git-utils"
	releasesAPI    = "https://api.github.com/repos/%s/%s/releases/latest"
	currentVersion = "v0.4.0"
	cacheFile      = "git-utils-cache.json"
)

type release struct {
	TagName string `json:"tag_name"`
}

func init() {
	UpdateChecker()
}

func UpdateChecker() {
	// Read the cached version and publication time
	cachedVersion, err := readCachedVersion()
	if err != nil {
		fmt.Println("Failed to read cached version:", err)
	}

	latestVersion := ""
	// Check if the cached version is up-to-date
	if cachedVersion != "" && compareVersions(cachedVersion, currentVersion) >= 0 {
		latestVersion = cachedVersion
	} else {
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
	}

	if compareVersions(latestVersion, currentVersion) > 0 {
		fmt.Printf(color.RedString("A newer version (%s) of the CLI is available. Please update to the latest version.")+color.GreenString("\nhttps://github.com/arzkar/git-utils#installation\n"), latestVersion)
	}

	// Cache the latest version and publication time
	err = cacheVersion(latestVersion, time.Now())
	if err != nil {
		fmt.Println("Failed to cache the version:", err)
	}
}

func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	return strings.Compare(v1, v2)
}

func getCacheFilePath() string {
	// Get the cache directory path
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	// Ensure the cache directory exists
	err = os.MkdirAll(filepath.Join(cacheDir, "git-utils"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Construct the cache file path
	cachePath := filepath.Join(cacheDir, "git-utils", cacheFile)

	return cachePath
}

func readCachedVersion() (string, error) {
	cachePath := getCacheFilePath()

	file, err := os.Open(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Cache file does not exist, return empty version
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	cachedData := struct {
		Version string    `json:"version"`
		Time    time.Time `json:"time"`
	}{}

	err = json.NewDecoder(file).Decode(&cachedData)
	if err != nil {
		return "", err
	}

	return cachedData.Version, nil
}

func cacheVersion(version string, published time.Time) error {
	cachePath := getCacheFilePath()

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create the cache data
	cachedData := struct {
		Version string    `json:"version"`
		Time    time.Time `json:"time"`
	}{
		Version: version,
		Time:    published,
	}

	// Encode the cache data to JSON and write it to the file
	err = json.NewEncoder(file).Encode(&cachedData)
	if err != nil {
		return err
	}

	return nil
}
