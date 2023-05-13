package main

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
)

func GitInit(path string) error {
	_, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}

	fmt.Printf("Initialized empty Git repository in %s\n", path)
	return nil
}
