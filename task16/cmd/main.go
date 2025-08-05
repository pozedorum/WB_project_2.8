package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"task16/internal/wget"
)

func main() {
	urlFrom, depth, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = wget.Wget(urlFrom, depth)
	if err != nil {
		fmt.Println(err)
	}
}

func printUsage() {
	fmt.Println("./mywget [URL] [Recursion Length]")
}

func parseArgs() (urlFrom string, depth int, err error) {
	if len(os.Args) != 3 {
		printUsage()
		err = errors.New("wrong count of arguments")
		return
	}
	urlFrom = os.Args[1]
	depth, err = strconv.Atoi(os.Args[2])
	return
}
