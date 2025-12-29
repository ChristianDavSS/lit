package commits

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
)

func FetchCommits(path string) {
	head, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(head)
}
