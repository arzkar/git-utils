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
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Displays git commit activity for the last working day",
	Run:   runStandup,
}

func init() {
	rootCmd.AddCommand(standupCmd)
}

func runStandup(cmd *cobra.Command, args []string) {
	fetchLatestCommits()

	commits := getCommits()

	if len(commits) == 0 {
		fmt.Println("No activity found.")
		return
	}

	for _, commit := range commits {
		fmt.Println(formatCommit(commit))
	}
}

func fetchLatestCommits() {
	cmd := execCommand("git", "fetch", "--all")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getCommits() []string {
	sinceDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	cmd := execCommand("git", "log", "--since="+sinceDate)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(output), "\ncommit ")
}

func formatCommit(commit string) string {
	commit = strings.TrimSpace(commit)
	lines := strings.Split(commit, "\n")
	if len(lines) < 5 {
		return ""
	}
	commitHash := lines[0]
	commitAuthor := extractCommitInfo(lines[4], "Author:")
	commitTime := extractCommitInfo(lines[3], "Date:")
	commitHeading := extractCommitHeading(lines[5:])
	return fmt.Sprintf("%s - %s (%s) <%s>", commitHash, commitHeading, commitTime, commitAuthor)
}

func extractCommitInfo(line, prefix string) string {
	startIndex := strings.Index(line, prefix)
	if startIndex == -1 {
		return ""
	}
	return strings.TrimSpace(line[startIndex+len(prefix):])
}

func extractCommitHeading(lines []string) string {
	var headingLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") || strings.HasPrefix(line, "Merge: ") ||
			strings.HasPrefix(line, "Author: ") || strings.HasPrefix(line, "Date: ") {
			break
		}
		headingLines = append(headingLines, line)
	}
	return strings.Join(headingLines, " ")
}

func execCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	return cmd
}
