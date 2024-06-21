package maze

import (
	"fmt"
)

// It takes a slice of ants as input
func (maze *Maze) walk() {
	var allpassed = true
	var EndRoomStr EndRoomUsed
	for i := 0; i < len(maze.ants); i++ {

		// If an ant has already reached the last room, skip to the next ant
		if maze.ants[i].curRoom.name == maze.endRoom.name {
			continue
		}
		// If an ant is not in the first room, mark the previous room as unoccupied
		if maze.ants[i].curRoom.name != maze.startRoom.name {
			maze.ants[i].curRoom.occupied = false
		}
		// If the next room an ant wants to move to is already occupied, skip to the next ant
		if maze.ants[i].pathOfAnt[maze.ants[i].step].occupied {
			continue
		}
		// If an ant has already used a particular path to reach the last room, skip to the next ant
		if EndRoomStr.used && EndRoomStr.whichPath[0].name == maze.ants[i].pathOfAnt[0].name && maze.ants[i].pathOfAnt[maze.ants[i].step].name == maze.endRoom.name {
			continue
		}
		// Move an ant to the next room in its path
		maze.ants[i].curRoom = maze.ants[i].pathOfAnt[maze.ants[i].step]
		maze.ants[i].step++
		if maze.ants[i].curRoom.name != maze.endRoom.name {
			maze.ants[i].curRoom.occupied = true
		}
		// If an ant has reached the last room, store the path used by the ant
		if maze.ants[i].curRoom.name == maze.endRoom.name {
			EndRoomStr.used = true
			EndRoomStr.whichPath = maze.ants[i].pathOfAnt
		}
		// Update the allpassed variable to false if an ant has not yet reached the last room
		allpassed = false
		// Print the movement of an ant in the ant farm
		fmt.Print("L", maze.ants[i].id, "-", maze.ants[i].curRoom.name, " ")

	}
	// If all ants have reached the last room, return
	if allpassed {
		return
		// If some ants have not yet reached the last room, set the end room used by an ant to false
		// and call the walk function recursively
	} else {
		EndRoomStr.used = false
		EndRoomStr.whichPath = nil

		fmt.Println(" ")
		maze.walk()
	}
}

// FindAllPossiblePaths is a recursive function that searches for all possible paths from the starting room to the ending room.
func (maze *Maze) FindAllPossiblePaths(path []*room, currentRoom room, paths *[][]*room) {
	// if the current room is the last room, append the current path to the paths slice
	if currentRoom.name == maze.endRoom.name {
		// Check if the path goes back to the first room, and skip it if it does
		var skipPath bool
		for i := 0; i < len(path); i++ {
			if path[i].name == maze.startRoom.name {
				skipPath = true
				break
			}
		}

		if len(*paths) == 0 {
			*paths = append(*paths, nil)
		} else if (*paths)[len(*paths)-1] != nil {
			*paths = append(*paths, nil)
		}

		for i := 0; i < len(path); i++ {
			if !skipPath {
				(*paths)[len(*paths)-1] = append((*paths)[len(*paths)-1], path[i])
			} else {
				break
			}
		}
	}

	// Recursively explore all possible paths from the current room to other rooms
	for i := 0; i < len(currentRoom.children); i++ {
		var toContinue bool

		for k := 0; k < len(path); k++ {
			if path[k].name == currentRoom.children[i].name {
				toContinue = true
				break
			}
		}

		if !toContinue {
			pathToPass := path
			pathToPass = append(pathToPass, currentRoom.children[i])
			maze.FindAllPossiblePaths(pathToPass, *currentRoom.children[i], paths)
			pathToPass = path
		}
	}

	// remove any empty paths from the paths slice
	for i := 0; i < len(*paths); i++ {
		if (*paths)[i] == nil {
			*paths = append((*paths)[:i], (*paths)[i+1:]...)
		}
	}
}

func SortPaths(ways [][]*room) [][]*room {
	for i := 0; i < len(ways)-1; i++ {
		if len(ways[i]) > len(ways[i+1]) {
			ways[i], ways[i+1] = ways[i+1], ways[i]
		}
	}

	for k := 0; k < len(ways)-1; k++ {
		if len(ways[len(ways)-1]) < len(ways[k]) {
			ways[len(ways)-1], ways[k] = ways[k], ways[len(ways)-1]
		}
	}
	return ways
}

func (solution *solution) ClearPath(maze *Maze, ways [][]*room) [][]*room {
	var somebool = false
	var anotherbool = false
	solution.childrenOfFirstRoom = len(maze.startRoom.children)
	if solution.appendWays == nil {
		solution.appendWays = ways[0]
	}
	if solution.CombinatedRooms == nil {
		solution.CombinatedRooms = append(solution.CombinatedRooms, ways[0])
	}
	if solution.countloop == len(ways)-1 {
		return solution.CombinatedRooms
	}
	for i := 0; i < len(ways[solution.countloop+1]); i++ {
		for k := 0; k < len(solution.appendWays)-1; k++ {
			if ways[solution.countloop+1][i].name == solution.appendWays[k].name {
				somebool = true
			}
			if !somebool && i == len(ways[solution.countloop+1])-1 && k == len(solution.appendWays)-2 {
				solution.appendWays = append(solution.appendWays, ways[solution.countloop+1]...)
				anotherbool = true
			}
		}
	}
	if anotherbool {
		solution.CombinatedRooms = append(solution.CombinatedRooms, ways[solution.countloop+1])
	}
	solution.countloop++
	solution.ClearPath(maze, ways)
	return solution.CombinatedRooms
}

