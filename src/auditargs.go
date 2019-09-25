package gitsearch

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	searchOptsFile = "searchopts.json"
)

//GSArgs used to store arguments received from the user
type GSArgs struct {
	gitURL    *string
	rulesFile *string
	json      *bool
	worker    *int
}

//Validate vaildates the arguments
func (gsArgs *GSArgs) Validate() error {

	gsArgs.gitURL = flag.String("giturl", "", "git repository URL")
	gsArgs.rulesFile = flag.String("rulesfile", "", fmt.Sprintf("rules file path. \nFor json format refer file defaultrule.json."))
	gsArgs.json = flag.Bool("json", true, fmt.Sprintf("Output format to be json (true or false)."))
	gsArgs.worker = flag.Int("worker", 1, "number of workers for parallel processing (max "+strconv.Itoa(maxWorker)+")")

	flag.Parse()

	if *gsArgs.worker > maxWorker {
		*gsArgs.worker = workers
	}

	if "" == *gsArgs.gitURL {
		flag.PrintDefaults()
		err := fmt.Errorf("argument giturl should not be empty")
		return err
	}

	if "" == *gsArgs.rulesFile {
		insDir, _ := os.Getwd()
		defRuleFile := insDir + "/" + searchOptsFile
		if _, err := os.Stat(defRuleFile); err != nil {
			if os.IsNotExist(err) {
				err := fmt.Errorf("default rulesfile %s not found", searchOptsFile)
				return err
			}
		}
		*gsArgs.rulesFile = defRuleFile
	}

	return nil
}
