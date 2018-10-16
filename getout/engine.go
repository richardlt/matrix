package getout

import (
	"math/rand"
	"time"
)

type coord struct{ x, y int }

func randCoord(r *rand.Rand, w, h int) coord { return coord{r.Intn(w-2) + 1, r.Intn(h-2) + 1} }

func newEngine(w, h int) *engine {
	e := engine{
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
		gridWidth: w, gridHeight: h,
		grid: make([][]int, w),
	}
	for i := 0; i < w; i++ {
		e.grid[i] = make([]int, h)
	}

	e.player = randCoord(e.r, w, h)
	e.grid = e.explore(e.grid, e.player, 1)

	return &e
}

type engine struct {
	r                     *rand.Rand
	gridWidth, gridHeight int
	grid                  [][]int
	player, end           coord
	depth                 int
}

func (e *engine) isGameOver() bool { return e.player == e.end }

func (e *engine) explore(grid [][]int, origin coord, depth int) [][]int {
	// set to 1 the current cell
	grid[origin.x][origin.y] = 1

	// while there is available cells around the origin
	for {
		var available []coord

		// compute available cells around origin x-2/+2 and y-2/+2
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

		// randomly get the next cell from availables
		next := available[e.r.Intn(len(available))]

		// if exploration depth superior than max depth set the end to next coord
		if depth > e.depth {
			e.end = next
			e.depth = depth
		}

		// set to 1 the cell between origin and next
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

		// if next is an edge cell set to 2 else explore from next
		isEdge := 0 == next.x || next.x == len(grid)-1 || 0 == next.y || next.y == len(grid[next.x])-1
		if isEdge {
			grid[next.x][next.y] = 2
		} else {
			grid = e.explore(grid, next, depth+1)
		}
	}
}

func (e *engine) move(direction string) {
	newCoord := e.player

	if direction == "left" && newCoord.x > 0 {
		newCoord.x--
	} else if direction == "up" && newCoord.y > 0 {
		newCoord.y--
	} else if direction == "right" && newCoord.x < e.gridWidth-1 {
		newCoord.x++
	} else if direction == "down" && newCoord.y < e.gridHeight-1 {
		newCoord.y++
	}

	if e.grid[newCoord.x][newCoord.y] == 0 {
		return
	}

	e.player = newCoord
}
