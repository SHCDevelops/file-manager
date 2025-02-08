package filesystem

import (
	"bufio"
	"github.com/SHCDevelops/file-manager/lib/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CodeStats struct {
	Languages map[string]*LanguageStat
	mu        sync.Mutex
}

type LanguageStat struct {
	TotalLines   int
	CommentLines int
	CodeLines    int
}

func CountCodeLines(root string, ignoreList []string) (*CodeStats, error) {
	stats := &CodeStats{
		Languages: make(map[string]*LanguageStat),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	errWalk := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(root, path)

		if info.IsDir() {
			if utils.IsIgnored(relPath, ignoreList, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if utils.IsIgnored(relPath, ignoreList, false) {
			return nil
		}

		lang := getLanguage(path)
		if lang == "" {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			total, comments, err := analyzeFile(path, lang)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			stats.mu.Lock()
			defer stats.mu.Unlock()

			if _, exists := stats.Languages[lang]; !exists {
				stats.Languages[lang] = &LanguageStat{}
			}

			stats.Languages[lang].TotalLines += total
			stats.Languages[lang].CommentLines += comments
			stats.Languages[lang].CodeLines = stats.Languages[lang].TotalLines - stats.Languages[lang].CommentLines
		}()

		return nil
	})

	wg.Wait()

	if errWalk != nil {
		return nil, errWalk
	}

	mu.Lock()
	defer mu.Unlock()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return stats, nil
}

func getLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".html", ".htm":
		return "HTML"
	case ".css":
		return "CSS"
	case ".js":
		return "JavaScript"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".go":
		return "Go"
	default:
		return ""
	}
}

func analyzeFile(path, lang string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var parser LineParser
	switch lang {
	case "HTML":
		parser = &htmlParser{}
	case "CSS":
		parser = &cssParser{}
	case "JavaScript", "TypeScript":
		parser = &jsParser{}
	case "Go":
		parser = &goParser{}
	default:
		return 0, 0, nil
	}

	return parser.Parse(file)
}

type LineParser interface {
	Parse(*os.File) (total int, comments int, err error)
}

type htmlParser struct{}

func (p *htmlParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inComment := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inComment:
			comments++
			if strings.Contains(line, "-->") {
				inComment = false
			}
		case strings.Contains(line, "<!--"):
			comments++
			if !strings.Contains(line, "-->") {
				inComment = true
			}
		}
	}
	return total, comments, scanner.Err()
}

type cssParser struct{}

func (p *cssParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inComment := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inComment:
			comments++
			if strings.Contains(line, "*/") {
				inComment = false
			}
		case strings.HasPrefix(line, "/*"):
			comments++
			if !strings.Contains(line, "*/") {
				inComment = true
			}
		}
	}
	return total, comments, scanner.Err()
}

type jsParser struct{}

func (p *jsParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.Contains(line, "*/") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "//"):
			comments++
		case strings.HasPrefix(line, "/*"):
			comments++
			if !strings.Contains(line, "*/") {
				inMultiLine = true
			}
		}
	}
	return total, comments, scanner.Err()
}

type goParser struct{}

func (p *goParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.Contains(line, "*/") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "//"):
			comments++
		case strings.HasPrefix(line, "/*"):
			comments++
			if !strings.Contains(line, "*/") {
				inMultiLine = true
			}
		}
	}
	return total, comments, scanner.Err()
}
