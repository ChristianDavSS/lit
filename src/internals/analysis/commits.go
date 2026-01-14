package analysis

import (
	"CLI_App/src/internals"
	"CLI_App/src/internals/analysis/utils"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var repo = internals.GetGitRepository()

// FetchCommits main function to keep all the flags on.
func FetchCommits(who string, verbose, stats bool, since, until string, commitSize bool) {
	r, err := repo.Log(&git.LogOptions{
		All:   true,
		Since: utils.ValidateDate(since),
		Until: utils.ValidateDate(until),
	})
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
			printContributorData(contrib, verbose, stats, commitSize)
			return
		}
		// We search by the git email
		contrib, err := findContribBy(authors, who)
		if err != nil {
			panic(err)
		}
		printContributorData(contrib, verbose, stats, commitSize)
		return
	}

	// If there´s not a 'who' flag, we print all the data
	for _, v := range authors {
		printContributorData(v, verbose, stats, commitSize)
	}
}

// ----------------------------------------------------------------------------------------
func findContribBy(contributors map[string]*utils.Contributor, target string) (*utils.Contributor, error) {
	for _, v := range contributors {
		if v.Email == target {
			return v, nil
		}
	}
	return nil, fmt.Errorf("couldn't find a contributor with that email")
}

func printContributorData(contributor *utils.Contributor, verbose, stats, commitSize bool) {
	fmt.Printf("- %s <%s>\n", contributor.Name, contributor.Email)
	fmt.Printf("  Total of commits: %d\n", len(contributor.Commits))
	for _, c := range contributor.Commits {
		// Calculate the number of files changed
		totalChanges := len(c.Stats)
		if verbose {
			fmt.Printf("  * %s %s", c.When, c.Message)
			fmt.Printf("  Hash: %s\n", c.Hash)
			fmt.Printf("  Modified files: %d\n", totalChanges)
		}
		if commitSize {
			addition, deletion := 0, 0
			for _, v := range c.Stats {
				addition += v.Addition
				deletion += v.Deletion
			}
			fmt.Printf("  Total lines of code changed (mean): %d\n", (addition+deletion)/totalChanges)
			fmt.Printf("  Additions (mean): %d\nDeletions (mean): %d\n",
				addition/totalChanges, deletion/totalChanges)
		}
		if stats {
			fmt.Printf("  Stats:\n")
			formatPrintStats(c.Stats)
		}
	}
	fmt.Println()
}

// Function to print a formatted version of the commit stats.
func formatPrintStats(stats object.FileStats) {
	for _, v := range stats {
		fmt.Printf("  %s", v)
	}
	fmt.Println()
}

// ------------------------------------------------------------------------------------------

// Function to get all the contributors and their data from a git branch (main default)
func getContributors(iter object.CommitIter) map[string]*utils.Contributor {
	// authors slice to save up the contributors of the project
	var authors = make(map[string]*utils.Contributor)
	// iterate through the CommitIter
	err := iter.ForEach(func(commit *object.Commit) error {
		// obj create a Contributor object
		email := commit.Author.Email
		// Get the user from the map
		c, ok := authors[email]
		// Executes if the user isn't on the map keys
		if !ok {
			commits := make([]*utils.Commit, 0)
			commits = append(commits, getCommitInformation(commit))
			authors[email] = &utils.Contributor{
				Name:    commit.Author.Name,
				Email:   commit.Author.Email,
				Commits: commits,
			}
			// Executes if the user email IS on the map
		} else {
			c.Commits = append(c.Commits, getCommitInformation(commit))
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
func getCommitInformation(commit *object.Commit) *utils.Commit {
	year, month, day := commit.Author.When.Date()
	hour, minute, sec := commit.Author.When.Clock()
	stats, err := commit.Stats()
	if err != nil {
		fmt.Println("Couldn´t get the stats of the commit.")
		os.Exit(1)
	}
	return &utils.Commit{
		Hash: commit.Hash,
		When: fmt.Sprintf("%s %s %d %d %d:%d:%d",
			commit.Author.When.Weekday(), month, day, year, hour, minute, sec),
		Message: commit.Message,
		Stats:   stats,
	}
}
