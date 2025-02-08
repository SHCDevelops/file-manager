package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/spf13/cobra"
	"strings"
)

var FindDuplicatesCmd = &cobra.Command{
	Use:   "find-duplicates [directory]",
	Short: "Find duplicate files in the specified directory",
	Long: `This command scans the specified directory and finds duplicate files based on their content.
It uses file hashes to identify duplicates.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		ignorePattern, _ := cmd.Flags().GetString("ignore")
		ignoreList := strings.Split(ignorePattern, ",")

		duplicates, err := filesystem.FindDuplicates(directory, ignoreList)
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

func init() {
	FindDuplicatesCmd.Flags().StringP("ignore", "i", "", "Comma-separated list of directories or patterns to ignore (e.g., temp,.git)")
}
