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

	"github.com/arzkar/git-utils/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-utils",
	Short: "A CLI for performing various operations on git repositories",
	Long: `git-utils v0.5.0
Copyright (c) Arbaaz Laskar <arzkar.dev@gmail.com>

A CLI for performing various operations on git repositories
`,
	Run: runRoot,
}

var configFlag bool

func init() {
	rootCmd.Flags().BoolVar(&configFlag, "config", false, "Show app config")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	if configFlag {
		// Print app directory and config file path
		appDir := utils.GetAppDir()
		configFilePath := utils.GetConfigFilePath()
		fmt.Println("App Directory:", appDir)
		fmt.Println("Config File Path:", configFilePath)
	} else {
		cmd.Help()
	}
}
