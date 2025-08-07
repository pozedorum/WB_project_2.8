package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pozedorum/WB_project_2/task17/telnet"
)

func main() {
	timeout := flag.Int64("timeout", 10, "connection timeout in seconds")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		printUsage()
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	if err := telnet.RunTelnetClient(host, port, time.Duration(*timeout)*time.Second); err != nil {
		fmt.Printf("Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: ./mytelnet [--timeout=10] host port")
	fmt.Println("Example: ./mytelnet --timeout=5 example.com 23")
}
