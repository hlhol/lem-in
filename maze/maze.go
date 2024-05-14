package maze

import (
	"bufio"
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
	distance    int
	name        string
	occupancies int
	paths       []string
}

var collection [][]string

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
		name:        roomName,
		occupancies: 0,
		distance:    0,
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
	} else {
		log.Println("Paths found:")
		for i, path := range paths {
			log.Printf("Path %d: %s", i+1, path)
		}
	}
}

func (maze *Maze) explorePaths(currentRoom Room, endRoom string, currentPath []string, paths *[]string) {
	currentPath = append(currentPath, currentRoom.name) //  is the end room

	if currentRoom.name == endRoom { //  is the end room
		*paths = append(*paths, strings.Join(currentPath, " -> "))
		return
	}

	// Explore all possible paths from current room
	for _, nextRoomName := range currentRoom.paths {
		for _, room := range maze.rooms {
			if room.name == nextRoomName {
				// Explore next room if it's not already visited
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
