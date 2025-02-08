package utils

import "path/filepath"

func IsIgnored(path string, ignoreList []string) bool {
	for _, pattern := range ignoreList {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}
