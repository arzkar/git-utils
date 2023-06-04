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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var grepCmd *cobra.Command

func init() {
	grepCmd = &cobra.Command{
		Use:   "grep",
		Short: "Search for a pattern in files",
		Long:  "Recursively search for a pattern in files within the specified directory",
		Args:  cobra.ExactArgs(1),
		Run:   runGrep,
	}

	var dir string
	var bufferSize int
	grepCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to search in")
	grepCmd.Flags().IntVarP(&bufferSize, "buffer", "b", 4096, "Buffer size for reading files (default: 4096)")
	rootCmd.AddCommand(grepCmd)
}

func runGrep(cmd *cobra.Command, args []string) {
	pattern := args[0]
	dir, _ := cmd.Flags().GetString("dir")
	bufferSize, _ := cmd.Flags().GetInt("buffer")

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

	foundMatch := false

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git folder
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			// Get the relative path from the repository root directory to the file
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			buf := make([]byte, bufferSize)
			scanner.Buffer(buf, bufferSize)
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, pattern) {
					foundMatch = true
					line = strings.ReplaceAll(line, pattern, color.RedString(pattern))
					fmt.Printf("\n%s\nL%d: %s\n", relPath, lineNumber, line)
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				if err.Error() == bufio.ErrTooLong.Error() {
					fmt.Println("Encountered buffer error. Please consider increasing the buffer size.")
					fmt.Println("To increase the buffer size, use --buffer buffer_size (default: 4096)")
					os.Exit(1)
				} else {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !foundMatch {
		fmt.Println("No matches found.")
	}
}
