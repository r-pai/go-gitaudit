package gitsearch

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"sync"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//GitBranch contains details about branch and its search options
type GitBranch struct {
	branchName          plumbing.ReferenceName
	gitRepo             *git.Repository
	searchAllBranchOpt  BranchSearchOptions
	searchCurrBranchOpt BranchSearchOptions
}

//SearchDetails contains search details
type SearchDetails struct {
	branch     *GitBranch
	commit     *object.Commit
	patch      *object.Patch
	prevCommit *object.Commit
}

//SetSearchOpts sets the branch search options
func (branch *GitBranch) SetSearchOpts(gitRepo *GitRepository) bool {

	searchBranch := false
	for _, s := range gitRepo.SearchOpts {
		if allBranch == s.Branch {
			branch.searchAllBranchOpt = s
			branch.searchAllBranchOpt.Branch = allBranch
			if searchBranch {
				break
			}
			searchBranch = true
		} else if branch.branchName.String() == s.Branch {
			branch.searchCurrBranchOpt = s
			if searchBranch {
				break
			}
			searchBranch = true
		}
	}

	return searchBranch
}

//Search searches as per the rules through all commit differences
func (branch *GitBranch) Search() {

	//retrieve the HEAD reference
	ref, err := branch.gitRepo.Head()
	if err != nil {
		fmt.Println("Error: Unable to retrive HEAD for branch ", branch.branchName.String())
		return
	}

	searchChan, wg := branch.createPatchSearchers()

	var prevCommit *object.Commit
	prevCommit = nil

	//Retrive the commit history
	cIter, err := branch.gitRepo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	err = cIter.ForEach(func(c *object.Commit) error {

		//Creating a patch between current and previous commits
		if prevCommit != nil {

			patch, err := c.Patch(prevCommit)
			if err != nil {
				fmt.Println("Error: Unable to create patch for commits ", prevCommit.Hash.String(), c.Hash.String())
				return nil
			}
			//fmt.Println("Patch:",patch)
			noOfCommits++
			searchChan <- SearchDetails{
				branch:     branch,
				commit:     c,
				patch:      patch,
				prevCommit: prevCommit}

		}
		prevCommit = c
		return nil
	})

	close(searchChan)
	wg.Wait()
}

//createPatchSearchers creates go routines, which will search on the diff
//Workers are created for each branch
func (branch *GitBranch) createPatchSearchers() (chan SearchDetails, *sync.WaitGroup) {
	searchChan := make(chan SearchDetails, 100)
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 1; w <= workers; w++ {
		go branch.patchSearcher(w, &wg, searchChan)
	}

	return searchChan, &wg
}

//patchSearcher searches the patch aginst all the rules and prints the output
func (branch *GitBranch) patchSearcher(id int, wg *sync.WaitGroup, searchChan <-chan SearchDetails) {
	for search := range searchChan {

		var searchBranchOpts [2]*BranchSearchOptions
		searchBranchOpts[0] = &search.branch.searchAllBranchOpt
		searchBranchOpts[1] = &search.branch.searchCurrBranchOpt

		for sBranchIndex := range searchBranchOpts {

			//For allBranch rule(to be checked in all branches), we need only check a patche
			//once and need not check in rest of the branches
			//To keep track we would be using a allBranchallBranchCommitTracker,
			//this has the sha256 as hask as the key and value its key in string as value
			if searchBranchOpts[sBranchIndex].Branch == allBranch {

				sumStr := search.prevCommit.Hash.String() + search.commit.Hash.String()
				sum := [sha256.Size]byte(sha256.Sum256([]byte(sumStr)))

				_, ok := allBranchCommitTracker[sum]
				if ok {
					continue
				}
				allBranchCommitTracker[sum] = sumStr
			}

			for key, searchString := range searchBranchOpts[sBranchIndex].RuleSet.Rules {
				//fmt.Println(id, ": Searching ", searchString, "in ", c.branch, "...")
				match, fPath, matchStr := search.branch.SearchInPatch(searchString, search.patch)
				if match {
					safe.Inc()

					report := AuditReport{
						RuleName:  searchBranchOpts[sBranchIndex].RuleSet.RuleName,
						Rule:      key,
						Hash:      search.commit.Hash.String(),
						Time:      search.commit.Committer.When.Format("Mon Jan _2 15:04:05 2006"),
						Committer: search.commit.Committer.Name,
						FilePath:  fPath,
						Branch:    search.branch.branchName.String(),
						Match:     matchStr,
					}
					report.print()
				}

			}
		}
	}
	wg.Done()
}

//SearchInPatch searches the searchString in patch , if a match is found we will just come out.
//searchString can be a regular experssion
func (branch *GitBranch) SearchInPatch(searchString string, p *object.Patch) (bool, string, string) {
	fPath := ""
	for _, f := range p.FilePatches() {
		if !f.IsBinary() {
			for _, c := range f.Chunks() {

				re := regexp.MustCompile(searchString)
				match := re.Find([]byte(c.Content()))
				if match != nil {
					from, to := f.Files()
					if from != nil {
						fPath = from.Path()
					} else {
						if to != nil {
							fPath = to.Path()
						}
					}

					return true, fPath, string(match)
				}
			}
		}
	}

	return false, "", ""
}
