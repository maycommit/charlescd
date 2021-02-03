package main

import (
	"log"

	"github.com/maycommit/circlerr/cmd/cli/commands/circle"
	"github.com/maycommit/circlerr/cmd/cli/commands/project"
)

func main() {

	circle.CmdCircleRoot.AddCommand(circle.CircleList())
	circle.CmdCircleRoot.AddCommand(circle.CircleTree())

	project.CmdProjectRoot.AddCommand(project.ProjectList())

	rootCmd.AddCommand(circle.CmdCircleRoot)
	rootCmd.AddCommand(project.CmdProjectRoot)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
