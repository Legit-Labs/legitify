package cmd

import (
	"log"
	"os"
)

func setErrorFile(path string) (*os.File, error) {
	file, err := openForWrite(path)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, err
}

func openForWrite(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func setOutputFile(path string) error {
	if path == "" { // default to stdout
		return nil
	}

	file, err := openForWrite(path)
	if err != nil {
		return err
	}

	if err := os.Stdout.Close(); err != nil {
		return err
	}

	os.Stdout = file
	return nil
}
