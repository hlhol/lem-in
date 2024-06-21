package maze

import (
	"fmt"
	"os"
	"strings"
)

// Load rooms for maze lines
func (maze *Maze) Load() {
	// Initialize the slices to store the names of the rooms, connections, beginning and destination rooms
	roomNames := []string{}
	connections := []string{}
	beginConnRmNames := []string{}
	destConnRmNames := []string{}

	// Loop through the roomsandConnections array to extract the names of the rooms
	for i := 0; i < len(maze.lines); i++ {
		for j := 0; j < len(maze.lines[i]); j++ {
			// If there's a space, split the string and store the name of the room in roomNames slice
			if maze.lines[i][j] == ' ' {
				roomName := strings.Split(maze.lines[i], " ")[0]
				roomNames = append(roomNames, roomName)

				break // because there are 2 spaces
			}
		}
	}

	// Loop through the roomsandConnections array to extract the connections between the rooms
	for i := 0; i < len(maze.lines); i++ {
		for j := 0; j < len(maze.lines[i]); j++ {
			if maze.lines[i][j] == '-' {
				connections = append(connections, maze.lines[i])

				beginDestSlice := strings.Split(maze.lines[i], "-")
				// If the beginning room is the ending room, add the destination room to beginConnRmNames
				if beginDestSlice[0] == maze.endName {
					beginConnRmNames = append(beginConnRmNames, beginDestSlice[1])
					destConnRmNames = append(destConnRmNames, beginDestSlice[0])
				}
				// Add the beginning and destination rooms to their respective slices
				beginConnRmNames = append(beginConnRmNames, beginDestSlice[0])
				destConnRmNames = append(destConnRmNames, beginDestSlice[1])
			}
		}
	}

	// Loop through the beginConnRmNames slice to check if any room is linked to itself
	for i := 0; i < len(beginConnRmNames); i++ {
		if beginConnRmNames[i] == destConnRmNames[i] {
			fmt.Println("ERROR\nSome Rooms Linked to Themselves")
			os.Exit(0)

		}
	}

	// Add the startName room to the antFarmRooms
	maze.addRoom(nil, maze.startName, beginConnRmNames, destConnRmNames)

	return
}

