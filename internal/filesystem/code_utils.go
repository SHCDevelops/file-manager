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

func CountCodeLines(root string, ignoreList []string, ignoreLanguages []string) (*CodeStats, error) {
	stats := &CodeStats{
		Languages: make(map[string]*LanguageStat),
	}

	ignoredLangs := make(map[string]bool)
	for _, lang := range ignoreLanguages {
		ignoredLangs[strings.ToLower(lang)] = true
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

		lang := getLanguage(path, ignoredLangs)
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

func getLanguage(path string, ignoreLanguages map[string]bool) string {
	ext := strings.ToLower(filepath.Ext(path))
	extToLang := map[string]string{
		".html":  "HTML",
		".htm":   "HTML",
		".css":   "CSS",
		".js":    "JavaScript",
		".mjs":   "JavaScript",
		".cjs":   "JavaScript",
		".ts":    "TypeScript",
		".tsx":   "TypeScript",
		".jsx":   "JavaScript",
		".go":    "Go",
		".py":    "Python",
		".pyw":   "Python",
		".rb":    "Ruby",
		".java":  "Java",
		".cpp":   "C++",
		".cc":    "C++",
		".cxx":   "C++",
		".hpp":   "C++",
		".h":     "C",
		".c":     "C",
		".php":   "PHP",
		".swift": "Swift",
		".kt":    "Kotlin",
		".kts":   "Kotlin",
		".rs":    "Rust",
		".dart":  "Dart",
		".sh":    "Shell",
		".bash":  "Shell",
		".zsh":   "Shell",
		".pl":    "Perl",
		".pm":    "Perl",
		".lua":   "Lua",
		".sql":   "SQL",
		".cs":    "C#",
		".vb":    "Visual Basic",
		".fs":    "F#",
		".scala": "Scala",
		".hs":    "Haskell",
		".lhs":   "Haskell",
		".ml":    "OCaml",
		".mli":   "OCaml",
		".pas":   "Pascal",
		".pp":    "Pascal",
		".json":  "JSON",
		".xml":   "XML",
		".yaml":  "YAML",
		".yml":   "YAML",
		".toml":  "TOML",
	}

	lang := extToLang[ext]
	if lang == "" {
		return ""
	}

	if ignoreLanguages[strings.ToLower(lang)] {
		return ""
	}

	return extToLang[ext]
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
	case "JavaScript", "TypeScript", "Java", "C++", "C", "PHP", "Swift", "Kotlin", "Rust", "Dart", "C#", "Scala":
		parser = &cStyleParser{}
	case "Go":
		parser = &cStyleParser{}
	case "Python", "Shell", "Perl", "YAML":
		parser = &hashParser{}
	case "Ruby":
		parser = &rubyParser{}
	case "Haskell":
		parser = &haskellParser{}
	case "SQL":
		parser = &sqlParser{}
	case "Lua":
		parser = &luaParser{}
	case "Pascal":
		parser = &pascalParser{}
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

type cStyleParser struct{}

func (p *cStyleParser) Parse(file *os.File) (int, int, error) {
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

type hashParser struct{}

func (p *hashParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			comments++
		}
	}
	return total, comments, scanner.Err()
}

type rubyParser struct{}

func (p *rubyParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.HasPrefix(line, "=end") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "=begin"):
			comments++
			inMultiLine = true
		case strings.HasPrefix(line, "#"):
			comments++
		}
	}
	return total, comments, scanner.Err()
}

type haskellParser struct{}

func (p *haskellParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.Contains(line, "-}") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "--"):
			comments++
		case strings.HasPrefix(line, "{-"):
			comments++
			if !strings.Contains(line, "-}") {
				inMultiLine = true
			}
		}
	}
	return total, comments, scanner.Err()
}

type sqlParser struct{}

func (p *sqlParser) Parse(file *os.File) (int, int, error) {
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
		case strings.HasPrefix(line, "--"):
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

type luaParser struct{}

func (p *luaParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.Contains(line, "]]") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "--"):
			if strings.HasPrefix(line, "--[[") {
				comments++
				if !strings.Contains(line, "]]") {
					inMultiLine = true
				}
			} else {
				comments++
			}
		}
	}
	return total, comments, scanner.Err()
}

type pascalParser struct{}

func (p *pascalParser) Parse(file *os.File) (int, int, error) {
	scanner := bufio.NewScanner(file)
	inMultiLine := false
	total, comments := 0, 0

	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())

		switch {
		case inMultiLine:
			comments++
			if strings.Contains(line, "}") {
				inMultiLine = false
			}
		case strings.HasPrefix(line, "//"):
			comments++
		case strings.HasPrefix(line, "{"):
			comments++
			if !strings.Contains(line, "}") {
				inMultiLine = true
			}
		}
	}
	return total, comments, scanner.Err()
}
