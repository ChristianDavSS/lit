package analysis

import (
	"CLI_App/src/internals"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Contributor struct {
	name, email string
	commits     []*Commit
}

type Commit struct {
	hash          plumbing.Hash
	when, message string
	stats         object.FileStats
}

var repo = internals.GetGitRepository()

func FetchCommits() {
	r, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	getContributors(r)
}

// Function to get all the contributors and their data from a git branch (main default)
func getContributors(iter object.CommitIter) {
	// authors slice to save up the contributors of the project
	var authors = make(map[string]*Contributor)
	// iterate through the CommitIter
	err := iter.ForEach(func(commit *object.Commit) error {
		// obj create a Contributor object
		email := commit.Author.Email
		// Get the user from the map
		c, ok := authors[email]
		// Executes if the user isn't on the map keys
		if !ok {
			commits := make([]*Commit, 0)
			commits = append(commits, getCommitInformation(commit))
			authors[email] = &Contributor{
				name:    commit.Author.Name,
				email:   commit.Author.Email,
				commits: commits,
			}
			// Executes if the user email IS on the map
		} else {
			c.commits = append(c.commits, getCommitInformation(commit))
		}
		// Return nil as err
		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print out the contributors
	for _, v := range authors {
		fmt.Println(v)
	}
}

// Function that receives an object.Commit and returns our Commit object filled up with the data we needed.
func getCommitInformation(commit *object.Commit) *Commit {
	year, month, day := commit.Author.When.Date()
	hour, minute, sec := commit.Author.When.Clock()
	return &Commit{
		hash: commit.Hash,
		when: fmt.Sprintf("%s %s %d %d %d:%d:%d",
			commit.Author.When.Weekday(), month, day, year, hour, minute, sec),
	}
}
