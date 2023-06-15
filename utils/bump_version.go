package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

var bump_cfg = ".git-utils-bump.cfg"

func GetCurrentVersion() (string, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return "", err
	}

	return config.Section("bumpversion").Key("current_version").String(), nil
}

func IncrementVersion(version string) string {
	parts := strings.Split(version, ".")
	lastPart := parts[len(parts)-1]
	lastPart = strings.TrimPrefix(lastPart, "v")
	updatedLastPart := strconv.Itoa(parseInt(lastPart) + 1)
	parts[len(parts)-1] = updatedLastPart
	return strings.Join(parts, ".")
}

func BumpMajorVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[0] = strconv.Itoa(parseInt(parts[0]) + 1)
	parts[1] = "0"
	parts[2] = "0"
	return strings.Join(parts, ".")
}

func BumpMinorVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[1] = strconv.Itoa(parseInt(parts[1]) + 1)
	parts[2] = "0"
	return strings.Join(parts, ".")
}

func BumpPatchVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[2] = strconv.Itoa(parseInt(parts[2]) + 1)
	return strings.Join(parts, ".")
}

func UpdateFiles(currentVersion, newVersion string) error {

	// Check if the Git directory is dirty
	if err := checkGitDirectoryStatus(); err != nil {
		fmt.Println("Git directory is dirty.\nPlease stage, commit, or stash your changes before running the bump.")
		os.Exit(1)
	}
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return err
	}

	sections := config.SectionStrings()

	// Check if search pattern exists in all the files
	for _, section := range sections {
		if strings.HasPrefix(section, "bumpversion:file:") {
			filePath := strings.TrimPrefix(section, "bumpversion:file:")
			searchPattern := config.Section(section).Key("search").String()

			// Check if the search pattern exists in the file
			if !fileContainsString(filePath, strings.ReplaceAll(searchPattern, "{current_version}", currentVersion)) {
				return fmt.Errorf("search pattern not found in file: %s", filePath)
			}
		}
	}
	for _, section := range sections {
		if strings.HasPrefix(section, "bumpversion:file:") {
			filePath := strings.TrimPrefix(section, "bumpversion:file:")
			searchPattern := config.Section(section).Key("search").String()
			replacePattern := config.Section(section).Key("replace").String()

			err := replaceStringInFile(filePath, searchPattern, replacePattern, currentVersion, newVersion)
			if err != nil {
				return err
			}

			// Stage the updated file
			cmd := exec.Command("git", "add", filePath)
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to stage file: %s", filePath)
			}
		}
	}
	config.Section("bumpversion").Key("current_version").SetValue(newVersion)
	err = config.SaveTo(bump_cfg)
	if err != nil {
		return err
	}

	// Stage the updated file
	cmd := exec.Command("git", "add", bump_cfg)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stage file: %s", bump_cfg)
	}

	return nil
}

func fileContainsString(filename, searchPattern string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), searchPattern)
}

func replaceStringInFile(filename, searchPattern, replacePattern, currentVersion, newVersion string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Replace the search pattern with the new version
	newData := strings.ReplaceAll(string(data), strings.ReplaceAll(searchPattern, "{current_version}", currentVersion), strings.ReplaceAll(replacePattern, "{new_version}", newVersion))
	err = os.WriteFile(filename, []byte(newData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetCommitOption() (bool, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return false, err
	}

	return config.Section("bumpversion").Key("commit").Bool()
}

func GetTagOption() (bool, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return false, err
	}

	return config.Section("bumpversion").Key("tag").Bool()
}

func CommitChanges(currentVersion, newVersion, commitMessage string) error {
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateTag(version, message string) error {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return err
	}

	tagFormat := config.Section("bumpversion").Key("tag_format").String()
	tagName := strings.ReplaceAll(tagFormat, "{tag}", version)

	cmd := exec.Command("git", "tag", "-a", tagName, "-m", fmt.Sprintf("Version %s", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func checkGitDirectoryStatus() error {
	// Check if there are unstaged changes
	cmd := exec.Command("git", "diff", "--exit-code")
	err := cmd.Run()
	if err != nil {
		// If there are unstaged changes, return an error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("git directory has unstaged changes")
		}
		return err
	}

	return nil
}
