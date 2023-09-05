package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TviewInputSelector struct{}

type InputSelector interface {
	SelectItem(keys []string) (string, error)
}

func (uis *TviewInputSelector) SelectItem(keys []string) (string, error) {
	app := tview.NewApplication()

	var selectedItem string

	// Create a list widget and add the items to it
	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			selectedItem = keys[index]
			returnValue(selectedItem, &BlackholeDestination{})
			// returnValue(selectedItem, &ClipboardDestination{})
			// returnValue(selectedItem, &ConsoleDestination{})
			// returnValue(selectedItem, &FileDestination{FilePath: "items.txt"})
			app.Stop()
		})
	for _, item := range keys {
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
	return selectedItem, nil
}
