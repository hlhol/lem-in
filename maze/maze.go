package maze

import (
	"bufio"
	"fmt"
	"log"
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

	for _, room := range maze.rooms { // Iterate through rooms to find paths
		if room.name == maze.startRoom {
			maze.explorePaths(room, maze.endRoom, []string{}, &paths)
		}
	}

	if len(paths) == 0 {
		log.Println("No paths found from start room to end room.")
		return
	}

	paths = filteries(paths)
	log.Println("Paths found:")
	for i, path := range paths {
		log.Printf("Path %v: %s", i+1, path)
	}

	antsPerPath := make([]int, len(paths))

	for ants := maze.antCount; ants > 0; ants-- {
		minPathIndex := findShortestPathIndex(paths, antsPerPath)

		antsPerPath[minPathIndex]++
	}

	antPositions := make([][]string, maze.antCount)
	antIndex := 0
	for i, ants := range antsPerPath {
		pathRooms := strings.Split(paths[i], " -> ")
		for j := 0; j < ants; j++ {
			antPositions[antIndex] = pathRooms
			antIndex++
		}
	}

	occupiedRooms := make(map[string]bool)

	for i := 0; i < maze.antCount; i++ {
		occupiedRooms[maze.startRoom] = true
	}

	steps := 0
	for {
		moved := false
		stepOutput := []string{}
		newOccupiedRooms := make(map[string]bool)

		for i, pos := range antPositions {
			if len(pos) > 1 {
				if !occupiedRooms[pos[1]] || pos[1] == maze.endRoom {
					stepOutput = append(stepOutput, fmt.Sprintf("L%v-%s", i+1, pos[1]))
					antPositions[i] = pos[1:]
					newOccupiedRooms[pos[1]] = true
					moved = true
				} else {
					newOccupiedRooms[pos[0]] = true
				}
			} else {
				stepOutput = append(stepOutput, fmt.Sprintf("L%v-%s", i+1, pos[0]))
			}
		}

		if !moved {
			break
		}

		occupiedRooms = newOccupiedRooms
		steps++
		fmt.Printf("Step %v: %s\n", steps, strings.Join(stepOutput, " "))
	}
}

func findShortestPathIndex(paths []string, antsPerPath []int) int {
	minIndex := 0
	minSteps := calculateSteps(paths[0], antsPerPath[0])

	for i := 1; i < len(paths); i++ {
		steps := calculateSteps(paths[i], antsPerPath[i])
		if steps < minSteps {
			minSteps = steps
			minIndex = i
		}
	}

	return minIndex
}

func calculateSteps(path string, ants int) int {
	rooms := roompath(path)
	return rooms + ants
}

func roompath(path string) int {
	parts := strings.Split(path, " -> ")
	return len(parts) + 1
}

func (maze *Maze) explorePaths(currentRoom Room, endRoom string, currentPath []string, paths *[]string) {
	currentPath = append(currentPath, currentRoom.name)

	if currentRoom.name == endRoom { //  is the end room
		*paths = append(*paths, strings.Join(currentPath, " -> "))
		return
	}

	for _, nextRoomName := range currentRoom.paths {
		for _, room := range maze.rooms {
			if room.name == nextRoomName {
				if !contains(currentPath, room.name) {
					maze.explorePaths(room, endRoom, currentPath, paths)
				}
			}
		}
	}
}

func contains(path []string, room string) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}


func filteries(paths []string) []string{
	filtered := make([]string, 0, len(paths))

	for _,p1 := range paths {
		include := true
		for _,p2 := range paths {
			rooms := strings.Split(p2, " ")
			
			if len(p1) > len(p2) {
				for _,room := range rooms {
					if strings.Contains(p1 ,room) {
						include = false
						break
					}
				}
			}

		}
		if include {
			filtered = append(filtered, p1)
		}
	}

	return filtered
}

