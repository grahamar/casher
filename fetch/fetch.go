package fetch

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/grahamar/casher/root"
)

// aws region.
var region string

// s3 cache bucket.
var bucket string

// cache directory.
var cacheDir string

// check keys.
var keys []string

// Command config.
var Command = &cobra.Command{
	Use:   "fetch [options]",
	Short: "Fetch cache from S3",
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
	f.StringArrayVarP(&keys, "keys", "k", []string{}, "S3 keys to check")
}

func downloadArchive(key string, downloader *s3manager.Downloader) error {
	file, err := os.Create(filepath.Join(cacheDir, "fetch.tar.gz"))
	if err != nil {
		log.Errorf("Unable to open cache file %v", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	return err
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if bucket == "" {
		return errors.New("You must supply an S3 bucket")
	}

	os.MkdirAll(cacheDir, 0755)

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	downloader := s3manager.NewDownloader(sess)
	svc := s3.New(sess)

	log.Info("attempting to download cache archive")

	var archiveFound bool
	for _, key := range keys {
		log.Infof("fetching %s", key)
		_, err := svc.HeadObject(&s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
		if err == nil {
			log.Info("found cache")
			archiveFound = true
			if downloadArchive(key, downloader) != nil {
				log.Errorf("Unable to download cache file %v", err)
				return err
			}
			break
		}
	}

	if !archiveFound {
		log.Info("could not download cache")
	}

	return nil
}
