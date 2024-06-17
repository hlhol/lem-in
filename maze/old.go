package maze

/*
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

package maze

import (
"bufio"
"fmt"
"log"
"os"
"strconv"
"strings"
)

type Maze struct {
	antCount  int
	rooms     []Room
	startRoom string
	endRoom   string
}

type Room struct {
	name  string
	paths []string
}

// Load method to load the maze from input
func Load(inputStream *bufio.Reader) *Maze {
	maze := &Maze{}

	// Read ant count
	line, err := inputStream.ReadString('\n')
	if err != nil || len(line) == 0 {
		return nil
	}
	maze.antCount, err = strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return nil
	}

	if maze.antCount <= 0 {
		fmt.Printf("bad input the ant should be more than 0\n")
		os.Exit(0)
	}

	// Load rooms and connections
	for {
		line, err := inputStream.ReadString('\n')
		if err != nil || len(line) == 0 {
			break
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "##start") {
			// Read next line to get start room name
			nextLine, err := inputStream.ReadString('\n')
			if err != nil {
				log.Println("Error reading next line after ##start:", err)
				return nil
			}
			parts := strings.Fields(nextLine)
			if len(parts) > 0 {
				maze.startRoom = parts[0]
				log.Println("Start room:", maze.startRoom)
				// Ensure the start room is loaded as a regular room
				if !maze.loadRooms(nextLine) {
					return nil
				}
			} else {
				log.Println("Unexpected format for next line after ##start")
				return nil
			}
			continue
		} else if strings.HasPrefix(line, "##end") {
			// Read next line to get end room name
			nextLine, err := inputStream.ReadString('\n')
			if err != nil {
				log.Println("Error reading next line after ##end:", err)
				return nil
			}
			parts := strings.Fields(nextLine)
			if len(parts) > 0 {
				maze.endRoom = parts[0]
				log.Println("End room:", maze.endRoom)
				if !maze.loadRooms(nextLine) {
					return nil
				}
			} else {
				log.Println("Unexpected format for next line after ##end")
				return nil
			}
			continue
		}

		if strings.Contains(line, "-") {
			// Connection line
			log.Println("Mapping connection:", line)
			if !maze.mapFarm(line) {
				return nil
			}
		} else {
			// Room line
			log.Println("Loading room:", line)
			if !maze.loadRooms(line) {
				return nil
			}
		}
	}

	return maze
}

func (maze *Maze) loadRooms(line string) bool {
	// Parse room data
	parts := strings.Fields(line)
	if len(parts) != 3 {
		log.Println("Unexpected format for room line:", line)
		return false
	}
	roomName := parts[0]
	_, err := strconv.Atoi(parts[1]) // x-coordinate
	if err != nil {
		log.Println("Error parsing x-coordinate:", err)
		return false
	}

	_, err = strconv.Atoi(parts[2]) // y-coordinate
	if err != nil {
		log.Println("Error parsing y-coordinate:", err)
		return false
	}

	// Create room
	room := Room{
		name: roomName,
	}

	maze.rooms = append(maze.rooms, room)
	return true
}

func (maze *Maze) mapFarm(line string) bool {
	// Parse connection data
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		log.Println("Unexpected format for connection line:", line)
		return false
	}
	room1Name := parts[0]
	room2Name := parts[1]

	log.Printf("Mapping connection between rooms: %s and %s\n", room1Name, room2Name)

	// Find room objects by name
	var room1, room2 *Room
	for i := range maze.rooms {
		if maze.rooms[i].name == room1Name {
			room1 = &maze.rooms[i]
		}
		if maze.rooms[i].name == room2Name {
			room2 = &maze.rooms[i]
		}
	}

	// Handle if rooms are not found
	if room1 == nil || room2 == nil {
		log.Println("Room not found for connection:", line)
		return false
	}

	// Establish connection between rooms
	room1.paths = append(room1.paths, room2Name)
	room2.paths = append(room2.paths, room1Name)
	return true
}

func (maze *Maze) Solve() {
	log.Println("Solving the maze...")

	var paths []string

	// Explore all paths from start to end room using BFS
	maze.mapNodes(&paths)

	if len(paths) == 0 {
		log.Println("No paths found from start room to end room.")
		return
	}

	// Filter paths using the custom filtries function
	paths2 := filtries(paths, maze, maze.antCount)
	if len(paths2) == 0 {
		log.Println("No valid paths after filtering.")
		return
	}

	log.Println("Paths found:")
	for i, path := range paths2 {
		log.Printf("Path %v: %s", i+1, path)
	}

	// Distribute ants along the paths
	antsPerPath := distributeAnts(paths2, maze.antCount)

	antPositions := make([][]string, maze.antCount)
	antIndex := 0
	for i, ants := range antsPerPath {
		pathRooms := strings.Split(paths2[i], " -> ")
		for j := 0; j < ants; j++ {
			antPositions[antIndex] = append([]string{}, pathRooms...)
			antIndex++
		}
	}

	steps := 0
	for {
		moved := false
		stepOutput := []string{}

		for i := range antPositions {
			pos := antPositions[i]
			if len(pos) > 1 {
				nextRoom := pos[1]

				if nextRoom == maze.endRoom || !isRoomOccupied(antPositions, nextRoom, i) {
					stepOutput = append(stepOutput, fmt.Sprintf("L%v-%s", i+1, nextRoom))
					antPositions[i] = pos[1:]
					moved = true
				}
			}
		}

		if !moved {
			break
		}

		steps++
		fmt.Printf("Step %v: %s\n", steps, strings.Join(stepOutput, " "))
	}

	for _, v := range paths {
		fmt.Println(v)
	}
}

func distributeAnts(paths []string, antCount int) []int {
	antsPerPath := make([]int, len(paths))
	rooms := make([]int, len(paths))
	totalSteps := make([]int, len(paths))

	for i, path := range paths {
		rooms[i] = roomCount(path)
		totalSteps[i] = rooms[i]
	}

	remainingAnts := antCount
	for remainingAnts > 0 {
		minStepsIndex := findMinSteps(totalSteps)
		antsPerPath[minStepsIndex]++
		totalSteps[minStepsIndex] = rooms[minStepsIndex] + antsPerPath[minStepsIndex]
		remainingAnts--
	}

	return antsPerPath
}

func findMinSteps(totalSteps []int) int {
	minSteps := totalSteps[0]
	minIndex := 0

	for i, steps := range totalSteps {
		if steps < minSteps {
			minSteps = steps
			minIndex = i
		}
	}

	return minIndex
}

func roomCount(path string) int {
	parts := strings.Split(path, " -> ")
	return len(parts)
}

func (maze *Maze) mapNodes(paths *[]string) {
	type PathNode struct {
		currentRoom string
		path        []string
	}

	queue := []PathNode{{currentRoom: maze.startRoom, path: []string{maze.startRoom}}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.currentRoom == maze.endRoom {
			*paths = append(*paths, strings.Join(node.path, " -> "))
			continue
		}

		for _, nextRoomName := range maze.getPaths(node.currentRoom) {
			if !contains(node.path, nextRoomName) {
				newPath := append([]string{}, node.path...)
				newPath = append(newPath, nextRoomName)
				queue = append(queue, PathNode{currentRoom: nextRoomName, path: newPath})
				visited[nextRoomName] = true
			}
		}
	}
}

func (maze *Maze) getPaths(roomName string) []string {
	for _, room := range maze.rooms {
		if room.name == roomName {
			return room.paths
		}
	}
	return []string{}
}

func contains(path []string, room string) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}

// func filtries(paths []string, maze *Maze, antCount int) []string {
// 	if len(paths) == 0 {
// 		log.Println("No paths to filter")
// 		return paths
// 	}

// 	var finalPaths []string
// 	minSteps := -1

// 	subsets := gensub(paths)

// 	for _, subset := range subsets {
// 		if len(subset) == 0 { // Skip the empty subset
// 			continue
// 		}

// 		steps := calculateStepsForSubset(subset, maze, antCount)
// 		if minSteps == -1 || steps < minSteps {
// 			minSteps = steps
// 			finalPaths = append([]string{}, subset...)
// 		}
// 	}

// 	return finalPaths
// }

func filtries(paths []string, maze *Maze, antCount int) []string {
	if len(paths) == 0 {
		log.Println("No paths to filter")
		return paths
	}

	/*for _,p := range paths{
		for _,p1 := range paths {

		}
	}star/

	var finalPaths []string
	minSteps := -1

	for subsetSize := 1; subsetSize <= len(paths); subsetSize++ {
		currentSubset := paths[:subsetSize]

		steps := calculateStepsForSubset(currentSubset, maze, antCount)
		if minSteps == -1 || steps < minSteps {
			minSteps = steps
			finalPaths = append([]string{}, currentSubset...)
		}
	}

	return finalPaths
}

func gensub(paths []string) [][]string {
	allSubs := [][]string{{}}
	for _, p := range paths {
		var newSubsets [][]string
		for _, subset := range allSubs {
			newSubset := append([]string{}, subset...)
			newSubset = append(newSubset, p)
			newSubsets = append(newSubsets, newSubset)
		}
		allSubs = append(allSubs, newSubsets...)
	}
	return allSubs
}

func calculateStepsForSubset(paths []string, maze *Maze, antCount int) int {
	antsPerPath := distributeAnts(paths, antCount)
	antPositions := make([][]string, maze.antCount)
	antIndex := 0

	for i, ants := range antsPerPath {
		pathRooms := strings.Split(paths[i], " -> ")
		for j := 0; j < ants; j++ {
			antPositions[antIndex] = append([]string{}, pathRooms...)
			antIndex++
		}
	}

	steps := 0
	for {
		moved := false
		for i := range antPositions {
			pos := antPositions[i]
			if len(pos) > 1 {
				nextRoom := pos[1]

				if nextRoom == maze.endRoom || !isRoomOccupied(antPositions, nextRoom, i) {
					antPositions[i] = pos[1:]
					moved = true
				}
			}
		}
		if !moved {
			break
		}
		steps++
	}

	return steps
}

func isRoomOccupied(positions [][]string, room string, currentAnt int) bool {
	for i, pos := range positions {
		if i != currentAnt && len(pos) > 0 && pos[0] == room {
			return true
		}
	}
	return false
}
*/
