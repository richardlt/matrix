package zigzag

import (
	"math/rand"
)

func newEngine(sc, w, h uint64) *engine {
	// prepare snakes
	ss := []*snake{}
	if sc > 0 {
		ss = append(ss, newSnake(coord{2, 0}, "right", 3))
	}
	if sc > 1 {
		ss = append(ss, newSnake(coord{w - 3, h - 1}, "left", 3))
	}
	if sc > 2 {
		ss = append(ss, newSnake(coord{w - 1, 2}, "down", 3))
	}
	if sc > 3 {
		ss = append(ss, newSnake(coord{0, h - 3}, "up", 3))
	}

	// prepare candies
	cs := []*candy{}
	for i := uint64(0); i < sc*5; i++ {
	findGoodCoords:
		co := coord{uint64(rand.Intn(int(w))), uint64(rand.Intn(int(h)))}
		for _, c := range cs {
			if c.CheckIfOver(co) {
				goto findGoodCoords
			}
		}
		cs = append(cs, newCandy(co))
	}

	return &engine{snakes: ss, candies: cs, gridHeight: h, gridWidth: w}
}

type engine struct {
	snakes                []*snake
	candies               []*candy
	gridHeight, gridWidth uint64
}

func (e *engine) ChangePlayerDirection(playerSlot int, direction string) {
	if playerSlot < 0 || len(e.snakes) <= playerSlot {
		return
	}

	e.snakes[playerSlot].Direction = direction
}

func (e *engine) MovePlayers() {
	for playerSlot, s := range e.snakes {
		e.movePlayer(playerSlot, s)
	}
}

func (e *engine) IsGameOver() bool {
	// check first that there is more than one snake alive
	if len(e.snakes) > 1 {
		countAliveSnakes := 0
		for _, s := range e.snakes {
			if s.Length() > 0 {
				countAliveSnakes++
			}
		}
		if countAliveSnakes == 1 {
			return true
		}
	}

	// then check if there are still candies
	for _, c := range e.candies {
		if c.State {
			return false
		}
	}
	return true
}

func (e *engine) GetWinners() []int {
	max := e.snakes[0].Length()
	winners := []int{0}

	for i := 1; i < len(e.snakes); i++ {
		l := e.snakes[i].Length()
		if l > max {
			max = l
			winners = []int{i}
		} else if l == max {
			winners = append(winners, i)
		}
	}

	return winners
}

func (e *engine) movePlayer(playerSlot int, s *snake) {
	if s.Length() > 0 {
		s.Body = append(
			[]coord{s.Body[0].GetNear(s.Direction, e.gridWidth-1, e.gridHeight-1)},
			s.Body[:len(s.Body)-1]...,
		)

		for i, os := range e.snakes {
			if i != playerSlot && os.Length() > 0 && s.CheckIfOverAnotherSnake(*os) {
				s.LooseHead()
				return
			}
		}

		if s.CheckIfOverItself() {
			s.LooseHead()
			return
		}

		for _, c := range e.candies {
			if c.State && c.CheckIfOver(s.Body[0]) {
				c.State = false
				for i := 0; i < c.Points+1; i++ {
					s.GrowUp()
				}
				return
			}
		}
	}
}

func (e engine) GetSnakes() []snake {
	ss := make([]snake, len(e.snakes))
	for i, s := range e.snakes {
		ss[i] = *s
	}
	return ss
}

func (e engine) GetCandies() []candy {
	cs := make([]candy, len(e.candies))
	for i, c := range e.candies {
		cs[i] = *c
	}
	return cs
}
