package flashbiter

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/taylormonacelli/aeryavenue"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
)

var randNameCount = 35

func Main() int {
	selectedPath, err := GetUniquePath()
	if err != nil {
		slog.Error("GetUniquePath", "error", err)
		return 1
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
	paths, err := genPathsBySubDir()
	if err != nil {
		slog.Error("pathsBySubDir", "error", err)
		return "", err
	}

	inputSelector := aeryavenue.GetInputSelector()
	selectedPath, err := selectPath(paths, inputSelector)
	if err != nil {
		return "", err
	}

	return selectedPath, nil
}
