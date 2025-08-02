package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		printUsage()
		return
	}

}

func printUsage() {
	fmt.Println("./mywget [URL] [Recursion Length]")
}
