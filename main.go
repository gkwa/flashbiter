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
		&SententiaPathNamer{},
		&RandomdataPathNamer{},
		&GofakeitPathNamer{},
		&Combo1{},
		&Combo2{},
		&Combo4{},
		&Combo5{},
		&Combo6{},
		&Combo7{},
		&Combo8{},
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	var uniquePaths map[string]string

	for {
		randomIndex := rng.Intn(len(namers))
		pn = namers[randomIndex]
		x := generateUniquePaths(baseDir, 2, pn)
		uniquePaths = mergeMaps(uniquePaths, x)
		if len(uniquePaths) >= count {
			break
		}
	}

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

	absPath, err := getAbsPath(selectedPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(absPath)
	clipboard.WriteAll(absPath)
}
