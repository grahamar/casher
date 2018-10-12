package utils

import (
	"os"
	"sort"

	"github.com/apex/log"
	"github.com/google/go-cmp/cmp"
)

var trans = cmp.Transformer("Sort", func(in []string) []string {
	out := append(in[:0:0], in...) // Copy input to avoid mutating it
	sort.Strings(out)
	return out
})

// IsArchiveFetched returns true if cache archive exists
func IsArchiveFetched(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// IsChanged returns true if cached directories have changed
func IsChanged(cacheDir string, paths []string) (bool, error) {
	err := WriteCacheChecksums("after", cacheDir, paths)
	if err != nil {
		log.Errorf("Unable to write checks file %v", err)
		return false, err
	}

	before, err := ReadCacheChecksums("before", cacheDir)
	if err != nil {
		log.Errorf("Unable to read before checks file %v", err)
		return false, err
	}

	after, err := ReadCacheChecksums("after", cacheDir)
	if err != nil {
		log.Errorf("Unable to read after checks file %v", err)
		return false, err
	}

	if len(before) != len(after) {
		log.Info("change detected (content changed, file is created, or file is deleted)")
		return true, nil
	}

	b := struct{ Strings []string }{before}
	a := struct{ Strings []string }{after}
	return !cmp.Equal(b, a, trans), nil
}
