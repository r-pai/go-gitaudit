package gitsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

//AuditReport contains reporting fields to be
type AuditReport struct {
	RuleName  string `json:"rulename"`
	Rule      string `json:"rule"`
	Branch    string `json:"branch"`
	Hash      string `json:"hash"`
	Time      string `json:"time"`
	Committer string `json:"committer"`
	FilePath  string `json:"file"`
	Match     string `json:"match"`
}

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	branch         string
	noOfViolations int
	mux            sync.Mutex
}

// Inc increments the number of Violations
func (c *SafeCounter) Inc() {
	c.mux.Lock()
	c.noOfViolations++
	c.mux.Unlock()
}

// GetViolationCount returns the no of Violations
func (c *SafeCounter) GetViolationCount() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.noOfViolations
}

// Reset resets the safe values
func (c *SafeCounter) Reset() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.noOfViolations = 0
}

//print prints the output
func (report *AuditReport) print() {
	if jsonOutput {
		data, err := json.Marshal(*report)
		if err != nil {
			log.Println("JSON Marshal error: ", err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, data, "", "\t")
		if error != nil {
			log.Println("JSON parse error: ", error)
			return
		}

		fmt.Println(prettyJSON.String())
	} else {
		fmt.Println(
			"RuleName:", report.RuleName,
			"\nRule\t:", report.Rule,
			"\nBranch\t:", report.Branch,
			"\nHash\t:", report.Hash,
			"\nTime\t:", report.Time,
			"\nCommiter:", report.Committer,
			"\nFile\t:", report.FilePath,
			"\nMatch\t:", report.Match,
			"\n ")
	}

}
