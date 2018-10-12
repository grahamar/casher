package fetch

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/grahamar/casher/root"
	"github.com/grahamar/casher/utils"
)

// aws region.
var region string

// s3 cache bucket.
var bucket string

// cache directory.
var cacheDir string

// cache archive key.
var key string

// Command config.
var Command = &cobra.Command{
	Use:   "push [options]",
	Short: "Push cache to S3",
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

	defaultAwsRegion := os.Getenv("AWS_REGION")
	if defaultAwsRegion == "" {
		defaultAwsRegion = os.Getenv("AWS_DEFAULT_REGION")
	}

	f := Command.Flags()
	f.StringVarP(&region, "region", "r", defaultAwsRegion, "AWS region")
	f.StringVarP(&bucket, "bucket", "b", "", "Set S3 bucket")
	f.StringVarP(&cacheDir, "directory", "d", filepath.Join(home, ".casher"), "Set cache directory")
	f.StringVarP(&key, "key", "k", "", "S3 key to upload cache to")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if bucket == "" || key == "" {
		return errors.New("You must supply an S3 bucket and key")
	}

	cachedDirectories, err := utils.ReadPathsFile(cacheDir)
	if err != nil {
		log.Errorf("Unable to read cache paths file %v", err)
		return err
	}

	changed := true
	fetchFilename := filepath.Join(cacheDir, "fetch.tar.gz")
	if utils.IsArchiveFetched(fetchFilename) {
		changed, err = utils.IsChanged(cacheDir, cachedDirectories)
		if err != nil {
			log.Errorf("Unable to check for changes %v", err)
			return err
		}
	}

	if !changed {
		log.Info("nothing changed, not updating cache")
		return nil
	}

	log.Info("changes detected, packing new archive")

	pushTarGzFilename := filepath.Join(cacheDir, "push.tar.gz")
	file, err := os.Create(pushTarGzFilename)
	if err != nil {
		log.Errorf("Unable to open cache file %v", err)
		return err
	}

	err = utils.BuildTarGz(file, cachedDirectories)
	if err != nil {
		log.Errorf("Unable to build cache archive %v", err)
		return err
	}

	file.Close()
	file, err = os.OpenFile(pushTarGzFilename, os.O_RDONLY, 0755)
	if err != nil {
		log.Errorf("Unable to read cache archive %v", err)
		return err
	}
	defer file.Close()

	log.Info("uploading archive")

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		ContentType: aws.String("application/gzip"),
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
	})

	return err
}
