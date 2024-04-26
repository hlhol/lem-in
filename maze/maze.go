package maze

import (
	"bufio"
	"strconv"
	"strings"
)

type Maze struct {
    antCount int
    rooms []Room
} 

type Room struct {
    connections []Room
    distance int

    name string
    occupancies int

    isEnd bool
}

func Load(inputStream *bufio.Reader) *Maze {

    maze := &Maze{}

    // first line (ant count)
    line, err := inputStream.ReadString('\n')
    if (err != nil || len(line) == 0) {
        return nil
    }

    maze.antCount, err = strconv.Atoi(line[:len(line)-1])
    if (err != nil) {
        return nil
    }

    // second line (##start)
    line, err = inputStream.ReadString('\n')
    if (err != nil || len(line) == 0 || line != "##start\n") {
        return nil
    }
    
    if !maze.loadRooms(inputStream) {
        return nil
    }
    if !maze.mapFarm(inputStream) { // connect rooms and calculate distance
       return nil 
    }

    return maze
}

func (maze *Maze) Solve() {


}

func (maze *Maze) loadRooms(inputStream *bufio.Reader) bool {

    var line string
    var err error
    var roomName string

    readName := func() {

        roomName  = ""

        line, err = inputStream.ReadString('\n')
        if err != nil || len(line) == 0 {
            // print("can't read")
            return
        }

        // print(line)

        parts := strings.Split(line[:len(line)-1], " ")
        if len(parts) != 3 {
            // println("parts error ", line[:len(line)-1], len(parts))
            return
        }

        _, err1 := strconv.Atoi(parts[1])
        _, err2 := strconv.Atoi(parts[2])

        // anti fun measure
        if err1 != nil || err2 != nil {
            // print(parts[1], err1 == nil, parts[2], err2 == nil)
            return
        }

        roomName = parts[0]
    }


    // first Room
    readName()
    if roomName == "" {
        return false
    }

    room := Room{ name: roomName, occupancies: maze.antCount }
    maze.rooms = append(maze.rooms, room)

    // rest of the farm
    readName()
    for err == nil {
        if (line == "##end\n") {
            // print("end?")
            maze.rooms[len(maze.rooms)-1].isEnd = true
            return true

        } else if (roomName == "") {
            // print("loop empty")
            return false
        }

        room = Room{ name: roomName, occupancies: 0 }
        maze.rooms = append(maze.rooms, room)

        readName()
    }

    // print("bad exit")
    return false
}

func (maze *Maze) mapFarm(inputStream *bufio.Reader) bool {

    for _, v := range maze.rooms {
        println(v.name, v.occupancies, v.isEnd)
    }


    return true
}

