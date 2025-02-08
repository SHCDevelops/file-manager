package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "file-manager",
		Short: "CLI tool for managing and analyzing files",
		Long:  `File Manager is a powerful CLI tool to analyze and manage files and directories.`,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
