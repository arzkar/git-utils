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
	cmd := exec.Command("git", "-C", path, "branch", "--list", fmt.Sprintf("%s*", branch))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list matching branches in repository '%s': %s\n%s", path, err, string(output))
	}

	branches := parseBranches(output, branch)
	if len(branches) == 0 {
		return fmt.Errorf("no matching branches found in repository '%s'", path)
	}

	// Display menu to choose between branches
	choice, err := displayBranchMenu(branches)
	if err != nil {
		return fmt.Errorf("failed to display branch menu: %s", err)
	}

	selectedBranch := branches[choice]
	cmd = exec.Command("git", "-C", path, "checkout", selectedBranch)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s' in repository '%s': %s\n%s", selectedBranch, path, err, string(output))
	}

	fmt.Printf(color.GreenString("Successfully checked out branch '%s' in repository '%s'\n\n", selectedBranch, path))
	return nil
}

// parseBranches parses the output of the "git branch --list" command
func parseBranches(output []byte, branch string) []string {
	var branches []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, fmt.Sprintf("origin/%s", branch)) {
			branches = append(branches, line)
		}
	}
	return branches
}

// displayBranchMenu displays a menu with the list of branches and returns the selected choice
func displayBranchMenu(branches []string) (int, error) {
	fmt.Println("Choose a branch:")
	for i, branch := range branches {
		fmt.Printf("[%d] %s\n", i+1, branch)
	}

	fmt.Print("Enter your choice: ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(branches) {
		return 0, fmt.Errorf("invalid choice")
	}

	return choice - 1, nil
}
