package main

import (
	"fmt"

	gitsearch "github.com/r-pai/go-gitaudit/src"
)

func main() {

	gsArgs := &gitsearch.GSArgs{}
	err := gsArgs.Validate()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	gitsearch.Start(gsArgs)
}
