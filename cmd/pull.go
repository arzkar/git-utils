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

var pullCmd *cobra.Command

func init() {
	pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull all or specified branches",
		Long:  "Pull all or comma-separated list of branches",
		Args:  cobra.ExactArgs(1),
		Run:   runPull,
	}

	var dryRun bool
	pullCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without actually pulling the changes")
	pullCmd.Flags().StringP("dir", "d", "", "Directory to perform the pull operation")

	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) {
	pull := args[0]
	dir, _ := cmd.Flags().GetString("dir")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

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
			err := pullRepository(path, pull, dryRun)
			if err != nil {
				fmt.Printf("Error pulling repository '%s': %s\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func pullRepository(path string, pull string, dryRun bool) error {
	if pull == "all" {
		err := pullAllBranches(path, dryRun)
		if err != nil {
			return err
		}
	} else {
		branches := strings.Split(pull, ",")
		for _, branch := range branches {
			err := pullBranch(path, branch, dryRun)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pullAllBranches(path string, dryRun bool) error {
	cmd := exec.Command("git", "-C", path, "branch", "--format", "%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get local branches: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, branch := range branches {
		err := pullBranch(path, branch, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func pullBranch(path string, branch string, dryRun bool) error {
	fmt.Printf("Pulling branch '%s' in repository '%s'\n", branch, path)
	execDryRun(dryRun, path, branch)

	// Show the changes made by the pull operation
	cmd := exec.Command("git", "-C", path, "diff", "--stat", branch+"..origin/"+branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get changes made by pull for branch '%s' in repository '%s': %w\n%s", branch, path, err, string(output))
	}

	if len(output) > 0 {
		fmt.Println("Changes made by pull:")
		colorizedOutput := utils.ColorizeDiffStat(string(output))
		fmt.Println(colorizedOutput)

		cmd = exec.Command("git", "-C", path, "pull")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to pull branch '%s' in repository '%s': %w\n%s", branch, path, err, string(output))
		}

		fmt.Printf(color.GreenString("Successfully pulled branch '%s' in repository '%s'\n\n", branch, path))
	} else {
		fmt.Println(color.GreenString("No changes made by pull\n"))
	}

	return nil
}

func execDryRun(dryRun bool, path string, branch string) {
	if dryRun {
		fmt.Printf("Dry run: Changes for branch '%s' in repository '%s':\n", branch, path)

		cmd := exec.Command("git", "-C", path, "fetch", "origin", branch)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("failed to fetch branch '%s' in repository '%s': %s\n%s", branch, path, err, string(output))
			os.Exit(1)
		}

		cmd = exec.Command("git", "-C", path, "diff", "--stat", branch+"..origin/"+branch)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("failed to get changes for branch '%s' in repository '%s': %s\n%s", branch, path, err, string(output))
			os.Exit(1)
		}

		if len(output) > 0 {
			colorizedOutput := utils.ColorizeDiffStat(string(output))
			fmt.Println(colorizedOutput)
		} else {
			fmt.Println("No changes")
		}

		os.Exit(0)
	}
}
