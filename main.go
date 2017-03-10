package main

import "os"

// Name of this program
var Name string

// Version of this program
var Version string

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
