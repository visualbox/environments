package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/tls"
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
	app     = "./app"
)

var (
	// Proc ...
	Proc *exec.Cmd
	tr   = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
	envCmd     = map[string]map[string]string{
		// Backward compat
		"nodejs": map[string]string{
			"prepare": "yarn install",
			"run":     "node index.js",
		},
		"node": map[string]string{
			"prepare": "yarn install",
			"run":     "node index.js",
		},
		"python3": map[string]string{
			"prepare": "pip3 install -r requirements.txt",
			"run":     "python3 main.py",
		},
		"golang": map[string]string{
			"prepare": "glide install",
			"run":     "go run *.go",
		},
	}
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
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	signedURL, err := ioutil.ReadAll(resp.Body)

	return string(signedURL), err
}

func downloadFile(url string, dest string) error {
	resp, err := httpClient.Get(url)
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

func pipeToStream(
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
	Proc = exec.Command("/bin/sh", "-c", command)
	Proc.Dir = fmt.Sprintf("%v/app", os.Getenv("HOME"))
	Proc.Env = os.Environ()
	Proc.Env = append(Proc.Env, fmt.Sprintf("MODEL=%s", EnvModel))

	stdoutChan := make(chan struct{})
	stderrChan := make(chan struct{})

	var wgStream sync.WaitGroup
	wgStream.Add(2)

	go pipeToStream(true, onStdout, &wgStream, Proc, stdoutChan)
	go pipeToStream(false, onStderr, &wgStream, Proc, stderrChan)

	stdoutChan <- struct{}{}
	stderrChan <- struct{}{}
	Proc.Start()

	wgStream.Wait()

	return nil
}

func prepare() error {
	cmd := envCmd[EnvRuntime]["prepare"]
	go Status(StatusTypeInfo, cmd)

	onStdout := func(m string) {
		go Status(StatusTypeInfo, m)
	}
	onStderr := func(m string) {
		go Status(StatusTypeWarning, m)
	}

	return cmdStream(cmd, onStdout, onStderr)
}

func run() error {
	cmd := envCmd[EnvRuntime]["run"]
	go Status(StatusTypeInfo, cmd)

	// Parse output string as JSON.
	// If it fails, send a status instead.
	onStdout := func(m string) {
		go Status(StatusTypeInfo, m)
	}
	onStderr := func(m string) {
		go Status(StatusTypeError, m)
	}

	return cmdStream(cmd, onStdout, onStderr)
}

// StartIntegration ...
func StartIntegration() {
	go Status(StatusTypeInfo, "starting integration")

	signedURL, err := getSignedURL()
	if err != nil {
		Status(StatusTypeError, "unable to load integration files")
		log.Fatal(err)
		Terminate(true)
	}

	err = downloadFile(signedURL, archive)
	if err != nil {
		Status(StatusTypeError, "failed to download integration files")
		log.Fatal(err)
		Terminate(true)
	}

	_, err = unzip(archive, app)
	if err != nil {
		Status(StatusTypeError, "failed to write integration files to disk")
		log.Fatal(err)
		Terminate(true)
	}

	// Remove ZIP archive
	err = os.Remove(archive)
	if err != nil {
		log.Println(err)
	}

	// Kill any previous process
	Terminate(false)

	// Run prepare CMD
	err = prepare()
	if err != nil {
		Status(StatusTypeError, "unable to run prepare command")
		log.Fatal(err)
		Terminate(true)
	}

	// Run integration CMD
	err = run()
	if err != nil {
		Status(StatusTypeError, "unable to run integration command")
		log.Fatal(err)
		Terminate(true)
	}
}
