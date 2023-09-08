package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

type (
	ClipboardDestination struct{}
	ConsoleDestination   struct{}
	BlackholeDestination struct{}
)

type RandomItemSelector struct{}

type OutputDestination interface {
	Write(data string) error
}

type FileDestination struct {
	FilePath string
}

// fixme: never see otuput of this, dunno why
func (cd *BlackholeDestination) Write(data string) error {
	return nil
}

func (cd *ClipboardDestination) Write(data string) error {
	return writeToClipboard(data)
}

func writeToClipboard(s string) error {
	err := clipboard.WriteAll(s)
	if err != nil {
		fmt.Println("Error writing to clipboard:", err)
	}
	return err
}

func (fd *FileDestination) Write(data string) error {
	file, err := os.OpenFile(fd.FilePath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(data + "\n"); err != nil {
		return err
	}
	return nil
}

func generateUniquePaths(baseDir string, numPaths int, pn PathNamer) map[string]string {
	myMap := make(map[string]string)
	for i := 0; i < numPaths; {
		subdir := pn.GetName()
		fullpath := filepath.Join(baseDir, subdir)
		if _, keyExists := myMap[subdir]; keyExists || pathExists(fullpath) {
			continue
		}
		myMap[subdir] = fullpath
		i++
	}
	return myMap
}

func selectPath(paths map[string]string, is InputSelector) (string, error) {
	sortedKeys := sortedKeys(paths)
	item, err := is.SelectItem(sortedKeys)
	if err != nil {
		return "", err
	}
	return paths[item], nil
}

func stringToBool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

// fixme: never see otuput of this, dunno why
func (cd *ConsoleDestination) Write(data string) error {
	// time.Sleep(2000 * time.Millisecond)
	fmt.Printf("Selected value: %s\n", data)
	return nil
}

func returnValue(val string, output OutputDestination) {
	if err := output.Write(val); err != nil {
		fmt.Println("Error writing to output:", err)
	}
}

func getInputSelector() InputSelector {
	ris := &RandomItemSelector{}
	uis := &TviewInputSelector{}

	// don't prompt for input while in automated pipeline
	envVars := []string{"GITHUB_ACTIONS", "GITLAB_CI"}

	for _, envVar := range envVars {
		s := os.Getenv(envVar)
		b, err := stringToBool(s)
		if err != nil {
			return uis
		}
		if b {
			return ris
		}
	}
	return uis
}

func (ris *RandomItemSelector) SelectItem(keys []string) (string, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(keys))
	return keys[index], nil
}
