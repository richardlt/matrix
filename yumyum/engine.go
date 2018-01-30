package yumyum

import "math/rand"

func newEngine(pc, w, h uint64) *engine {
	// prepare players
	ps := []*player{}
	if pc > 0 {
		ps = append(ps, newPlayer(coord{0, 0}))
	}
	if pc > 1 {
		ps = append(ps, newPlayer(coord{w - 1, h - 1}))
	}
	if pc > 2 {
		ps = append(ps, newPlayer(coord{w - 1, 0}))
	}
	if pc > 3 {
		ps = append(ps, newPlayer(coord{0, h - 1}))
	}

	// prepare candies
	cs := []*candy{}
	for i := uint64(0); i < pc*5; i++ {
	findGoodCoords:
		co := coord{uint64(rand.Intn(int(w))), uint64(rand.Intn(int(h)))}
		for _, c := range cs {
			if c.CheckIfOver(co) {
				goto findGoodCoords
			}
		}
		cs = append(cs, newCandy(co))
	}

	return &engine{ps, cs, h, w}
}

type engine struct {
	players               []*player
	candies               []*candy
	gridHeight, gridWidth uint64
}

func (e *engine) MovePlayer(playerSlot int, direction string) {
	if playerSlot < 0 || len(e.players) <= playerSlot {
		return
	}

	p := e.players[playerSlot]

	newCoord := p.Coord

	if direction == "left" && newCoord.X > 0 {
		newCoord.X--
	} else if direction == "up" && newCoord.Y > 0 {
		newCoord.Y--
	} else if direction == "right" && newCoord.X < e.gridWidth-1 {
		newCoord.X++
	} else if direction == "down" && newCoord.Y < e.gridHeight-1 {
		newCoord.Y++
	}

	for i, op := range e.players {
		if i != playerSlot && op.Coord.Equals(newCoord) {
			return
		}
	}

	p.Coord = newCoord

	for _, c := range e.candies {
		if c.State && c.CheckIfOver(p.Coord) {
			c.State = false
			p.Score += c.Points
			return
		}
	}
}

func (e engine) GetPlayers() []player {
	ps := make([]player, len(e.players))
	for i, p := range e.players {
		ps[i] = *p
	}
	return ps
}

func (e engine) GetCandies() []candy {
	cs := make([]candy, len(e.candies))
	for i, c := range e.candies {
		cs[i] = *c
	}
	return cs
}
