package utils

import (
	"path/filepath"
	"strings"
)

// IsIgnored проверяет, должен ли быть проигнорирован указанный путь
// Поддерживает:
// - Паттерны в стиле .gitignore
// - Относительные пути от корня сканирования
// - Рекурсивные паттерны (**)
// - Негативные паттерны (!)
func IsIgnored(path string, ignorePatterns []string, isDir bool) bool {
	normalizedPath := normalizePath(path)
	var matched bool

	for i := len(ignorePatterns) - 1; i >= 0; i-- {
		pattern := normalizePattern(ignorePatterns[i])

		if pattern == "" {
			continue
		}

		if strings.HasPrefix(pattern, "!") {
			negatePattern := strings.TrimPrefix(pattern, "!")
			if match(negatePattern, normalizedPath, isDir) {
				return false
			}
			continue
		}

		if match(pattern, normalizedPath, isDir) {
			matched = true
		}
	}

	return matched
}

func normalizePath(path string) string {
	return filepath.ToSlash(strings.TrimPrefix(path, string(filepath.Separator)))
}

func normalizePattern(pattern string) string {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return ""
	}
	return filepath.ToSlash(pattern)
}

func match(pattern, path string, isDir bool) bool {
	if isDir && !strings.HasSuffix(pattern, "/") {
		pattern += "/"
		path += "/"
	}

	if strings.HasPrefix(pattern, "**/") {
		pattern = strings.TrimPrefix(pattern, "**/")
	} else if strings.Contains(pattern, "/**/") {
		pattern = strings.ReplaceAll(pattern, "/**/", "/")
	}

	matched, _ := filepath.Match(pattern, path)
	if !matched {
		if isDir && strings.Contains(pattern, "/") {
			matched, _ = filepath.Match(pattern+"/**", path+"/")
		}
	}

	return matched
}
