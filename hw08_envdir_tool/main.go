package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := initCli()
	envDirPath := args[1]

	env, err := ReadDir(envDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(args[2:], env))
}

func initCli() []string {
	const MinArgs = 2

	if args := os.Args; len(args) > MinArgs {
		return args
	}

	fmt.Printf("Please give %v or more arguments.\n", MinArgs)
	os.Exit(0)
	return nil
}
