package cmd

import (
	"fmt"
	"github.com/SHCDevelops/file-manager/internal/filesystem"
	"github.com/spf13/cobra"
	"strings"
)

var CodeStatsCmd = &cobra.Command{
	Use:   "code-stats [directory]",
	Short: "Analyze code statistics for supported languages",
	Long: `This command analyzes code statistics including:
- Total lines of code
- Comment lines
- Code lines (total - comments)

Supported languages: Go, HTML, CSS, JavaScript, TypeScript`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		ignorePattern, _ := cmd.Flags().GetString("ignore")
		ignoreList := strings.Split(ignorePattern, ",")

		stats, err := filesystem.CountCodeLines(directory, ignoreList)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(stats.Languages) == 0 {
			fmt.Println("No code files found in supported formats")
			return
		}

		fmt.Println("Code Statistics:")
		for lang, data := range stats.Languages {
			fmt.Printf("\n%s:\n", lang)
			fmt.Printf("  Total lines: %d\n", data.TotalLines)
			fmt.Printf("  Comments:    %d (%.1f%%)\n",
				data.CommentLines,
				percent(data.CommentLines, data.TotalLines))
			fmt.Printf("  Code lines:  %d (%.1f%%)\n",
				data.CodeLines,
				percent(data.CodeLines, data.TotalLines))
		}
	},
}

func init() {
	CodeStatsCmd.Flags().StringP("ignore", "i", "", "Comma-separated list of directories or patterns to ignore")
}

func percent(part, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(part) * 100.0 / float64(total)
}
