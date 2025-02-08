package filesystem

import (
	"github.com/SHCDevelops/file-manager/lib/utils"
	"os"
	"path/filepath"
	"sync"
)

func SearchFiles(dir string, pattern string, ignoreList []string) ([]string, error) {
	var matchedFiles []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errors []error

	worker := func(path string, info os.FileInfo) {
		defer wg.Done()

		if utils.IsIgnored(path, ignoreList, false) {
			return
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}

		if matched {
			mu.Lock()
			matchedFiles = append(matchedFiles, path)
			mu.Unlock()
		}
	}

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

		wg.Add(1)
		go worker(path, info)
		return nil
	})

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if err != nil {
		return nil, err
	}
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return matchedFiles, nil
}
