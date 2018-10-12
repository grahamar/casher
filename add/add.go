package fetch

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/grahamar/casher/root"
	"github.com/grahamar/casher/utils"
)

// cache directory.
var cacheDir string

// symlink paths.
var paths []string

// Command config.
var Command = &cobra.Command{
	Use:   "add [options]",
	Short: "Add local paths to cache",
	RunE:  run,
}

// Initialize.
func init() {
	root.Register(Command)

	home, err := homedir.Dir()
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	f := Command.Flags()
	f.StringVarP(&cacheDir, "directory", "d", filepath.Join(home, ".casher"), "Set cache directory")
	f.StringArrayVarP(&paths, "paths", "p", []string{}, "paths to add")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	os.MkdirAll(cacheDir, 0755)

	for _, path := range paths {
		fi, err := os.Lstat(path)
		if fi != nil && fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			log.Infof("%s is a symbolic link; not following", path)
			continue
		}

		log.Infof("adding %s to cache", path)
		if os.IsNotExist(err) {
			log.Infof("creating directory %s", path)
			os.MkdirAll(path, 0755)
		}
	}

	err := utils.WritePathsFile(cacheDir, paths)
	if err != nil {
		log.Errorf("Unable to write paths file %v", err)
		return err
	}

	fetchFilename := filepath.Join(cacheDir, "fetch.tar.gz")
	if utils.IsArchiveFetched(fetchFilename) {
		f, err := os.Open(fetchFilename)
		defer f.Close()
		err = utils.ExtractTarGz(f, paths)
		if err != nil {
			log.Errorf("Unable to extract cache archive %v", err)
			return err
		}

		err = utils.WriteCacheChecksums("before", cacheDir, paths)
		if err != nil {
			log.Errorf("Unable to write checks file %v", err)
			return err
		}
	} else {
		log.Infof("No fetched archive [%s] to extract", fetchFilename)
	}

	return nil
}