func (solution *solution) FirstChildren(maze *Maze, ways [][]*room) []Path {
	var PathStruct Path
	var PathStruct2 []Path
	solution.childrenOfFirstRoom = len(maze.startRoom.children)
	for i := 0; i < len(ways); i++ {
		for k := 0; k < solution.childrenOfFirstRoom; k++ {
			if ways[i][0] == maze.startRoom.children[k] {
				PathStruct.id = k
				PathStruct.paths = ways[i]
				PathStruct.intersect = true
				PathStruct2 = append(PathStruct2, PathStruct)

			}
		}
	}
	return PathStruct2
}

func SortedPaths(way []Path, idChildren int) []Path {
	var SepPath []Path
	for i := 0; i < len(way); i++ {
		if way[i].id == idChildren {
			SepPath = append(SepPath, way[i])
		}
	}
	return SepPath
}

func SortAgain(way [][]Path) [][]Path {
	for i := 0; i < len(way)-1; i++ {
		if len(way[i]) < len(way[i+1]) {
			way[i], way[i+1] = way[i+1], way[i]
		}
	}
	return way
}

func (solution *solution) AllCombinations(way [][]Path) [][]Path {
	var AnotherPath []Path
	var AnotherPath2 [][]Path
	if solution.childrenOfFirstRoom == 2 {
		for i := 0; i < len(way[0]); i++ {
			for k := 0; k < len(way[1]); k++ {
				AnotherPath = append(AnotherPath, way[0][i], way[1][k])
				AnotherPath2 = append(AnotherPath2, AnotherPath)
				AnotherPath = nil

			}
		}
	} else if solution.childrenOfFirstRoom == 3 {
		for i := 0; i < len(way[0]); i++ {
			for k := 0; k < len(way[1]); k++ {
				for l := 0; l < len(way[2]); l++ {
					AnotherPath = append(AnotherPath, way[0][i], way[1][k], way[2][l])

					AnotherPath2 = append(AnotherPath2, AnotherPath)
					AnotherPath = nil
				}
			}
		}
	} else if solution.childrenOfFirstRoom == 4 {
		for i := 0; i < len(way[0]); i++ {
			for k := 0; k < len(way[1]); k++ {
				for t := 0; t < len(way[3]); t++ {
					AnotherPath = append(AnotherPath, way[0][i], way[1][k], way[3][t])

					AnotherPath2 = append(AnotherPath2, AnotherPath)
					AnotherPath = nil
				}
			}
		}
	} else if solution.childrenOfFirstRoom == 1 {
		AnotherPath2 = append(AnotherPath2, way[0])
	}

	return AnotherPath2
}

func FindIntersect(way [][]Path) [][]Path {
	for i := 0; i < len(way); i++ {
		for k := 0; k < len(way[i])-1; k++ {
			for l := 0; l < len(way[i][k].paths)-1; l++ {
				for t := 0; t < len(way[i][k+1].paths)-1; t++ {
					if way[i][k].paths[l].name == way[i][k+1].paths[t].name {
						way[i][k].intersect = false
						way[i][k+1].intersect = false

					}
				}
			}
		}
	}
	return way
}

func (solution *solution) FindBestCombinations(way [][]Path) {
	for i := 0; i < len(way); i++ {
		intersectbool := false
		for k := 0; k < len(way[i]); k++ {
			for t := 0; t < len(way[i][k].paths); t++ {
			}
			if !way[i][k].intersect {
				intersectbool = true
			} else if k == len(way[i])-1 && !intersectbool {
				solution.BestCombinations = append(solution.BestCombinations, way[i])
			}
		}
	}
	if solution.BestCombinations == nil {
		solution.BestCombinations = append(solution.BestCombinations, way[0])
	}

	return
}

func (solution *solution) PathtoRoom(way [][]Path) {
	if len(way) > 1 {
		for i := 0; i < len(way)-1; i++ {
			countofpaths := 0
			for k := 0; k < len(way[i]); k++ {
				countofpaths += len(way[i][k].paths)
				way[i][k].totalLen = countofpaths
			}
		}
		for i := 0; i < len(way)-1; i++ {
			if way[i][len(way[i])-1].totalLen > way[i+1][len(way[i])-1].totalLen {
				way[i], way[i+1] = way[i+1], way[i]
			}
		}
	}
	for i := 0; i < len(way[0]); i++ {
		solution.BestPath = append(solution.BestPath, way[0][i].paths)
	}

	return
}

func SortBestPath(way [][]*room) [][]*room {
	for i := 0; i < len(way)-1; i++ {
		if len(way[i]) > len(way[i+1]) {
			way[i], way[i+1] = way[i+1], way[i]
		}
	}
	return way
}

func (solution *solution) EqNum(ants []ant, Path [][]*room) []ant {
	countofants := len(ants)
	for i := 0; i < len(Path)-1; i++ {
		if len(Path[i])+Path[i][0].queue > len(Path[i+1])+Path[i+1][0].queue {
			Path[i], Path[i+1] = Path[i+1], Path[i]
		}
	}
	ants[solution.counter].pathOfAnt = Path[0]
	Path[0][0].queue++
	solution.counter++
	if solution.counter != countofants {
		solution.EqNum(ants, Path)
	}
	return ants
}

func FindAnotherIntersect(way [][]*room) [][]*room {
	for i := 0; i < len(way)-1; i++ {
		for k := 0; k < len(way[i])-1; k++ {
			for t := 0; t < len(way[i+1])-1; t++ {
				if way[i][k].name == way[i+1][t].name {
					way = append(way[:i+1], way[i+2:]...)
					break
				}
			}
		}
	}
	return way
}
