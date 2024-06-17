package maze

import (
	"fmt"
	"log"
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
	reachedEnd := maze.start.CreateNodes(maze.end, nil)

	if !reachedEnd {
		log.Println("No paths found from start room to end room.")
		return nil
	}

	solution := &Solution{}

	for i := 0; i < 1000000; i++ {
		steps := Steps{}
		visitedRoom := Paths{}

		for _, room := range maze.rooms {
			if room == maze.end || room.antCount == 0 || visitedRoom.dejavu(room) {
				continue
			}

			fmt.Print("visiting rooms: ")
			for _, node := range room.nodes {
				if room.antCount == 0 || (node.room != maze.end && node.room.isOccupied()) {
					continue
				}
				// TODO select shortest path
				room.antCount--
				node.room.antCount++

				//visitedRoom = append(visitedRoom, room)
				visitedRoom = append(visitedRoom, node.room)
				fmt.Print(" ", node.room.name)

				step := fmt.Sprintf("L%v-%s", i+1, node.room.name)
				steps.Append(step)
			}
			fmt.Println()
		}

		if len(steps) == 0 {
			break
		}

		solution.AddSteps(steps)
	}

	return solution
}

func (room *Room) CreateNodes(end *Room, parent *Node) bool {

	if room == end {
		parent.CalculateDistance()
		return true
	}

	for _, path := range room.paths {

		if parent != nil && parent.haveRoom(room) {
			continue
		}

		node := Node{
			parentRoom: room,
			parentNode: parent,
			room:       path,
		}

		haveEnd := path.CreateNodes(end, &node)

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
