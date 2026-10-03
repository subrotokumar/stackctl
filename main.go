/*
Copyright © 2025 Subroto Kumar <subrotokumar@outlook.in>
*/
package main

import (
	"github.com/subrotokumar/stackctl/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}
