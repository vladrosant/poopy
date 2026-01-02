package main

import (
	"flag"
	"fmt"
	"os"
)

func handleUser() {
	addCmd := flag.NewFlagSet("user", flag.ExitOnError)

	name := addCmd.string("n", "name", "Name of the user")
	age := addCmd.int("a", "age", "Age of the user")
	id := addCmd.int("id", "id", "ID of the user")

	flag.Parse()

	fmt.Printf("User %s, %d years old, ID %d", *name, *age, *id)

}

func main() {
	fmt.Println("Number of arguments:", len(os.Args))
	fmt.Println("Arguments:", os.Args)

	for i, arg := range os.Args {
		fmt.Printf("Argument %d: %s\n", i, arg)
	}

	command := os.Args[1]

	if command == "user" {
		handleUser()
	}
}
