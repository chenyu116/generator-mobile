package main

import (
	"fmt"

	"github.com/chenyu116/generator-mobile/cli"
)

var (
	_version_   = ""
	_branch_    = ""
	_commitId_  = ""
	_buildTime_ = ""
)

func main() {
	fmt.Printf("Version: %s, Branch: %s, Build: %s, Build time: %s\n\n",
		_version_, _branch_, _commitId_, _buildTime_)
	cli.Start()
}
