package filesystem

import (
	"bufio"
	"bytes"
	"github.com/SHCDevelops/file-manager/lib/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	maxScanTokenSize = 1024 * 1024 // 1 MB
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

	semaphore := make(chan struct{}, 10) // Пул из 10 горутин

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
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
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

	reader := bufio.NewReaderSize(file, maxScanTokenSize)
	var parser LineParser
	switch lang {
	case "HTML":
		parser = &htmlParser{}
	case "CSS":
		parser = &cssParser{}
	case "JavaScript", "TypeScript", "Java", "C++", "C", "PHP", "Swift", "Kotlin", "Rust", "Dart", "C#", "Scala", "Go":
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
	return parser.Parse(reader)
}

type LineParser interface {
	Parse(*bufio.Reader) (total int, comments int, err error)
}

type htmlParser struct{}

func (p *htmlParser) Parse(reader *bufio.Reader) (int, int, error) {
	inComment := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inComment:
			comments++
			if strings.Contains(lineStr, "-->") {
				inComment = false
			}
		case strings.Contains(lineStr, "<!--"):
			comments++
			if !strings.Contains(lineStr, "-->") {
				inComment = true
			}
		}
	}
	return total, comments, nil
}

type cssParser struct{}

func (p *cssParser) Parse(reader *bufio.Reader) (int, int, error) {
	inComment := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inComment:
			comments++
			if strings.Contains(lineStr, "*/") {
				inComment = false
			}
		case strings.HasPrefix(lineStr, "/*"):
			comments++
			if !strings.Contains(lineStr, "*/") {
				inComment = true
			}
		}
	}
	return total, comments, nil
}

type cStyleParser struct{}

func (p *cStyleParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.Contains(lineStr, "*/") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "//"):
			comments++
		case strings.HasPrefix(lineStr, "/*"):
			comments++
			if !strings.Contains(lineStr, "*/") {
				inMultiLine = true
			}
		}
	}
	return total, comments, nil
}

type hashParser struct{}

func (p *hashParser) Parse(reader *bufio.Reader) (int, int, error) {
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		if strings.HasPrefix(lineStr, "#") {
			comments++
		}
	}
	return total, comments, nil
}

type rubyParser struct{}

func (p *rubyParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.HasPrefix(lineStr, "=end") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "=begin"):
			comments++
			inMultiLine = true
		case strings.HasPrefix(lineStr, "#"):
			comments++
		}
	}
	return total, comments, nil
}

type haskellParser struct{}

func (p *haskellParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.Contains(lineStr, "-}") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "--"):
			comments++
		case strings.HasPrefix(lineStr, "{-"):
			comments++
			if !strings.Contains(lineStr, "-}") {
				inMultiLine = true
			}
		}
	}
	return total, comments, nil
}

type sqlParser struct{}

func (p *sqlParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.Contains(lineStr, "*/") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "--"):
			comments++
		case strings.HasPrefix(lineStr, "/*"):
			comments++
			if !strings.Contains(lineStr, "*/") {
				inMultiLine = true
			}
		}
	}
	return total, comments, nil
}

type luaParser struct{}

func (p *luaParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.Contains(lineStr, "]]") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "--"):
			if strings.HasPrefix(lineStr, "--[[") {
				comments++
				if !strings.Contains(lineStr, "]]") {
					inMultiLine = true
				}
			} else {
				comments++
			}
		}
	}
	return total, comments, nil
}

type pascalParser struct{}

func (p *pascalParser) Parse(reader *bufio.Reader) (int, int, error) {
	inMultiLine := false
	total, comments := 0, 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		total++
		lineStr := strings.TrimSpace(string(line))
		if isPrefix {
			var buf bytes.Buffer
			buf.Write(line)
			for isPrefix {
				line, isPrefix, err = reader.ReadLine()
				if err != nil {
					return 0, 0, err
				}
				buf.Write(line)
			}
			lineStr = strings.TrimSpace(buf.String())
		}
		switch {
		case inMultiLine:
			comments++
			if strings.Contains(lineStr, "}") {
				inMultiLine = false
			}
		case strings.HasPrefix(lineStr, "//"):
			comments++
		case strings.HasPrefix(lineStr, "{"):
			comments++
			if !strings.Contains(lineStr, "}") {
				inMultiLine = true
			}
		}
	}
	return total, comments, nil
}
