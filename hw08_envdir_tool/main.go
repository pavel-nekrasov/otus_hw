package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Wrong number of arguments. Expected: <path to env dir> <command to run> [arg1] [arg2] ...")
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := RunCmd(os.Args[2:], env)
	os.Exit(exitCode)
}
