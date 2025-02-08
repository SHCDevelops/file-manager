package filesystem

import (
	"github.com/SHCDevelops/file-manager/lib/utils"
	"os"
	"path/filepath"
	"sort"
)

type FileSize struct {
	Path string
	Size int64
}

func AnalyzeSpace(dir string, top int, ignoreList []string) ([]FileSize, error) {
	var files []FileSize

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && utils.IsIgnored(path, ignoreList) {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, FileSize{Path: path, Size: info.Size()})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > top {
		files = files[:top]
	}

	return files, nil
}
