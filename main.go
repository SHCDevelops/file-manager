package main

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "file-manager",
		Short: "CLI tool for managing and analyzing files",
		Long:  `File Manager is a powerful CLI tool to analyze and manage files and directories.`,
	}

	rootCmd.AddCommand(cmd.AnalyzeSpaceCmd)
	rootCmd.AddCommand(cmd.FindDuplicatesCmd)
	rootCmd.AddCommand(cmd.SearchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
