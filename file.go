package main

import (
	"log/slog"
	"os"

	mymazda "github.com/taylormonacelli/forestfish/mymazda"
	"github.com/taylormonacelli/oliveluck"
)

func getBaseDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "."
}

func pathsBySubDir() (map[string]string, error) {
	baseDir := getBaseDir()
	r, err := mymazda.ExpandTilde(baseDir)
	if err != nil {
		slog.Error("expandTilde", "error", err)
		return nil, err
	}
	baseDir = r

	var mergeMap map[string]string

	for len(mergeMap) <= randNameCount {
		pathsMap := genUniquePaths(baseDir, 2, oliveluck.GetRandNamer())
		mergeMap = mergeMaps(mergeMap, pathsMap)
	}

	return mergeMap, nil
}
