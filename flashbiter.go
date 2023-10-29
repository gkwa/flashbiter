package main

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/atotto/clipboard"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
)

var randNameCount = 35

func Main() int {
	selectedPath, err := GetUniquePath()
	if err != nil {
		panic(err)
	}

	// human canceled tview
	if selectedPath == "" {
		return 0
	}

	GitInit(selectedPath)

	absPath, err := filepath.Abs(selectedPath)
	if err != nil {
		slog.Error("filepath.Abs", "error", err)
	}
	if !mymazda.DirExists(absPath) {
		panic(err)
	}

	fmt.Println(absPath)
	clipboard.WriteAll(absPath)

	return 0
}

func GetUniquePath() (string, error) {
	uniquePaths, err := pathsBySubDir()
	if err != nil {
		slog.Error("pathsBySubDir", "error", err)
		return "", err
	}

	inputSelector := getInputSelector()
	selectedPath, err := selectPath(uniquePaths, inputSelector)
	if err != nil {
		return "", err
	}

	return selectedPath, nil
}
