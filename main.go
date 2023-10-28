package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
)

func main() {
	status := Main()
	os.Exit(status)
}

var randNameCount = 35

func Main() int {
	uniquePaths, err := pathsBySubDir()
	if err != nil {
		slog.Error("doit", "error", err)
	}

	inputSelector := getInputSelector()
	selectedPath, err := selectPath(uniquePaths, inputSelector)
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
