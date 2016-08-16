package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadTemplates(templateDir string) (map[string]string, error) {
	newTemplates := map[string]string{}

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return maskAny(err)
		}

		fileName := filepath.Base(path)
		fileName = fileName[0:len(fileName)-len(filepath.Ext(fileName))] + ".go"

		newTemplates[fileName] = string(raw)

		return nil
	})

	if err != nil {
		return nil, maskAny(err)
	}

	return newTemplates, nil
}
