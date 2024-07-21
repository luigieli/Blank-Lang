package main

import (
	"blank/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Blank Lang: Welcome %s! \n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
