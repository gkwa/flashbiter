package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	log "github.com/taylormonacelli/reactnut/cmd/logging"
)

func pathExists(path string) bool {
	err := expandTilde(&path)
	if err != nil {
		log.Logger.Fatalf("expanding tilde creates error for path: %s, error: %s",
			path, err)
	}
	log.Logger.Traceln(path) // output: /Users/username/Documents/example.txt

	// Use os.Stat() to get information about the path
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the error value is nil, which indicates that the path exists
	if err == nil {
		// Check if the path is a directory
		if fileInfo.Mode().IsDir() {
			log.Logger.Tracef("%s is a directory", path)
		} else {
			log.Logger.Tracef("%s is a file", path)
		}
	} else {
		log.Logger.Tracef("Path %s does not exist", path)
	}
	return true
}

func getBaseDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "."
}

func expandTilde(path *string) error {
	if strings.HasPrefix(*path, "~/") || *path == "~" {
		currentUser, err := user.Current()
		if err != nil {
			log.Logger.Warningf("checking current user results in error: %s", err)
			return err
		}
		*path = strings.Replace(*path, "~", currentUser.HomeDir, 1)
		log.Logger.Tracef("path is expanded to %s", *path)
		return nil
	}
	return nil
}

func getAbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !fileInfo.IsDir() {
		msg := fmt.Sprintf("%s is not a directory\n", absPath)
		panic(msg)
	}

	return absPath, nil
}
