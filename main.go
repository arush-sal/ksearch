package main

import "github.com/arush-sal/ksearch/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
