package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var tagCmd *cobra.Command

func init() {
	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Create a new tag for the repository",
		Run:   runTag,
	}

	var tagName, tagMessage, dir string
	tagCmd.Flags().StringVarP(&tagName, "tag_name", "a", "", "Name of the tag")
	tagCmd.Flags().StringVarP(&tagMessage, "tag_message", "m", "", "Message for the tag")
	tagCmd.Flags().StringVarP(&dir, "dir", "", "", "Git directory (optional)")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) {
	tagName, _ := cmd.Flags().GetString("tag_name")
	tagMessage, _ := cmd.Flags().GetString("tag_message")
	dir, _ := cmd.Flags().GetString("dir")

	if tagName == "" || tagMessage == "" {
		tagCmd.Help()
		return
	}

	fmt.Println("dir", dir)
	// Check if the tag message contains the @changelog keyword
	if strings.Contains(tagMessage, "@changelog") {
		// Retrieve the previous and new tags
		prevTag, _ := getPreviousTag(dir)
		newTag := tagName

		// Construct the changelog URL
		changelogURL := fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", getUsername(dir), getRepositoryName(dir), prevTag, newTag)

		// Construct the full changelog message
		fullTagMessage := fmt.Sprintf("Full changelog: %s", changelogURL)

		// Open the default Git editor for tag message editing
		editedTagMessage := openGitEditor(fullTagMessage)
		if editedTagMessage != "" {
			tagMessage = editedTagMessage
		}
	}

	// Create the tag using the git command in the specified directory
	fmt.Println("tagName", tagName)
	fmt.Println("tagMessage", tagMessage)
	// cmdGit := exec.Command("git", "-C", dir, "tag", tagName, "-a", "-m", tagMessage)
	// cmdGit.Stdout = os.Stdout
	// cmdGit.Stderr = os.Stderr
	// err := cmdGit.Run()
	// if err != nil {
	// 	fmt.Println("Failed to create tag:", err)
	// 	return
	// }

	// fmt.Println("Tag created successfully. Push it by running: git push --tags")
}

func getPreviousTag(dir string) (string, error) {
	cmdGit := exec.Command("git", "-C", dir, "describe", "--abbrev=0", "--tags", "--exclude=*-*")
	output, err := cmdGit.Output()
	if err != nil {
		return "", err
	}

	prevTag := strings.TrimSpace(string(output))
	return prevTag, nil
}

func getUsername(dir string) string {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get remote URL:", err)
		os.Exit(1)
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")

	remote := parts[len(parts)-2]
	remote = strings.TrimSuffix(remote, ".git")

	if len(parts) == 2 {
		remoteParts := strings.Split(remote, ":")
		if len(remoteParts) != 2 {
			fmt.Println("Invalid remote URL format:", remote)
			os.Exit(1)
		}
		return strings.TrimPrefix(remoteParts[1], "git@")
	} else if len(parts) == 5 {
		parts = strings.Split(remote, ":")
		username := parts[0]
		return username

	} else {
		fmt.Println("Invalid remote URL:", url)
		os.Exit(1)
	}
	return ""
}

func getRepositoryName(dir string) string {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get remote URL:", err)
		os.Exit(1)
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")

	remote := parts[len(parts)-1]
	return strings.TrimSuffix(remote, ".git")
}

func openGitEditor(initialContent string) string {
	editor := getGitEditor()
	if editor == "" {
		fmt.Println("No default Git editor found.")
		setEditorOption := promptYesNo("Do you want to set the Git core editor now and continue editing or commit the default message? (y/n): ")
		if setEditorOption == "y" {
			setGitEditor()
			editor = getGitEditor()
		} else {
			return initialContent
		}
	}

	tmpFile, err := os.CreateTemp("", "tagmessage")
	if err != nil {
		fmt.Println("Failed to create temporary file for editing tag message:", err)
		return initialContent
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(initialContent)); err != nil {
		fmt.Println("Failed to write initial content to temporary file:", err)
		return initialContent
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to open tag message editor:", err)
		return initialContent
	}

	editedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println("Failed to read edited tag message:", err)
		return initialContent
	}

	return strings.TrimSpace(string(editedContent))
}

func getGitEditor() string {
	cmd := exec.Command("git", "config", "--get", "core.editor")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func setGitEditor() {
	fmt.Println("Please enter the command for the Git core editor:")
	editor := promptInput("Editor command: ")

	cmd := exec.Command("git", "config", "--global", "core.editor", editor)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to set Git core editor:", err)
		return
	}

	fmt.Println("Git core editor set successfully.")
}

func promptYesNo(prompt string) string {
	var response string
	for {
		fmt.Print(prompt)
		fmt.Scanln(&response)
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "n" {
			break
		}
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}
	return response
}

func promptInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}
