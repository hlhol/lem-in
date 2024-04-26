package main

import (
	"fmt"
	"os"
    "bufio"
    "lem/maze"
)


func main()  {
    
    if len(os.Args) != 2 {
        fmt.Println("only file please")
        return
    }

    filename := os.Args[1]

    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("can't open file" + filename)
        return
    }

    reader := bufio.NewReader(file)

    m := maze.Load(reader)
    if (m == nil) {
        fmt.Println("Bad input file")
        return
    }

    m.Solve()

}
