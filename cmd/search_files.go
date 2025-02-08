package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

var SearchCmd = &cobra.Command{
	Use:   "search [pattern] [directory]",
	Short: "Search for files matching a pattern in the specified directory",
	Long: `This command searches for files that match a given pattern in the specified directory.
You can ignore specific directories using the --ignore flag.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]
		directory := args[1]

		ignorePattern, _ := cmd.Flags().GetString("ignore")
		ignoreList := strings.Split(ignorePattern, ",")

		matchedFiles, err := filesystem.SearchFiles(directory, pattern, ignoreList)

		if err != nil {
			color.Red("Error: %v\n", err)
			return
		}

		if len(matchedFiles) == 0 {
			color.Yellow("No files found.")
		} else {
			header := color.New(color.FgHiGreen, color.Bold).SprintFunc()
			fileColor := color.New(color.FgHiWhite).SprintFunc()

			fmt.Printf("\n%s\n", header("Matching files:"))
			for _, file := range matchedFiles {
				fmt.Printf("▸ %s\n", fileColor(file))
			}
		}
	},
}

func init() {
	SearchCmd.Flags().StringP("ignore", "i", "", "Comma-separated list of directories or patterns to ignore (e.g., temp,.git)")
}
