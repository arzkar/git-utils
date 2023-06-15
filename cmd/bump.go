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

	"github.com/arzkar/git-utils/utils"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Version bump the version",
	Run:   bumpVersion,
}

var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Version bump the major version",
	Run:   bumpMajor,
}

var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Version bump the minor version",
	Run:   bumpMinor,
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Version bump the patch version",
	Run:   bumpPatch,
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.AddCommand(majorCmd)
	bumpCmd.AddCommand(minorCmd)
	bumpCmd.AddCommand(patchCmd)
}

func bumpVersion(cmd *cobra.Command, args []string) {
	currentVersion, err := utils.GetCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := utils.IncrementVersion(currentVersion)

	err = utils.UpdateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := utils.GetCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}
	commitMessage := fmt.Sprintf("Bump version: %s → %s", currentVersion, newVersion)
	if commitEnabled {
		err = utils.CommitChanges(currentVersion, newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := utils.GetTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = utils.CreateTag(newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpMajor(cmd *cobra.Command, args []string) {
	currentVersion, err := utils.GetCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := utils.BumpMajorVersion(currentVersion)

	err = utils.UpdateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := utils.GetCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	commitMessage := fmt.Sprintf("Bump version: %s → %s", currentVersion, newVersion)
	if commitEnabled {
		err = utils.CommitChanges(currentVersion, newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := utils.GetTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = utils.CreateTag(newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpMinor(cmd *cobra.Command, args []string) {
	currentVersion, err := utils.GetCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := utils.BumpMinorVersion(currentVersion)

	err = utils.UpdateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := utils.GetCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	commitMessage := fmt.Sprintf("Bump version: %s → %s", currentVersion, newVersion)
	if commitEnabled {
		err = utils.CommitChanges(currentVersion, newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := utils.GetTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = utils.CreateTag(newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}

func bumpPatch(cmd *cobra.Command, args []string) {
	currentVersion, err := utils.GetCurrentVersion()
	if err != nil {
		fmt.Println("Failed to read current version:", err)
		return
	}

	newVersion := utils.BumpPatchVersion(currentVersion)

	err = utils.UpdateFiles(currentVersion, newVersion)
	if err != nil {
		fmt.Println("Failed to update files:", err)
		return
	}

	commitEnabled, err := utils.GetCommitOption()
	if err != nil {
		fmt.Println("Failed to read commit option:", err)
		return
	}

	commitMessage := fmt.Sprintf("Bump version: %s → %s", currentVersion, newVersion)
	if commitEnabled {
		err = utils.CommitChanges(currentVersion, newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to commit changes:", err)
			return
		}
	}

	tagEnabled, err := utils.GetTagOption()
	if err != nil {
		fmt.Println("Failed to read tag option:", err)
		return
	}

	if tagEnabled {
		err = utils.CreateTag(newVersion, commitMessage)
		if err != nil {
			fmt.Println("Failed to create tag:", err)
			return
		}
	}

	fmt.Printf("Bump version: %s → %s\n", currentVersion, newVersion)
}
