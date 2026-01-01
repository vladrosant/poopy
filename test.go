package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Number of arguments:", len(os.Args))
	fmt.Println("Arguments:", os.Args)

	for i, arg := range os.Args {
		fmt.Printf("Argument %d: %s\n", i, arg)
	}
}
