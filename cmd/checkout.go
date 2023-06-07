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

	"github.com/arzkar/git-utils/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var checkoutCmd *cobra.Command

func init() {
	checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a branch in all repositories",
		Long:  "Checkout a branch in all repositories",
		Args:  cobra.ExactArgs(1),
		Run:   runCheckout,
	}

	checkoutCmd.Flags().StringP("dir", "d", "", "Directory to perform the checkout operation")

	rootCmd.AddCommand(checkoutCmd)
}

func runCheckout(cmd *cobra.Command, args []string) {
	branch := args[0]
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
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only process directories
		if !d.IsDir() {
			return nil
		}

		// Check if the directory contains a .git subdirectory
		gitDir := filepath.Join(path, ".git")
		_, err = os.Stat(gitDir)
		if os.IsNotExist(err) {
			return nil
		}

		if d.IsDir() && utils.IsGitRepository(path) {
			err := checkoutBranch(path, branch)
			if err != nil {
				fmt.Printf("Error checking out branch '%s' in repository '%s': %s\n", branch, path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func checkoutBranch(path string, branch string) error {
	fmt.Printf("Checking out branch '%s' in repository '%s'\n", branch, path)
	cmd := exec.Command("git", "-C", path, "checkout", "--track", fmt.Sprintf("origin/%s", branch))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s' in repository '%s': %s\n%s", branch, path, err, string(output))
	}

	fmt.Printf(color.GreenString("Successfully checked out branch '%s' in repository '%s'\n\n", branch, path))
	return nil
}
