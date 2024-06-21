package main

import (
	"fmt"
	"lem/maze"
	"os"
)

func main() {

	if len(os.Args) != 2 && false {
		fmt.Println("only file please")
		return
	}

	//filename := os.Args[1]
	filename := "examples/example04.txt"

	m := maze.ReadFile(filename)
	m.Start()

	return
	// this will work after full clean up
	// it's not working right now because the maze.start change global and have side effects
	for i := 0; i < 8; i++ {
		//filename := fmt.Sprintf("examples/example0%d.txt", i)
		//maze.Start(filename)
	}

}
