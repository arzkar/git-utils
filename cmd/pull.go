package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull all or specified branches",
	Long:  "Pull all or comma-separated list of branches",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pull := args[0]
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

			if info.IsDir() && isGitRepository(path) {
				err := pullRepository(path, pull)
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
	},
}

func isGitRepository(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func pullRepository(path string, pull string) error {
	if pull == "all" {
		err := pullAllBranches(path)
		if err != nil {
			return err
		}
	} else {
		branches := strings.Split(pull, ",")
		for _, branch := range branches {
			err := pullBranch(path, branch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pullAllBranches(path string) error {
	cmd := exec.Command("git", "-C", path, "branch", "--format", "%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get local branches: %s", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, branch := range branches {
		err := pullBranch(path, branch)
		if err != nil {
			return err
		}
	}

	return nil
}

func pullBranch(path, branch string) error {
	log.Printf("Pulling branch '%s' in repository '%s'\n", branch, path)

	cmd := exec.Command("git", "-C", path, "pull", "origin", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull branch '%s': %s", branch, string(output))
	}

	log.Printf("Branch '%s' pulled successfully in repository '%s'\n", branch, path)

	return nil
}

func init() {
	var dir string
	pullCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to perform the pull operation")
}
