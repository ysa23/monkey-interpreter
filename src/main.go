package main

import (
	"fmt"
	"os"
	"os/user"
	"./repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! I'll bet your wondering why the red suit\n", user.Name)
	fmt.Printf("Lets get the ball rolling - you can now start typing commands:\n")
	repl.Start(os.Stdin, os.Stdout)
}
