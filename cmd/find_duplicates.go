package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/spf13/cobra"
)

var FindDuplicatesCmd = &cobra.Command{
	Use:   "find-duplicates [directory]",
	Short: "Find duplicate files in the specified directory",
	Long: `This command scans the specified directory and finds duplicate files based on their content.
It uses file hashes to identify duplicates.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		duplicates, err := filesystem.FindDuplicates(directory)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(duplicates) == 0 {
			fmt.Println("No duplicates found.")
		} else {
			fmt.Println("Duplicates found:")
			for _, group := range duplicates {
				fmt.Println(group)
			}
		}
	},
}
