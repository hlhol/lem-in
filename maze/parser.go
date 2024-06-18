package maze

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Load method to load the maze from input
func Load(inputStream *bufio.Reader) *Maze {
	maze := &Maze{}

	var antCount int

	// Read ant count
	line, err := inputStream.ReadString('\n')
	if err != nil || len(line) == 0 {
		return nil
	}

	antCount, err = strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return nil
	}

	if antCount <= 0 {
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
			// Read the next line to get start room name
			nextLine, err := inputStream.ReadString('\n')
			if err != nil {
				log.Println("Error reading next line after ##start:", err)
				return nil
			}
			parts := strings.Fields(nextLine)
			if len(parts) <= 0 {
				log.Println("Unexpected format for next line after ##start")
				return nil
			}

			newRoom := loadRooms(nextLine)
			if newRoom == nil {
				return nil
			}
			newRoom.ants = generateAnts(antCount)
			maze.rooms = append(maze.rooms, newRoom)

			maze.start = newRoom
			log.Println("Start room:", nextLine)
			// Ensure the start room is loaded as a regular room

			continue

		} else if strings.HasPrefix(line, "##end") {
			// Read the next line to get end room name
			nextLine, err := inputStream.ReadString('\n')
			if err != nil {
				log.Println("Error reading next line after ##end:", err)
				return nil
			}

			parts := strings.Fields(nextLine)
			if len(parts) <= 0 {
				log.Println("Unexpected format for next line after ##end")
				return nil
			}

			newRoom := loadRooms(nextLine)
			if newRoom == nil {
				return nil
			}
			maze.rooms = append(maze.rooms, newRoom)

			maze.end = newRoom
			log.Println("End room:", nextLine)
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

			newRoom := loadRooms(line)
			if newRoom == nil {
				return nil
			}
			maze.rooms = append(maze.rooms, newRoom)
		}
	}

	return maze
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
			room1 = maze.rooms[i]
		}
		if maze.rooms[i].name == room2Name {
			room2 = maze.rooms[i]
		}
	}

	// Handle if rooms are not found
	if room1 == nil || room2 == nil {
		log.Println("Room not found for connection:", line)
		return false
	}

	// Establish connection between rooms
	room1.paths = append(room1.paths, room2)
	room2.paths = append(room2.paths, room1)
	return true
}

func loadRooms(line string) *Room {
	// Parse room data
	parts := strings.Fields(line)
	if len(parts) != 3 {
		log.Println("Unexpected format for room line:", line)
		return nil
	}

	roomName := parts[0]
	_, err := strconv.Atoi(parts[1]) // x-coordinate
	if err != nil {
		log.Println("Error parsing x-coordinate:", err)
		return nil
	}

	_, err = strconv.Atoi(parts[2]) // y-coordinate
	if err != nil {
		log.Println("Error parsing y-coordinate:", err)
		return nil
	}

	// Create room
	room := &Room{
		name: roomName,
	}

	return room
}

func generateAnts(count int) []*Ant {
	ants := make([]*Ant, count)

	for i := range ants {
		name := strconv.Itoa(i + 1)

		ant := &Ant{name: name}
		ants[i] = ant
	}

	return ants
}
