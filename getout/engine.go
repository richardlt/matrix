package getout

import (
	"math/rand"
	"time"
)

type coord struct{ x, y int }

func randCoord(r *rand.Rand, w, h int) coord { return coord{r.Intn(w-2) + 1, r.Intn(h-2) + 1} }

func newEngine(w, h int) *engine {
	grid := make([][]int, w)
	for i := 0; i < w; i++ {
		grid[i] = make([]int, h)
	}

	e := engine{
		r:          rand.New(rand.NewSource(time.Now().UnixNano())),
		gridHeight: h, gridWidth: w,
		grid: grid,
	}

	e.start = randCoord(e.r, w, h)
	e.grid = e.explore(grid, e.start, 1)

	return &e
}

func (e *engine) explore(grid [][]int, origin coord, depth int) [][]int {
	grid[origin.x][origin.y] = 1

	for {
		var available []coord

		for _, x := range []int{origin.x - 2, origin.x + 2} {
			if 0 <= x && x < len(grid) {
				cell := grid[x][origin.y]
				if cell == 0 {
					available = append(available, coord{x, origin.y})
				}
			}
		}

		column := grid[origin.x]
		for _, y := range []int{origin.y - 2, origin.y + 2} {
			if 0 <= y && y < len(column) {
				cell := column[y]
				if cell == 0 {
					available = append(available, coord{origin.x, y})
				}
			}
		}

		if len(available) == 0 {
			return grid
		}

		next := available[e.r.Intn(len(available))]

		if depth > e.depth {
			e.end = next
			e.depth = depth
		}

		moveX := next.x - origin.x
		if 0 < moveX {
			grid[next.x-1][next.y] = 1
		} else if moveX < 0 {
			grid[next.x+1][next.y] = 1
		}

		moveY := next.y - origin.y
		if 0 < moveY {
			grid[next.x][next.y-1] = 1
		} else if moveY < 0 {
			grid[next.x][next.y+1] = 1
		}

		isEdge := 0 == next.x || next.x == len(grid)-1 || 0 == next.y || next.y == len(grid[next.x])-1
		if isEdge {
			grid[next.x][next.y] = 2
		} else {
			grid = e.explore(grid, next, depth+1)
		}
	}
}

type engine struct {
	r                     *rand.Rand
	gridHeight, gridWidth int
	grid                  [][]int
	start, end            coord
	depth                 int
}
