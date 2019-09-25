package gitsearch

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const (
	//CloneDir is the temporary directory name used while
	//Checking out git repository
	CloneDir = "gitaudit-"
)

//GitRepository has the details about the repository to be searched
// and search options
type GitRepository struct {
	repoURL    string
	repo       *git.Repository
	cloneDir   string
	SearchOpts []BranchSearchOptions
}

//CloneRepo clones the repository
func (gitRepo *GitRepository) CloneRepo(repoURL string) error {

	if searchMaster {
		return gitRepo.cloneRepoInMemory(repoURL)
	}

	return gitRepo.cloneRepoInDir(repoURL)

}

//cloneRepoInMemory clones the repository in memory
func (gitRepo *GitRepository) cloneRepoInMemory(repoURL string) error {

	fmt.Println("Cloning repository ", repoURL)

	var err error
	gitRepo.repo, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	return err
}

//cloneRepoInDir clones the git repository into a local temporary directory
func (gitRepo *GitRepository) cloneRepoInDir(repoURL string) error {

	dir, err := ioutil.TempDir("", CloneDir)

	if err != nil {
		return err
	}

	gitRepo.cloneDir = dir
	fmt.Println("Cloning ", repoURL, " in ", dir)

	gitRepo.repo, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}

func (gitRepo *GitRepository) searchInBranch(branchName plumbing.ReferenceName) *GitBranch {

	branch := &GitBranch{branchName, gitRepo.repo, BranchSearchOptions{}, BranchSearchOptions{}}

	if branch.SetSearchOpts(gitRepo) {
		branch.Search()
	} else {
		fmt.Println("No rule to search in branch ", branch.branchName)
	}

	return branch
}

//Search starts the search in gitrepo
func (gitRepo *GitRepository) Search() {

	//Search Only Master Branch
	if searchMaster {
		gitRepo.searchInBranch(masterBranch)
		return
	}

	//Loop though all branches
	refs, err := gitRepo.repo.References()
	if err != nil {
		fmt.Println("Error searching for branches, error ", err.Error())
		return
	}

	start := time.Now()
	time.Sleep(time.Second * 2)
	allBranchCommitTracker = make(map[[sha256.Size]byte]string)
	err = refs.ForEach(func(ref *plumbing.Reference) error {

		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		w, err := gitRepo.repo.Worktree()
		if err != nil {
			fmt.Println("Error getting worktree ", err.Error())
			return nil
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash: ref.Hash(),
		})
		if err != nil {
			fmt.Println("Unable to checkout branch ", ref.Name(), err.Error())
			return nil
		}

		gitRepo.searchInBranch(ref.Name())
		return nil
	})

	fmt.Printf("Time taken for searching the repo %s\n", time.Since(start))
}
