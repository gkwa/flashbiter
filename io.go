package flashbiter

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/taylormonacelli/aeryavenue"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
	"github.com/taylormonacelli/oliveluck"
)

type (
	ClipboardDestination struct{}
	ConsoleDestination   struct{}
	BlackholeDestination struct{}
)

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
		slog.Error("error writing to clipboard", "error", err)
	}

	return err
}

func (fd *FileDestination) Write(data string) error {
	err := os.WriteFile(fd.FilePath, []byte(data), 0o644)
	if err != nil {
		return err
	}

	return nil
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

func selectPath(paths map[string]string, selector aeryavenue.InputSelector) (string, error) {
	sortedKeys := sortedKeys(paths)
	item, err := selector.SelectItem(sortedKeys)
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
