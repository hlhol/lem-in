package main

import (
	"fmt"
	"lem/maze"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("only file please")
		return
	}

	filename := os.Args[1]
	//filename := "examples/example04.txt"

	m := maze.ReadFile(filename)
	m.Start()
	return

	for i := 0; i < 8; i++ {
		filename := fmt.Sprintf("examples/example0%d.txt", i)
		m := maze.ReadFile(filename)
		m.Start()
	}

}
