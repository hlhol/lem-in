package maze

type Maze struct {
	rooms []*Room

	start *Room
	end   *Room
}

type Room struct {
	ants []*Ant

	name  string
	paths Paths
	nodes []*Node
}

type Ant struct {
	name string
}

type Paths []*Room

type Node struct {
	parentRoom *Room
	parentNode *Node

	room     *Room
	distance int
	siblings []*Node
}

type Solution struct {
	steps []Steps
}

type Steps []string
