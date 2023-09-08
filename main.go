package main

import (
	"fmt"
	"math/rand"

	"github.com/atotto/clipboard"
)

func main() {
	baseDir := getBaseDir()
	expandTilde(&baseDir)

	count := 35

	var pn PathNamer

	randomNumber := rand.Intn(3)
	switch randomNumber {
	case 0:
		pn = &SententiaPathNamer{}
	case 1:
		pn = &RandomdataPathNamer{}
	case 2:
		pn = &GofakeitPathNamer{}
	}

	uniquePaths := generateUniquePaths(baseDir, count, pn)

	inputSelector := getInputSelector()
	selectedPath, err := selectPath(uniquePaths, inputSelector)
	if err != nil {
		panic(err)
	}

	// human canceled tview
	if selectedPath == "" {
		return
	}

	GitInit(selectedPath)
	_ = gitCommitReadme(selectedPath)

	absPath, err := getAbsPath(selectedPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(absPath)
	clipboard.WriteAll(absPath)
}
