package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/atotto/clipboard"
	"github.com/castillobgr/sententia"
	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rivo/tview"
	log "github.com/taylormonacelli/reactnut/cmd/logging"
)

type Item struct {
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

type BlackholeDestination struct{}

// fixme: never see otuput of this, dunno why
func (cd *BlackholeDestination) Write(data string) error {
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

func getBaseDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "."
}

type PathNamer interface {
	GetName() string
}

type RandomPathNamer struct{}

func (rpn *RandomPathNamer) GetName() string {
	adjective := randomdata.Adjective()
	noun := randomdata.Noun()
	return strings.ToLower(adjective + noun)
}

type SententiaPathNamer struct{}

func (spn *SententiaPathNamer) GetName() string {
	str, err := sententia.Make("{{ adjective }}{{ nouns }}")
	if err != nil {
		panic(err)
	}
	return strings.ToLower(str)
}

func generateUniquePaths(baseDir string, numPaths int, pn PathNamer) map[string]string {
	myMap := make(map[string]string)
	for i := 0; i < numPaths; {
		subdir := pn.GetName()
		fullpath := filepath.Join(baseDir, subdir)
		if _, keyExists := myMap[subdir]; keyExists || pathExists(fullpath) {
			continue
		}
		myMap[subdir] = fullpath
		i++
	}
	return myMap
}

type InputSelector interface {
	SelectItem(keys []string) (string, error)
}

type TviewInputSelector struct{}

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

type RandomItemSelector struct{}

func (ris *RandomItemSelector) SelectItem(keys []string) (string, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(keys))
	return keys[index], nil
}

func selectPath(paths map[string]string, is InputSelector) (string, error) {
	sortedKeys := sortedKeys(paths)
	item, err := is.SelectItem(sortedKeys)
	if err != nil {
		return "", err
	}
	return paths[item], nil
}

func GitInit(path string) error {
	_, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	return nil
}

func gitCommitReadme(dir string) error {
	// Define the template string
	const tmpl = `# {{.Title}}

{{.Description}}

## Installation

{{.Installation}}

## Usage

{{.Usage}}
`

	// Define the data for the template
	data := struct {
		Title        string
		Description  string
		Installation string
		Usage        string
	}{
		Title:        "My Project",
		Description:  "This is a description of my project.",
		Installation: "To install my project, run this command: go install myproject",
		Usage:        "To use my project, run this command: myproject",
	}

	// Create a new directory and its parent directories
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	filename := "README.md"
	fullPath := filepath.Join(dir, filename)

	// Create the new file
	file, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Parse the template
	tpl, err := template.New("README").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	// Execute the template with the data and write the output to the file
	err = tpl.Execute(file, data)
	if err != nil {
		panic(err)
	}

	// Open the repository
	repo, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the new file
	file, err = os.Create("README.md")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// Add the new file to the repository
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = worktree.Add("README.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Commit the changes
	_, err = worktree.Commit("Add README.md", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your.email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func stringToBool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

func getInputSelector() InputSelector {
	ris := &RandomItemSelector{}
	uis := &TviewInputSelector{}
	envVars := []string{"GITHUB_ACTIONS", "GITLAB_CI"}

	for _, envVar := range envVars {
		s := os.Getenv(envVar)
		b, err := stringToBool(s)
		if err != nil {
			return uis
		}
		if b {
			return ris
		}
	}
	return uis
}

func main() {
	baseDir := getBaseDir()
	candidateCount := 35

	var pn PathNamer
	if rand.Intn(2) == 0 {
		pn = &SententiaPathNamer{}
	} else {
		pn = &RandomPathNamer{}
	}
	uniquePaths := generateUniquePaths(baseDir, candidateCount, pn)

	inputSelector := getInputSelector()
	selectedPath, err := selectPath(uniquePaths, inputSelector)
	if err != nil {
		panic(err)
	}

	GitInit(selectedPath)
	_ = gitCommitReadme(selectedPath)

	absPath, err := filepath.Abs(selectedPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Printf("%s is not a directory\n", absPath)
		os.Exit(1)
	}

	fmt.Println(absPath)
}
