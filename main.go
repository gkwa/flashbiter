package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/taylormonacelli/reactnut/cmd/logging"
)

type Fruit struct {
	Name string
}

type OutputDestination interface {
	Write(data string) error
}

type FileDestination struct {
	FilePath string
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

type ConsoleDestination struct{}

// fixme: never see otuput of this, dunno why
func (cd *ConsoleDestination) Write(data string) error {
	// time.Sleep(2000 * time.Millisecond)
	fmt.Printf("Selected value: %s\n", data)
	return nil
}

type ClipboardDestination struct{}

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

func returnValue(val string, output OutputDestination) {
	if err := output.Write(val); err != nil {
		fmt.Println("Error writing to output:", err)
	}
}

func allowUserToSelectItem(selectables []string) (string, error) {
	app := tview.NewApplication()

	var selectedFruit string

	// Create a list widget and add the fruits to it
	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			selectedFruit = selectables[index]
			// returnValue(selectedFruit, &ClipboardDestination{})
			// returnValue(selectedFruit, &ConsoleDestination{})
			// returnValue(selectedFruit, &FileDestination{FilePath: "fruits.txt"})
			app.Stop()
		})
	for _, item := range selectables {
		list.AddItem(item, "", rune(0), nil)
	}

	// Set up key bindings to navigate the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'n':
				list.SetCurrentItem((list.GetCurrentItem() + 1) % list.GetItemCount())
				return nil
			case 'p':
				current := list.GetCurrentItem()
				if current == 0 {
					current = list.GetItemCount()
				}
				list.SetCurrentItem((current - 1) % list.GetItemCount())
				return nil
			case 'q':
				app.Stop()
				return nil
			}
		}
		return event
	})

	// Set the list widget as the root and run the application
	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return selectedFruit, nil
}

func concatWords() string {
	adjective := randomdata.Adjective()
	noun := randomdata.Noun()
	concat := strings.ToLower(fmt.Sprintf("%s%s", adjective, noun))
	return concat
}

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

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func genPathStr(basePath string, subdir string, fullPath *string) error {
	*fullPath = filepath.Join(basePath, subdir)

	err := expandTilde(fullPath)
	if err != nil {
		log.Logger.Fatalf("expanding path causes error: %s", err)
	}

	if filepath.IsAbs(*fullPath) {
		return nil
	}

	c, err := os.Getwd()
	if err != nil {
		return err
	}

	*fullPath = filepath.Join(c, *fullPath)
	return nil
}

func main() {
	var baseDir string
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	} else {
		baseDir = "."
	}

	var fullPath string
	myMap := make(map[string]string)

	for i := 0; i < 35; i++ {
		subdir := concatWords()
		err := genPathStr(baseDir, subdir, &fullPath)
		if err != nil {
			panic(err)
		}

		_, keyExists := myMap[subdir]

		if keyExists || pathExists(fullPath) {
			i--
			continue
		}

		myMap[subdir] = fullPath
	}

	sorted := sortedKeys(myMap)
	item, err := allowUserToSelectItem(sorted)
	if err != nil {
		panic(err)
	}
	fmt.Println(myMap[item])
}
