package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/tanema/cws/lib/manifest"
)

// Zip will create a new zip archive file and update the manifest for publishing
// with the new version added and the dev key removed
func Zip(dir, version string) (string, error) {
	file, err := os.Create("compiled_extension.zip")
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	return file.Name(), filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name, err = filepath.Rel(filepath.Dir(dir), path)
		if err != nil {
			return err
		}
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		data, err := getFile(path, version)
		if err != nil {
			return err
		}
		defer data.Close()
		_, err = io.Copy(headerWriter, data)
		return err
	})
}

func getFile(path, version string) (io.ReadCloser, error) {
	if filepath.Base(path) != "manifest.json" {
		return os.Open(path)
	}
	manifestBytes, err := manifest.UpdateBytes(path, version)
	return io.NopCloser(bytes.NewBuffer(manifestBytes)), err
}
