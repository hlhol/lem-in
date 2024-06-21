package maze

import (
	"fmt"
	"os"
	"strings"
)

// Load rooms form maze lines
func (maze *Maze) Load() {

	var roomNames []string

	var sourceRoomConnection []string
	var destRoomConnection []string

	// go through maze lines, and the rooms name and connections
	for _, line := range maze.lines {
		for _, char := range line {

			if char == ' ' {
				roomName := strings.Split(line, " ")[0]
				roomNames = append(roomNames, roomName)

				break

			} else if char == '-' {

				beginDestSlice := strings.Split(line, "-")

				if beginDestSlice[0] == maze.endName {
					sourceRoomConnection = append(sourceRoomConnection, beginDestSlice[1])
					destRoomConnection = append(destRoomConnection, beginDestSlice[0])
				} else {

					sourceRoomConnection = append(sourceRoomConnection, beginDestSlice[0])
					destRoomConnection = append(destRoomConnection, beginDestSlice[1])
				}

				break
			}
		}
	}

	// check if any room link to itself
	for i, roomName := range sourceRoomConnection {
		if roomName == destRoomConnection[i] {
			fmt.Println("ERROR\nSome Rooms Linked to Themselves")
			os.Exit(0)

		}
	}

	maze.addRoom(nil, maze.startName, sourceRoomConnection, destRoomConnection)
}

func (maze *Maze) addRoom(root *room, roomName string, sourceConnections, destinationConnections []string) *room {
	var newRoom *room
	maze.addedRooms = append(maze.addedRooms, roomName)

	endParent := false
	foundParent := false

	maze.addedRooms = append(maze.addedRooms, roomName)

	if roomName == maze.endName { // end room / base case
		if maze.endRoom == nil {

			maze.endRoom = &room{
				parent:   []*room{root},
				name:     roomName,
				children: nil,
				occupied: false,
			}

		} else {

			parentCount := len(maze.endRoom.parent)
			for t := 0; t < parentCount; t++ {
				// end room parent exists
				if maze.endRoom.parent[t].name == root.name {
					maze.endRoom.parent[t].children = root.children
					maze.endRoom.parent[t].parent = root.parent

					endParent = true
				} else if t == parentCount-1 && !endParent {
					maze.endRoom.parent = append(maze.endRoom.parent, root)
				}
			}
		}

		parentCount := len(maze.endRoom.parent)
		if !endParent {
			maze.endRoom.parent[parentCount-1].children = append(maze.endRoom.parent[parentCount-1].children, maze.endRoom)
		}

		return maze.endRoom

	} else if roomName == maze.startName { // startName Room special case

		maze.startRoom = &room{
			parent:   nil,
			name:     roomName,
			occupied: true,
		}

		maze.findChildren(maze.startRoom, roomName, sourceConnections, destinationConnections)
		maze.Farm = append(maze.Farm, *maze.startRoom)

	} else {

		if len(findDuplicates(maze.addedRooms)) != 0 {
			for i := 0; i < len(findDuplicates(maze.addedRooms)); i++ {
				if (findDuplicates(maze.addedRooms)[i]) == roomName {

					for t := 0; t < len(maze.Farm); t++ {
						if maze.Farm[t].name == roomName {
							// if room was found in the farm. it will be added and deduced later

							if newRoom == nil {
								newRoom = &room{
									parent:   []*room{root},
									name:     roomName,
									occupied: false,
								}
							}

							newRoom.parent = maze.Farm[t].parent
							newRoom.children = maze.Farm[t].children

							parentCount := len(newRoom.parent)

							for u := 0; u < len(newRoom.parent); u++ {
								if newRoom.parent[u].name == root.name {
									newRoom.parent[u].parent = root.parent
									newRoom.parent[u].children = root.children
									foundParent = true

								} else if u == parentCount-1 && !foundParent {
									newRoom.parent = append(newRoom.parent, root)
									foundParent = true
									break

								}
							}

							break

						}
					}
					break
				} else {
					newRoom = &room{
						parent:   []*room{root},
						name:     roomName,
						occupied: false,
					}
				}
			}
		} else {
			newRoom = &room{
				parent:   []*room{root},
				name:     roomName,
				occupied: false,
			}
		}

		parentCount := len(newRoom.parent)
		childrenCount := len(newRoom.parent[parentCount-1].children)
		didNotFindChildren := true

		if len(newRoom.parent[parentCount-1].children) == 0 {
			newRoom.parent[parentCount-1].children = append(newRoom.parent[parentCount-1].children, newRoom)
		} else {
			for r := 0; r < childrenCount; r++ {

				if newRoom.parent[parentCount-1].children[r].name == newRoom.name {
					newRoom.parent[parentCount-1].children[r] = newRoom
					didNotFindChildren = false
				}
			}
			if didNotFindChildren {
				newRoom.parent[parentCount-1].children = append(newRoom.parent[parentCount-1].children, newRoom)
			}

		}
		maze.findChildren(newRoom, roomName, sourceConnections, destinationConnections)
		maze.Farm = removeDuplicates(maze.Farm, roomName)
		maze.Farm = append(maze.Farm, *newRoom)
	}

	return newRoom
}

func (maze *Maze) findChildren(roomToAdd *room, newRoom string, sourceConnections, destinationConnections []string) {
	for c := 0; c < len(sourceConnections); c++ {
		sourceRoom := sourceConnections[c]

		if maze.startName == destinationConnections[c] {
			sourceConnections[c], destinationConnections[c] = destinationConnections[c], sourceConnections[c]
		}

		if newRoom != maze.startName {
			for m := 0; m < len(maze.startRoom.children); m++ {
				if maze.startRoom.children[m].name == destinationConnections[c] && maze.startName != sourceConnections[c] {
					sourceConnections[c], destinationConnections[c] = destinationConnections[c], sourceConnections[c]
				}
			}
		}

		// If the source room is the current room. Add its destination room
		if sourceRoom == newRoom {
			destinationRoom := destinationConnections[c]
			maze.addRoom(roomToAdd, destinationRoom, sourceConnections, destinationConnections)
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
		maze.ants[i] = ant{
			id:      i + 1,
			curRoom: maze.startRoom,
		}
	}
}

func findDuplicates(s []string) []string {
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

func removeDuplicates(s []room, t string) []room {
	for k := 0; k < len(s); k++ {
		if s[k].name == t {
			s = append(s[:k], s[k+1:]...)
			k--
		}
	}

	return s
}
