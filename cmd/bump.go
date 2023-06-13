package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the version",
	Run:   bumpVersion,
}

var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Bump the major version",
	Run:   bumpMajor,
}

var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Bump the minor version",
	Run:   bumpMinor,
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Bump the patch version",
	Run:   bumpPatch,
}

var bump_cfg = ".git-utils-bump.cfg"

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.AddCommand(majorCmd)
	bumpCmd.AddCommand(minorCmd)
	bumpCmd.AddCommand(patchCmd)
}

func bumpVersion(cmd *cobra.Command, args []string) {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := incrementVersion(currentVersion)

	err = updateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := getCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	if commitEnabled {
		err = commitChanges(currentVersion, newVersion)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := getTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = createTag(newVersion)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpMajor(cmd *cobra.Command, args []string) {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := bumpMajorVersion(currentVersion)

	err = updateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := getCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	if commitEnabled {
		err = commitChanges(currentVersion, newVersion)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := getTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = createTag(newVersion)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpMinor(cmd *cobra.Command, args []string) {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := bumpMinorVersion(currentVersion)

	err = updateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := getCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	if commitEnabled {
		err = commitChanges(currentVersion, newVersion)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := getTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = createTag(newVersion)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpPatch(cmd *cobra.Command, args []string) {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := bumpPatchVersion(currentVersion)

	err = updateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := getCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	if commitEnabled {
		err = commitChanges(currentVersion, newVersion)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := getTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = createTag(newVersion)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func getCurrentVersion() (string, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return "", err
	}

	return config.Section("bumpversion").Key("current_version").String(), nil
}

func incrementVersion(version string) string {
	parts := strings.Split(version, ".")
	lastPart := parts[len(parts)-1]
	lastPart = strings.TrimPrefix(lastPart, "v")
	updatedLastPart := strconv.Itoa(parseInt(lastPart) + 1)
	parts[len(parts)-1] = updatedLastPart
	return strings.Join(parts, ".")
}

func bumpMajorVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[0] = strconv.Itoa(parseInt(parts[0]) + 1)
	parts[1] = "0"
	parts[2] = "0"
	return strings.Join(parts, ".")
}

func bumpMinorVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[1] = strconv.Itoa(parseInt(parts[1]) + 1)
	parts[2] = "0"
	return strings.Join(parts, ".")
}

func bumpPatchVersion(version string) string {
	parts := strings.Split(version, ".")
	parts[2] = strconv.Itoa(parseInt(parts[2]) + 1)
	return strings.Join(parts, ".")
}

func updateFiles(currentVersion, newVersion string) error {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return err
	}

	config.Section("bumpversion").Key("current_version").SetValue(newVersion)

	err = config.SaveTo(bump_cfg)
	if err != nil {
		return err
	}

	// Update other files with new version if needed

	return nil
}

func getCommitOption() (bool, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return false, err
	}

	return config.Section("bumpversion").Key("commit").Bool()
}

func getTagOption() (bool, error) {
	config, err := ini.Load(bump_cfg)
	if err != nil {
		return false, err
	}

	return config.Section("bumpversion").Key("tag").Bool()
}

func commitChanges(currentVersion, newVersion string) error {
	cmd := exec.Command("git", "commit", "-am", fmt.Sprintf("Bump version: %s → %s", currentVersion, newVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createTag(version string) error {
	cmd := exec.Command("git", "tag", "-a", version, "-m", fmt.Sprintf("Version %s", version))
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
