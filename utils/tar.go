package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
)

// Contains returns whether a slice contains the prefix for an element
func Contains(s []string, e string) bool {
	for _, a := range s {
		if strings.HasPrefix(e, a) {
			return true
		}
	}
	return false
}

// ExtractTarGz Untar multiple files from tar, tar.gz and tar.bz2 File
func ExtractTarGz(r io.Reader, files []string) error {
	uncompressedStream, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal("NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf("Next() failed: %s", err.Error())
			return err
		}

		if !Contains(files, header.Name) {
			log.Warnf("not extracting file: %s", header.Name)
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(header.Name); os.IsNotExist(err) {
				if err = os.Mkdir(header.Name, 0755); err != nil {
					log.Errorf("Mkdir() failed: %s", err.Error())
					return err
				}
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Errorf("Create() failed: %s", err.Error())
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Errorf("Copy() failed: %s", err.Error())
				return err
			}
		default:
			log.Errorf("uknown type: %s in %s", header.Typeflag, header.Name)
			return err
		}
	}

	return nil
}

func addPath(tw *tar.Writer, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}

	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name = p

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
}

// BuildTarGz build .tar.gz from list of directories
func BuildTarGz(file *os.File, paths []string) error {
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for i, path := range paths {
		if err := addPath(tw, path); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}

		if i%5 == 0 {
			fmt.Print(".")
		}
	}

	fmt.Println("")

	return nil
}
