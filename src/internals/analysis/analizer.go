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

// FetchCommits main function to keep all the flags on.
func FetchCommits(who string, verbose bool) {
	r, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get contributors of the Git branch
	authors := getContributors(r)
	// If the 'who' flag was sent, we do the logics for that
	if len(who) > 0 {
		// We search by the git username
		contrib, ok := authors[who]
		if ok {
			printContributorData(contrib, verbose)
			return
		}
		// We search by the git email
		contrib, err := findContribBy(authors, who)
		if err != nil {
			panic(err)
		}
		printContributorData(contrib, verbose)
		return
	}

	// If there´s not a 'who' flag, we print all the data
	for _, v := range authors {
		printContributorData(v, verbose)
	}
}

// ----------------------------------------------------------------------------------------
func findContribBy(contributors map[string]*Contributor, target string) (*Contributor, error) {
	for _, v := range contributors {
		if v.email == target {
			return v, nil
		}
	}
	return nil, fmt.Errorf("couldn't find a contributor with that email")
}

func printContributorData(contributor *Contributor, verbose bool) {
	fmt.Printf("- %s <%s>\n", contributor.name, contributor.email)
	fmt.Printf("  Total of commits: %d\n", len(contributor.commits))
	if verbose {
		for _, c := range contributor.commits {
			fmt.Printf("* %s %s", c.when, c.message)
			fmt.Printf("  Hash: %s\n", c.hash)
			fmt.Printf("  Modified files: %d\n", len(c.stats))
			fmt.Printf("  Stats:\n")
			formatPrintStats(c.stats)
		}
	}
}

// Function to print a formatted version of the commit stats.
func formatPrintStats(stats object.FileStats) {
	for _, v := range stats {
		fmt.Printf("  %s", v)
	}
}

// ------------------------------------------------------------------------------------------

// Function to get all the contributors and their data from a git branch (main default)
func getContributors(iter object.CommitIter) map[string]*Contributor {
	// authors slice to save up the contributors of the project
	var authors = make(map[string]*Contributor)
	// iterate through the CommitIter
	err := iter.ForEach(func(commit *object.Commit) error {
		// obj create a Contributor object
		username := commit.Author.Name
		// Get the user from the map
		c, ok := authors[username]
		// Executes if the user isn't on the map keys
		if !ok {
			commits := make([]*Commit, 0)
			commits = append(commits, getCommitInformation(commit))
			authors[username] = &Contributor{
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

	return authors
}

// Function that receives an object.Commit and returns our Commit object filled up with the data we needed.
func getCommitInformation(commit *object.Commit) *Commit {
	year, month, day := commit.Author.When.Date()
	hour, minute, sec := commit.Author.When.Clock()
	stats, err := commit.Stats()
	if err != nil {
		fmt.Println("Couldn´t get the stats of the commit.")
		os.Exit(1)
	}
	return &Commit{
		hash: commit.Hash,
		when: fmt.Sprintf("%s %s %d %d %d:%d:%d",
			commit.Author.When.Weekday(), month, day, year, hour, minute, sec),
		message: commit.Message,
		stats:   stats,
	}
}
