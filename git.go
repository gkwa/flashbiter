package flashbiter

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitInit(path string) error {
	_, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	return nil
}

func gitCommitReadme(dir string) error {
	// Define the template string
	const tmpl = `* {{.Title}}

{{.Description}}

** Installation

{{.Installation}}

** Usage

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

	filename := "README.org"
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
	file, err = os.Create("README.org")
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
	_, err = worktree.Add("README.org")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Commit the changes
	_, err = worktree.Commit("Initial boilerplate", &git.CommitOptions{
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
