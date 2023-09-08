package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	baseDir := getBaseDir()
	expandTilde(&baseDir)

	count := 35

	var pn PathNamer

	namers := []PathNamer{
		// &SententiaPathNamer{},
		// &RandomdataPathNamer{},
		// &GofakeitPathNamer{},
		// &GofakeitPathNamer{},
		&Combo1{},
		&Combo2{},
		&Combo3{},
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	randomIndex := rng.Intn(len(namers))
	pn = namers[randomIndex]
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
