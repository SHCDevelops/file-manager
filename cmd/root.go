package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "file-manager",
	Short: "CLI tool for managing and analyzing files",
	Long:  `File Manager is a powerful CLI tool to analyze and manage files and directories.`,
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
