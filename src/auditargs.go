package gitsearch

import (
	"flag"
	"fmt"
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
		err := fmt.Errorf("Argument 'giturl' should not be empty")
		return err
	}

	if "" == *gsArgs.rulesFile {
		flag.PrintDefaults()
		err := fmt.Errorf("Argument 'rulesfile' should not be empty")
		return err
	}

	return nil
}
