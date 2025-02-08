package filesystem

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/SHCDevelops/file-manager/lib/utils"
)

type FileSize struct {
	Path string
	Size int64
}

func AnalyzeSpace(dir string, top int, ignoreList []string) ([]FileSize, error) {
	var files []FileSize
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errors []error

	worker := func(path string, info os.FileInfo) {
		defer wg.Done()

		if utils.IsIgnored(path, ignoreList, false) {
			return
		}

		mu.Lock()
		files = append(files, FileSize{Path: path, Size: info.Size()})
		mu.Unlock()
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

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > top {
		files = files[:top]
	}

	return files, nil
}