func (maze *Maze) addRoom(root *room, rmToAddName string, beginConnRmNames, destConnRmNames []string) *room {
	var roomToAdd *room
	var test []*room

	endroomparent := false
	anotherbool := false

	maze.duplicateRooms = append(maze.duplicateRooms, rmToAddName)

	if rmToAddName == maze.endName { // end room / base case
		test = append(test, root)
		if maze.endRoom == nil {
			maze.endRoom = &room{
				parent:   test,
				name:     rmToAddName,
				children: nil,
				occupied: false,
			}
		} else {
			parentnumber := len(maze.endRoom.parent)
			for t := 0; t < parentnumber; t++ {
				if maze.endRoom.parent[t].name == root.name {
					maze.endRoom.parent[t].children = root.children
					maze.endRoom.parent[t].parent = root.parent
					endroomparent = true
				} else if t == parentnumber-1 && !endroomparent {
					maze.endRoom.parent = append(maze.endRoom.parent, root)
				}
			}
		}

		countofparents := len(maze.endRoom.parent)
		if !endroomparent {
			maze.endRoom.parent[countofparents-1].children = append(maze.endRoom.parent[countofparents-1].children, maze.endRoom)
		}

		return maze.endRoom

	} else if rmToAddName == maze.startName { // startName Room special case
		maze.startRoom = &room{
			parent:   nil,
			name:     rmToAddName,
			occupied: true,
		}
		maze.findChildren(maze.startRoom, rmToAddName, beginConnRmNames, destConnRmNames)
		maze.Farm = append(maze.Farm, *maze.startRoom)
	} else {

		test = append(test, root)
		if len(dup(maze.duplicateRooms)) != 0 {
			for i := 0; i < len(dup(maze.duplicateRooms)); i++ {
				if (dup(maze.duplicateRooms)[i]) == rmToAddName {
					for t := 0; t < len(maze.Farm); t++ {
						if maze.Farm[t].name == rmToAddName {
							if roomToAdd == nil {
								roomToAdd = &room{
									parent:   test,
									name:     rmToAddName,
									occupied: false,
								}
							}
							roomToAdd.parent = maze.Farm[t].parent
							roomToAdd.children = maze.Farm[t].children
							numberofparents := len(roomToAdd.parent)
							for u := 0; u < len(roomToAdd.parent); u++ {
								if roomToAdd.parent[u].name == root.name {
									roomToAdd.parent[u].parent = root.parent
									roomToAdd.parent[u].children = root.children
									anotherbool = true

								} else if u == numberofparents-1 && !anotherbool {
									roomToAdd.parent = append(roomToAdd.parent, root)
									anotherbool = true
									break

								}
							}

							break

						}
					}
					break
				} else {
					roomToAdd = &room{
						parent:   test,
						name:     rmToAddName,
						occupied: false,
					}
				}
			}
		} else {
			roomToAdd = &room{
				parent:   test,
				name:     rmToAddName,
				occupied: false,
			}
		}

		countofparents := len(roomToAdd.parent)
		countofchildrens := 0
		var somebool = true
		countofchildrens = len(roomToAdd.parent[countofparents-1].children)
		if len(roomToAdd.parent[countofparents-1].children) == 0 {
			roomToAdd.parent[countofparents-1].children = append(roomToAdd.parent[countofparents-1].children, roomToAdd)
		} else {
			for r := 0; r < countofchildrens; r++ {

				if roomToAdd.parent[countofparents-1].children[r].name == roomToAdd.name {
					roomToAdd.parent[countofparents-1].children[r] = roomToAdd
					somebool = false
				}
			}
			if somebool {
				roomToAdd.parent[countofparents-1].children = append(roomToAdd.parent[countofparents-1].children, roomToAdd)
			}

		}

		maze.findChildren(roomToAdd, rmToAddName, beginConnRmNames, destConnRmNames)
		maze.Farm = CheckFarmDup(maze.Farm, roomToAdd.name)
		maze.Farm = append(maze.Farm, *roomToAdd)
	}

	return roomToAdd
}

func (maze *Maze) findChildren(roomToAdd *room, rmToAddName string, beginConnRmNames, destConnRmNames []string) {
	for c := 0; c < len(beginConnRmNames); c++ {
		beginRmName := beginConnRmNames[c]
		if maze.startName == destConnRmNames[c] {
			beginConnRmNames[c], destConnRmNames[c] = destConnRmNames[c], beginConnRmNames[c]
		}
		if rmToAddName != maze.startName {
			for m := 0; m < len(maze.startRoom.children); m++ {
				if maze.startRoom.children[m].name == destConnRmNames[c] && maze.startName != beginConnRmNames[c] {
					beginConnRmNames[c], destConnRmNames[c] = destConnRmNames[c], beginConnRmNames[c]
				}
			}
		}
		// Check if the beginning room of the current connection is the room to add.
		if beginRmName == rmToAddName {
			// If the above condition is true, retrieve the name of the destination room and add a new room to the ant farm.
			destRmName := destConnRmNames[c]
			maze.addRoom(roomToAdd, destRmName, beginConnRmNames, destConnRmNames)
		}
	}
}

func (maze *Maze) SwapRooms() {
	farmLength := len(maze.Farm) - 1

	for i := 0; i < (farmLength+1)/2; i++ {

		maze.Farm[i], maze.Farm[farmLength-i] = maze.Farm[farmLength-i], maze.Farm[i]
	}
}

func (maze *Maze) CreatingAnts() {
	maze.ants = make([]ant, maze.antCount)

	for i := 0; i < maze.antCount; i++ {
		var newAnt ant

		newAnt.id = i + 1
		newAnt.curRoom = maze.startRoom

		maze.ants[i] = newAnt
	}

	return
}

func dup(s []string) []string {
	var result []string
	duplicate := make(map[string]bool)

	for _, str := range s {
		if duplicate[str] {
			result = append(result, str)
		} else {
			duplicate[str] = true
		}
	}

	return result
}

func CheckFarmDup(s []room, t string) []room {
	for k := 0; k < len(s); k++ {
		if s[k].name == t {
			s = append(s[:k], s[k+1:]...)
			k--
		}
	}

	return s
}
