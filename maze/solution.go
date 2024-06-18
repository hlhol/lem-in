package maze

import (
	"fmt"
	"strings"
)

func (steps *Steps) Append(step string) {
	*steps = append(*steps, step)
}

func (solution *Solution) AddSteps(steps Steps) {
	//line := strings.Join(steps, " ")
	//fmt.Printf("Step %v: %s\n", 1, line)
	solution.steps = append(solution.steps, steps)
}

func (solution *Solution) Print() {
	if solution == nil {
		fmt.Println("Can't Solve Maze")
		return
	}

	for index, steps := range solution.steps {
		line := strings.Join(steps, " ")
		fmt.Printf("Step %v: %s\n", index, line)
	}
}
