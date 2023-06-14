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
	"path/filepath"
)

func init() {
	initAppDirectory()
}

func initAppDirectory() {
	appDir := GetAppDir()

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create the app directory:", err)
			os.Exit(1)
		}

		fmt.Println("App directory created:", appDir)
	}
}

func GetAppDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get the user's home directory:", err)
		os.Exit(1)
	}

	appDirName := filepath.Join("arzkar", "git-utils")
	appDir := filepath.Join(homeDir, appDirName)

	return appDir
}

func GetConfigFilePath() string {
	appDir := GetAppDir()
	return filepath.Join(appDir, "config.json")
}
