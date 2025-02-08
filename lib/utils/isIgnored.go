package utils

import "path/filepath"

func IsIgnored(path string, ignoreList []string, isDir bool) bool {
	for _, ignorePattern := range ignoreList {
		if isDir && filepath.Base(path) == ignorePattern {
			return true
		}

		if !isDir {
			matched, _ := filepath.Match(ignorePattern, filepath.Base(path))
			if matched {
				return true
			}
		}
	}
	return false
}
