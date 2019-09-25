package gitsearch

import (
	"crypto/sha256"
	"fmt"
	"os"
)

const (
	allBranch    = "all"
	masterBranch = "refs/heads/master"
	maxWorker    = 10
)

var (
	allBranchCommitTracker map[[sha256.Size]byte]string
	workers                = 1
	searchMaster           = true
	noOfCommits            = 0
	safe                   = SafeCounter{}
	jsonOutput             = true
)

//Start is the entry point for git search
func Start(gs *GSArgs) {

	repoSearchOptions, err := GetSearchOptions(*gs.rulesFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Lets Git Search
	if repoSearchOptions != nil {

		gitSearch(*gs.gitURL, repoSearchOptions)
		fmt.Println("Audited", noOfCommits, "commits, found", safe.GetViolationCount(), "violations.")
		fmt.Println("")
	}
}

//gitSearch searchs the repository
//1. Creates a temporay folder.
//2. Clone the repository
//3. Loop all the branches
//4. Each branch has a set of workers(routines), which does the
//   string searching in parallel.
//5. ony rules pplicable to a branch is applied.
//6. allbranch rule will be applied to all branches
//7. The output are printed to the stdout.
func gitSearch(gitURL string, repoSearhOpts *RepoSearchOptions) {

	gitRepo := &GitRepository{}
	err := gitRepo.CloneRepo(gitURL)
	defer os.RemoveAll(gitRepo.cloneDir)
	fmt.Println()

	if err != nil {
		fmt.Println("Cloning of repository ", gitURL, " failed, reason being ", err.Error())
		return
	}

	gitRepo.SearchOpts = repoSearhOpts.BranchSearchOptions[:]
	gitRepo.Search()
}
