package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-utils",
	Short: "A CLI for various git utils",
	Long: `git-utils v0.1
Copyright (c) Arbaaz Laskar <arzkar.dev@gmail.com>

A CLI for various git utilities
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
