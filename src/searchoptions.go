package gitsearch

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//RepoSearchOptions contains
type RepoSearchOptions struct {
	BranchSearchOptions []BranchSearchOptions `json:"repository"`
}

//BranchSearchOptions has the search options for each branch
type BranchSearchOptions struct {
	Branch  string  `json:"branch"`
	RuleSet RuleSet `json:"ruleset"`
}

//RuleSet contains rule details of a branch
type RuleSet struct {
	RuleName string            `json:"rulename"`
	Rules    map[string]string `json:"rules"`
}

//GetSearchOptions gets the search options from the seachoption file
func GetSearchOptions(searchOptFile string) (*RepoSearchOptions, error) {

	if _, err := os.Stat(searchOptFile); os.IsNotExist(err) {
		return nil, err
	}

	jsonFile, err := os.Open(searchOptFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	repoSearchOptions := &RepoSearchOptions{}
	if err := json.Unmarshal(byteValue, repoSearchOptions); err != nil {
		return nil, err
	}

	//If we are scanning only master branch we will clone the repo in memory
	for _, bOpt := range repoSearchOptions.BranchSearchOptions {
		if bOpt.Branch != masterBranch {
			searchMaster = false
			break
		}
	}

	return repoSearchOptions, nil
}
