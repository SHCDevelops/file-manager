package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/fatih/color"
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
			color.Red("Error: %v\n", err)
			return
		}

		if len(duplicates) == 0 {
			color.Green("No duplicates found. 🎉")
		} else {
			groupHeader := color.New(color.FgHiRed, color.Bold).SprintFunc()
			fileColor := color.New(color.FgHiYellow).SprintFunc()

			fmt.Printf("\n%s\n", groupHeader("Duplicates found:"))
			for i, group := range duplicates {
				fmt.Printf("\n%s %d\n", groupHeader("Group"), i+1)
				for _, file := range group {
					fmt.Printf("▸ %s\n", fileColor(file))
				}
			}
		}
	},
}

func init() {
	FindDuplicatesCmd.Flags().StringP("ignore", "i", "", "Comma-separated list of directories or patterns to ignore (e.g., temp,.git)")
}
