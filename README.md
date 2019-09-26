# go-gitaudit
  This project is inspired from [trufflehog](https://github.com/dxa4481/truffleHog) and the default basic search expressions are from trufflehog. 
  
  go-gitaudit searches deep in commit history, in any branch or entire repository to find any thing you looking for.
 
  **Features**
  - Uses golang go-git package
  - In Memory search for master branch
  - Configuration using json file.
  - Branch wise search support.
  - Faster, as search is done in parallel.

 
# Installation
  go-gitaudit is built using Googles [golang](https://golang.org) and require golang installed in the system to run.
  
  Dependency package
  - gopkg.in/src-d/go-git.v4
  
```sh
$ go get -u github.com/r-pai/go-gitaudit/...
$ cd $GOPATH/src/github.com/r-pai/go-gitaudit
$ go install
$ go-gitaudit --giturl <url>
```
# Help 
```
Usage : 
go-gitaudit 
  -giturl string
    	git repository URL
  -json
    	Output format to be json (true or false). (default true)
  -rulesfile string
    	rules file path. 
    	For json format refer file defaultrule.json.
  -worker int
    	number of workers for parallel processing (max 10) (default 1
```

# Example

Basic command for go-gitaudit. 

```
$./go-gitaudit --giturl=<url>
```

go-gitaudit with new Rules File
```
$./go-gitaudit --giturl=<url> --rulesfile=<rulesfilepath>
```

go-gitaudit output format is by default json (only for a rule). To change 
```
$./go-gitaudit --giturl=<url> --json=false
```

# rulesfile format

```
{
  "searchoptions": [
    {
      "branch": "all",
      "ruleset": {
        "rulename": "DefaultRule",
        "rules": {
          "rule1":   "regular[a-z]expression1",
          "rule2":    "searchstring"
        }
      }
    },
    {
      "branch": "refs/remotes/origin/dev",
      "ruleset": {
        "rulename": "DevRule",
        "rules": {
          "devrule1":   "regular[a-z]expression2"
        }
      }
    }
  ]
}
```


# Todos
  - Add more features in trufflehog to go-gitaudit (entropy,...)
  - write test cases

# Issues
  One of the issue encountered is, when a diff has more than 55K lines.This issue is not of 'gotrufflehog', its a package used go-git panics.
  More about the issue 
  - 'https://github.com/src-d/go-git/issues/973'
  - 'https://github.com/sergi/go-diff/issues/89' 
  