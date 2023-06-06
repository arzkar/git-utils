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
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var grepCmd *cobra.Command

func init() {
	grepCmd = &cobra.Command{
		Use:   "grep",
		Short: "Search for a pattern in files",
		Long:  "Recursively search for a pattern in files within the specified directory",
		Args:  cobra.ExactArgs(1),
		Run:   runGrep,
	}

	var dir string
	grepCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to search in")
	rootCmd.AddCommand(grepCmd)
}

func runGrep(cmd *cobra.Command, args []string) {
	pattern := args[0]
	dir, _ := cmd.Flags().GetString("dir")

	if dir == "" {
		// Use current working directory if --dir flag is not specified
		dir, _ = os.Getwd()
	} else {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			fmt.Printf("Directory '%s' does not exist\n", dir)
			os.Exit(1)
		}
	}

	foundMatch := false

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process directories
		if !info.IsDir() {
			return nil
		}

		// Check if the directory contains a .git subdirectory
		gitDir := filepath.Join(path, ".git")
		_, err = os.Stat(gitDir)
		if os.IsNotExist(err) {
			return nil
		}

		// Execute git grep in the git repository
		cmd := exec.Command("git", "grep", "-n", pattern)
		cmd.Dir = path
		output, err := cmd.Output()
		if err != nil {
			// Ignore "exit status 1" error when no matches are found
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				return nil
			}
			return err
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				foundMatch = true
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					coloredLine := strings.ReplaceAll(parts[1], pattern, color.RedString(pattern))
					fmt.Printf("\n%s:\nL%s\n", path, coloredLine)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !foundMatch {
		fmt.Println("No matches found.")
	}
}
