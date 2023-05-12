package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			returnValue(selectedFruit, &ClipboardDestination{})
			returnValue(selectedFruit, &ConsoleDestination{})
			returnValue(selectedFruit, &FileDestination{FilePath: "fruits.txt"})
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

func main() {
	items := []string{"apple", "pear"}

	item, err := allowUserToSelectItem(items)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(dir, item)
	fmt.Println(fullPath)
}
