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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/arzkar/git-utils/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Config struct {
	Tags struct {
		Messages map[string]string `json:"messages"`
	} `json:"tags"`
}

var tagCmd *cobra.Command

func init() {
	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Create a new tag with custom message for the repository",
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

	// Read the config file
	config, err := readConfigFile()
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		fmt.Println("Set the message values in the config file.\nRun: git-utils --config", err)
		os.Exit(1)
	}

	// Check if the tag message matches a configured message
	for key, message := range config.Tags.Messages {
		if tagMessage == key {
			prevTag, _ := getPreviousTag(dir)
			newTag := tagName

			templateVariables := utils.CreateTemplateVariables(dir, prevTag, newTag, message)
			tagMessage = parseTemplate(message, templateVariables)

			// Create the tag using the git command in the specified directory
			cmdGit := exec.Command("git", "-C", dir, "tag", tagName, "-a", "-m", tagMessage)
			cmdGit.Stdout = os.Stdout
			cmdGit.Stderr = os.Stderr
			err = cmdGit.Run()
			if err != nil {
				fmt.Println(color.RedString("Failed to create tag:", err))
				return
			}

			fmt.Println(color.GreenString("Tag created successfully. Push it by running: git push --tags"))
			return
		}
	}

	fmt.Println(color.RedString("No messages has been set in the config file. Set it up before running the tag command.") + color.GreenString("\nRun: git-utils --config"))

}

func parseTemplate(template string, variables map[string]string) string {
	for key, value := range variables {
		template = strings.ReplaceAll(template, "{"+key+"}", value)
	}
	return template
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

func readConfigFile() (Config, error) {
	config := Config{}
	filePath := utils.GetConfigFilePath()

	data, err := os.ReadFile(filePath)
	if err != nil {
		// If the file doesn't exist, create the default config file
		if os.IsNotExist(err) {
			err = createDefaultConfigFile()
			if err != nil {
				return config, err
			}
			// Read the newly created config file
			data, err = os.ReadFile(filePath)
			if err != nil {
				return config, err
			}
		} else {
			return config, err
		}
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func createDefaultConfigFile() error {
	config := Config{
		Tags: struct {
			Messages map[string]string `json:"messages"`
		}{
			Messages: make(map[string]string),
		},
	}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	filePath := utils.GetConfigFilePath()
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
