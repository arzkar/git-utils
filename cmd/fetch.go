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

	"github.com/arzkar/git-utils/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var fetchCmd *cobra.Command

func init() {
	fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch all or specified branches",
		Long:  "Fetch all or comma-separated list of branches",
		Args:  cobra.ExactArgs(1),
		Run:   runFetch,
	}

	fetchCmd.Flags().StringP("dir", "d", "", "Directory to perform the fetch operation")

	rootCmd.AddCommand(fetchCmd)
}

func runFetch(cmd *cobra.Command, args []string) {
	fetch := args[0]
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

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && utils.IsGitRepository(path) {
			err := fetchRepository(path, fetch)
			if err != nil {
				fmt.Printf("Error fetching repository '%s': %s\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func fetchRepository(path string, fetch string) error {
	if fetch == "all" {
		err := fetchAllBranches(path)
		if err != nil {
			return err
		}
	} else {
		branches := strings.Split(fetch, ",")
		for _, branch := range branches {
			err := fetchBranch(path, branch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fetchAllBranches(path string) error {
	cmd := exec.Command("git", "-C", path, "fetch", "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch all branches in repository '%s': %s\n%s", path, err, string(output))
	}

	fmt.Printf(color.GreenString("Successfully fetched all branches in repository '%s'\n\n", path))
	return nil
}

func fetchBranch(path string, branch string) error {
	fmt.Printf("Fetching branch '%s' in repository '%s'\n", branch, path)
	cmd := exec.Command("git", "-C", path, "fetch", "origin", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch branch '%s' in repository '%s': %s\n%s", branch, path, err, string(output))
	}

	fmt.Printf(color.GreenString("Successfully fetched branch '%s' in repository '%s'\n\n", branch, path))
	return nil
}
