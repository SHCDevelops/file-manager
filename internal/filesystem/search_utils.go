package filesystem

import (
	"github.com/SHCDevelops/file-manager/lib/utils"
	"os"
	"path/filepath"
)

func SearchFiles(dir string, pattern string, ignoreList []string) ([]string, error) {
	var matchedFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if utils.IsIgnored(path, ignoreList, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if utils.IsIgnored(path, ignoreList, false) {
			return nil
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}
		if matched {
			matchedFiles = append(matchedFiles, path)
		}

		return nil
	})
	
	if err != nil {
		return nil, err
	}

	return matchedFiles, nil
}
