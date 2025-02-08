package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/spf13/cobra"
	"strings"
)

var AnalyzeSpaceCmd = &cobra.Command{
	Use:   "analyze-space [directory]",
	Short: "Analyze disk space usage in the specified directory",
	Long:  `This command analyzes disk space usage and shows the largest files.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		top, _ := cmd.Flags().GetInt("top")
		ignorePattern, _ := cmd.Flags().GetString("ignore")
		ignoreList := strings.Split(ignorePattern, ",")

		files, err := filesystem.AnalyzeSpace(directory, top, ignoreList)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("No find files.")
		} else {
			fmt.Println("Top files by size:")
			for _, file := range files {
				fmt.Printf("%s (%d bytes)\n", file.Path, file.Size)
			}
		}
	},
}

func init() {
	AnalyzeSpaceCmd.Flags().IntP("top", "t", 10, "Number of files to display")
	AnalyzeSpaceCmd.Flags().StringP("ignore", "i", "", "Comma-separated list of directories or patterns to ignore (e.g., temp,.git)")
}
