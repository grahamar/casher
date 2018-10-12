package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/apex/log"
)

func hashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Errorf("%v", err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// WriteCacheChecksums writes md5sums file
func WriteCacheChecksums(suffix string, cacheDir string, paths []string) error {
	checkBeforeFilename := filepath.Join(cacheDir, "md5sums_"+suffix)
	checkBefore, err := os.OpenFile(checkBeforeFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer checkBefore.Close()

	hashes := []string{}
	for _, path := range paths {
		err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			hash, err := hashFile(p)
			hashes = append(hashes, fmt.Sprintf("%s  %s", hash, p))
			return err
		})
	}

	sort.Strings(hashes)

	w := bufio.NewWriter(checkBefore)
	for _, line := range hashes {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func ReadCacheChecksums(suffix string, cacheDir string) ([]string, error) {
	checkBeforeFilename := filepath.Join(cacheDir, "md5sums_"+suffix)
	file, err := os.Open(checkBeforeFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
