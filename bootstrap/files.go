package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	archive = "./lfp.zip"
	project = "./project"
)

// GetSignedURL - Get signed URL for integration
// ZIP based on instance token, integration ID
// and integration version.
func GetSignedURL() (string, error) {
	url := fmt.Sprintf("https://%s.execute-api.eu-west-1.amazonaws.com/prod/registry/lfp", EnvRestAPIID)
	body := map[string]string{
		"token":   EnvToken,
		"id":      EnvID,
		"version": EnvVersion,
	}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	signedURL, err := ioutil.ReadAll(resp.Body)

	return string(signedURL), err
}

// DownloadFile - Download a file from 'url' to 'dest'.
func DownloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create destination file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to destination file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return err
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// StartIntegration ...
func StartIntegration() {
	go Status(MessageTypeInfo, "starting integration")

	signedURL, err := GetSignedURL()
	if err != nil {
		Status(MessageTypeError, "unable to load integration files")
		log.Fatal(err)
	}

	err = DownloadFile(signedURL, archive)
	if err != nil {
		Status(MessageTypeError, "failed to start integration")
		log.Fatal(err)
	}

	_, err = Unzip(archive, project)
	if err != nil {
		Status(MessageTypeError, "failed to write all project files to disk")
		log.Fatal(err)
	}

	// Continue w/ sub-process
}
