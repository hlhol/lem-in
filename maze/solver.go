package maze

import (
	"fmt"
	"log"
	"sort"
)

func (maze *Maze) Solve() *Solution {
	log.Println("Solving the maze...")

	for _, v := range maze.rooms {
		fmt.Print("room ", v.name, " have connection:")
		for _, j := range v.paths {
			fmt.Print("  ", j.name)
		}
		fmt.Println()
	}

	// Explore all paths from start-to-end room using BFS
	reachedEnd := maze.start.CreateNodes(maze.start, maze.end, nil)

	if !reachedEnd {
		log.Println("No paths found from start room to end room.")
		return nil
	}

	maze.sort()

	solution := &Solution{}

	for i := 0; i < 10000; i++ {
		steps := Steps{}
		visitedRoom := Paths{}

		//maze.printRooms()
		for r := len(maze.rooms) - 1; r > -1; r-- {
			room := maze.rooms[r]

			if room == maze.end || room.antCount == 0 || visitedRoom.dejavu(room) {
				continue
			}

			var possibleNodes []Node

			//fmt.Print("visiting rooms: ")
			for _, node := range room.nodes {
				//fmt.Println(node.room.name, node.distance)
				if room.antCount == 0 || (node.room != maze.end && node.room.isOccupied()) {
					continue
				}

				possibleNodes = append(possibleNodes, node)
			}

			if len(possibleNodes) == 0 {
				continue
			}

			sort.Slice(possibleNodes, func(i, j int) bool {
				return possibleNodes[i].distance > possibleNodes[j].distance
			})
			node := possibleNodes[0]

			currentDistance := room.nodes[0].distance
			nextDistance := node.distance

			if currentDistance > nextDistance {
				continue
			}
			r = len(maze.rooms) - 1

			room.antCount--
			node.room.antCount++

			visitedRoom = append(visitedRoom, node.room)

			step := fmt.Sprintf("L%v-%s", i+1, node.room.name)
			steps.Append(step)

		}

		if len(steps) == 0 {
			break
		}

		solution.AddSteps(steps)
	}

	return solution
}

func (maze *Maze) sort() {
	for _, v := range maze.rooms {
		if v.nodes == nil {
			v.nodes = []Node{}
		}

		sort.Slice(v.nodes, func(i, j int) bool {
			return v.nodes[i].distance < v.nodes[j].distance
		})
	}

	sort.Slice(maze.rooms, func(i, j int) bool {

		roomI := maze.rooms[i]
		roomJ := maze.rooms[j]

		if len(roomI.nodes) == 0 {
			return false
		}

		if len(roomJ.nodes) == 0 {
			return true
		}

		return roomI.nodes[0].distance < roomJ.nodes[0].distance
	})
}

func (room *Room) CreateNodes(start, end *Room, parent *Node) bool {

	if room == end {
		return true
	}

	for _, path := range room.paths {

		if parent != nil && parent.haveRoom(room) || path == start {
			continue
		}

		node := Node{
			parentRoom: room,
			parentNode: parent,
			room:       path,
		}
		node.CalculateDistance()

		haveEnd := path.CreateNodes(start, end, &node)

		if !haveEnd {
			continue
		}

		room.nodes = append(room.nodes, node)
	}

	return len(room.nodes) > 0
}

func (room *Room) isOccupied() bool {
	return room.antCount > 0
}

func (node *Node) haveRoom(room *Room) bool {
	for node != nil {
		if node.parentRoom == room {
			return true
		}

		node = node.parentNode
	}
	return false
}

func (node *Node) CalculateDistance() {
	targetNode := node

	for node != nil {
		targetNode.distance++
		node = node.parentNode
	}
}

func (path *Paths) dejavu(room *Room) bool {
	for _, p := range *path {
		if p == room {
			return true
		}
	}

	return false
}

func (maze *Maze) printRooms() {
	for _, v := range maze.rooms {

		var distance int

		if v.nodes == nil || len(v.nodes) == 0 {
			distance = 0
		} else {
			distance = v.nodes[0].distance
		}

		fmt.Println("room: ", v.name, ", count: ", v.antCount, "distance", distance)
	}
}
