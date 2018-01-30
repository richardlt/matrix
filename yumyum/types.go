package yumyum

import (
	"math/rand"

	"github.com/richardlt/matrix/sdk-go/common"
)

type coord struct{ X, Y uint64 }

func (c coord) Equals(o coord) bool { return c.X == o.X && c.Y == o.Y }

func (c coord) Convert() common.Coord { return common.Coord{X: int64(c.X), Y: int64(c.Y)} }

func newCandy(c coord) *candy {
	return &candy{
		Coord:  c,
		Points: rand.Intn(2),
		State:  true,
	}
}

type candy struct {
	Coord  coord
	Points int
	State  bool
}

func (c candy) CheckIfOver(o coord) bool { return c.Coord.Equals(o) }

func newPlayer(c coord) *player { return &player{Coord: c} }

type player struct {
	Coord coord
	Score int
}
