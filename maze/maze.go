package maze

type Maze struct {
	lines          []string // lines describing the maze (no comments or commands)
	duplicateRooms []string // rooms already added to Farm (field bellow)

	startName string // start room name
	endName   string

	Farm      []room
	startRoom *room
	endRoom   *room

	antCount int
	ants     []ant
}

type ant struct {
	id        int
	curRoom   *room
	pathOfAnt []*room
	step      int
}

type Path struct {
	id        int     // unique identifier for the path.
	paths     []*room // the rooms in the path.
	intersect bool    // whether the path intersects with another path.
	queue     int     // the number of ants waiting to use the path.
	totalLen  int     // the total length of the path.
}

type room struct {
	name     string
	parent   []*room
	children []*room
	occupied bool
	queue    int
}

type solution struct {
	startRoomChildren int
	countloop         int
	appendWays        []*room
	CombinatedRooms   [][]*room
	BestCombinations  [][]Path
	BestPath          [][]*room
	counter           int
}
