package flashbiter

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"

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
	pathMap, err := genPathsBySubDir()
	if err != nil {
		slog.Error("pathsBySubDir", "error", err)
		return "", err
	}

	paths := make([]string, 0, len(pathMap))
	for path := range pathMap {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	inputSelector := aeryavenue.GetInputSelector()
	selectedPath, err := inputSelector.SelectItem(paths)
	if err != nil {
		return "", err
	}

	return selectedPath, nil
}

func mergeMaps(map1, map2 map[string]string) map[string]string {
	merged := make(map[string]string)

	for key, value := range map1 {
		merged[key] = value
	}

	for key, value := range map2 {
		merged[key] = value
	}

	return merged
}
