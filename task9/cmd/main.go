package main

import (
	"fmt"

	"github.com/pozedorum/WB_project_2.8/task9"
)

func main() {
	var input string
	fmt.Scan(&input)

	output, err := task9.UnpackString(input)
	if err != nil {
		fmt.Printf("error: UnpackString - %v\n", err)
	} else {
		fmt.Println(output)
	}
}
