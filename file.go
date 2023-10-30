package flashbiter

import (
	"log/slog"
	"os"
	"path/filepath"

	mymazda "github.com/taylormonacelli/forestfish/mymazda"
	"github.com/taylormonacelli/oliveluck"
)

func getBaseDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	return "."
}

func genPathsBySubDir() (map[string]string, error) {
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

func genUniquePaths(baseDir string, numPaths int, pn oliveluck.Namer) map[string]string {
	myMap := make(map[string]string)

	for count := 0; count < numPaths; {
		subdir := pn.GetName()
		fullpath := filepath.Join(baseDir, subdir)
		_, found := myMap[subdir]

		if found {
			continue
		}

		if mymazda.DirExists(fullpath) || mymazda.FileExists(fullpath) {
			continue
		}

		myMap[subdir] = fullpath
		count++
	}

	return myMap
}
