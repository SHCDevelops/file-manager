package filesystem

import (
	"crypto/md5"
	"fmt"
	"github.com/SHCDevelops/file-manager/lib/utils"
	"io"
	"os"
	"path/filepath"
)

func FindDuplicates(dir string, ignoreList []string) ([][]string, error) {
	hashes := make(map[string][]string)
	var duplicates [][]string

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

		if !info.IsDir() {
			hash, err := HashFile(path)
			if err != nil {
				return err
			}

			hashes[hash] = append(hashes[hash], path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, files := range hashes {
		if len(files) > 1 {
			duplicates = append(duplicates, files)
		}
	}

	return duplicates, nil
}

func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
