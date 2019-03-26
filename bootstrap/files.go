package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	archive = "./lfp.zip"
	project = "./project"
)

var (
	// Proc ...
	Proc *exec.Cmd
)

// GetSignedURL - Get signed URL for integration
// ZIP based on instance token, integration ID
// and integration version.
func getSignedURL() (string, error) {
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

func downloadFile(url string, dest string) error {
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
func unzip(src string, dest string) ([]string, error) {

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

func pipeToSocket(
	useStdOut bool,
	onStdio func(string),
	wgStream *sync.WaitGroup,
	cmd *exec.Cmd,
	c chan struct{},
) {
	defer wgStream.Done()

	var stdio io.ReadCloser
	var err error

	if useStdOut {
		stdio, err = cmd.StdoutPipe()
	} else {
		stdio, err = cmd.StderrPipe()
	}
	if err != nil {
		panic(err)
	}
	<-c
	scanner := bufio.NewScanner(stdio)
	for scanner.Scan() {
		onStdio(scanner.Text())
	}
}

func cmdStream(command string, onStdout func(string), onStderr func(string)) error {
	Proc := exec.Command("bash", "-c", command)
	stdoutChan := make(chan struct{})
	stderrChan := make(chan struct{})

	var wgStream sync.WaitGroup
	wgStream.Add(2)

	go pipeToSocket(true, onStdout, &wgStream, Proc, stdoutChan)
	go pipeToSocket(false, onStderr, &wgStream, Proc, stderrChan)

	stdoutChan <- struct{}{}
	stderrChan <- struct{}{}
	Proc.Start()

	wgStream.Wait()

	if err := Proc.Process.Kill(); err != nil {
		return err
	}
	return nil
}

func prepare() error {
	go Status(StatusTypeInfo, "yarn install")

	onStdout := func(m string) {
		go Status(StatusTypeInfo, m)
	}
	onStderr := func(m string) {
		go Status(StatusTypeWarning, m)
	}

	return cmdStream("cd project && yarn install", onStdout, onStderr)
}

func run() error {
	go Status(StatusTypeInfo, "node index.js")

	onStdout := func(m string) {
		go Status(StatusTypeInfo, m)
	}
	onStderr := func(m string) {
		go Status(StatusTypeError, m)
	}

	return cmdStream("cd project && node index", onStdout, onStderr)
}

// StartIntegration ...
func StartIntegration() {
	go Status(StatusTypeInfo, "starting integration")

	signedURL, err := getSignedURL()
	if err != nil {
		go Status(StatusTypeError, "unable to load integration files")
		log.Fatal(err)
	}

	err = downloadFile(signedURL, archive)
	if err != nil {
		go Status(StatusTypeError, "failed to start integration")
		log.Fatal(err)
	}

	_, err = unzip(archive, project)
	if err != nil {
		go Status(StatusTypeError, "failed to write all project files to disk")
		log.Fatal(err)
	}

	// Remove ZIP archive
	err = os.Remove(archive)
	if err != nil {
		log.Println(err)
	}

	// Run prepare CMD
	err = prepare()
	if err != nil {
		go Status(StatusTypeError, "unable to run prepare command")
		log.Fatal(err)
	}

	// Run integration CMD
	err = run()
	if err != nil {
		go Status(StatusTypeError, "unable to run integration command")
		log.Fatal(err)
	}
}
