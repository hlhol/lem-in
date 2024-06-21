package maze

import (
	"io/ioutil"
	"lem/utils"
	"log"
	"strconv"
	"strings"
)

func (maze *Maze) Start() {

	if maze.antCount == 0 {
		utils.Fatal("ERROR\nAnt Number 0 - Check The Ant Number")
	}

	// Create the ant farm based on the rooms and their connections
	maze.Load()

	// Swap rooms order
	maze.SwapRooms()

	// Add the last room to the ant farm
	maze.Farm = append(maze.Farm, *maze.endRoom)

	// Create the ants
	maze.CreatingAnts()

	// Find all possible paths between the rooms in the ant farm
	var allPaths [][]*room
	maze.FindAllPossiblePaths([]*room{}, maze.Farm[0], &allPaths)

	// Sort the paths based on their length
	for i := 0; i < len(allPaths); i++ {
		allPaths = SortPaths(allPaths)
	}

	s := solution{
		startRoomChildren: len(maze.startRoom.children),
	}

	// Clean the paths and remove any duplicate paths
	s.ClearPath(maze, allPaths)

	// Finding the best paths for first children of the first room
	var paths []Path
	var pathsCombinations [][]Path
	paths = s.firstRoomPaths(maze, allPaths)

	for i := 0; i < s.startRoomChildren; i++ {
		pathsCombinations = append(pathsCombinations, SortedPaths(paths, i))
	}

	pathsCombinations = SortAgain(pathsCombinations)
	pathsCombinations = s.AllCombinations(pathsCombinations)

	// Finding intersecting paths
	FindIntersect(pathsCombinations)

	// Finding the best combinations of paths
	s.FindBestCombinations(pathsCombinations)

	// Converting the best paths to rooms
	s.PathtoRoom(s.BestCombinations)
	s.BestPath = SortBestPath(s.BestPath)

	// Finding intersecting paths in the best path
	// intersect help reduce the time ant will wait and can provide a shorter path
	s.BestPath = FindAnotherIntersect(s.BestPath)
	s.BestPath = FindAnotherIntersect(s.BestPath)

	// Making sure that the number of ants in the maze is equal to the number of rooms in the best path
	maze.ants = s.EqNum(maze.ants, s.BestPath)

	// Making ants walk through the farm
	maze.walk()
}

func ReadFile(textfile string) *Maze {
	content, err := ioutil.ReadFile(textfile)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) == 0 {
		utils.Fatal("Empty input file")
	}

	antCount, _ := strconv.Atoi(lines[0])
	lines = lines[1:]

	rooms := make([]string, len(lines))
	validLines := 0

	startLine := -1
	endLine := -1

	for i, v := range lines {
		if len(v) == 0 {
			continue
		}

		if v[0] == '#' {
			command := testCommand(v)

			switch command {
			case 0:
				continue

			case 1:
				if startLine != -1 {
					utils.Fatal("multiple ##startName found")
				}
				startLine = validLines
				continue

			case 2:
				if endLine != -1 {
					utils.Fatal("multiple ##end found")
				}
				endLine = validLines
				continue

			default:
				utils.Fatal("base line at %d", i+2) // index startName from 0 + first line is skipped
			}
		}

		rooms[validLines] = v
		validLines++
	}

	rooms = rooms[:validLines]

	startName := strings.Split(rooms[startLine], " ")[0]
	endName := strings.Split(rooms[endLine], " ")[0]

	maze := &Maze{
		startName: startName,
		endName:   endName,
		lines:     rooms,
		antCount:  antCount,
	}

	return maze
}

// check if the given line is a comment or and label
// return 0 for comment, 1 for startName, 2 for end
// if bad syntax is used return -1
func testCommand(line string) int {
	if len(line) < 2 {
		return 0
	}

	if line[:2] == "##" {
		label := line[2:]

		switch label {
		case "start":
			return 1

		case "end":
			return 2

		default:
			return -1
		}
	}

	return 0
}
