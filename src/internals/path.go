package internals

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// GetWorkingDirectory Method to return the current working directory
func getWorkingDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error with the path: ", err)
		os.Exit(1)
	}

	return path
}

func GetGitRepository() *git.Repository {
	repo, err := git.PlainOpen(getWorkingDirectory())
	if err != nil {
		fmt.Println("I couldn´t detect a Git repository in this path. ", err)
		os.Exit(1)
	}

	return repo
}
