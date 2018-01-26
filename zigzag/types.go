package zigzag

import (
	"math/rand"

	"github.com/richardlt/matrix/sdk-go/common"
)

type coord struct{ X, Y uint64 }

func (c coord) Equals(o coord) bool { return c.X == o.X && c.Y == o.Y }

func (c coord) GetNear(direction string, maxWidth, maxHeight uint64) coord {
	switch direction {
	case "left":
		if c.X > 0 {
			c.X--
		} else {
			c.X = maxWidth
		}
	case "up":
		if c.Y > 0 {
			c.Y--
		} else {
			c.Y = maxHeight
		}
	case "right":
		if c.X < maxWidth {
			c.X++
		} else {
			c.X = 0
		}
	case "down":
		if c.Y < maxHeight {
			c.Y++
		} else {
			c.Y = 0
		}
	}
	return c
}

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

func newSnake(head coord, direction string, length uint64) *snake {
	s := &snake{[]coord{}, direction}

	for i := uint64(0); i < length; i++ {
		if direction == "left" {
			s.Body = append(s.Body, coord{head.X + i, head.Y})
		} else if direction == "up" {
			s.Body = append(s.Body, coord{head.X, head.Y + i})
		} else if direction == "right" {
			s.Body = append(s.Body, coord{head.X - i, head.Y})
		} else if direction == "down" {
			s.Body = append(s.Body, coord{head.X, head.Y - i})
		}
	}

	return s
}

type snake struct {
	Body      []coord
	Direction string
}

func (s *snake) GrowUp() {
	if len(s.Body) > 1 {
		end := s.Body[len(s.Body)-1]
		beforeEnd := s.Body[len(s.Body)-2]
		if end.X == beforeEnd.X {
			if end.Y < beforeEnd.Y {
				s.Body = append(s.Body, coord{end.X, end.Y - 1})
			} else {
				s.Body = append(s.Body, coord{end.X, end.Y + 1})
			}
		} else {
			if end.X < beforeEnd.X {
				s.Body = append(s.Body, coord{end.X - 1, end.Y})
			} else {
				s.Body = append(s.Body, coord{end.X + 1, end.Y})
			}
		}
	} else {
		end := s.Body[len(s.Body)-1]
		if s.Direction == "left" {
			s.Body = append(s.Body, coord{(end.X + 1), end.Y})
		} else if s.Direction == "up" {
			s.Body = append(s.Body, coord{end.X, (end.Y - 1)})
		} else if s.Direction == "right" {
			s.Body = append(s.Body, coord{(end.X - 1), end.Y})
		} else if s.Direction == "down" {
			s.Body = append(s.Body, coord{end.X, (end.Y - 1)})
		}
	}
}

func (s snake) CheckIfOverAnotherSnake(o snake) bool {
	for i := 0; i < len(o.Body); i++ {
		if s.Body[0].Equals(o.Body[i]) {
			return true
		}
	}
	return false
}

func (s *snake) LooseHead() { s.Body = s.Body[1:] }

func (s snake) CheckIfOverItself() bool {
	for i := 1; i < len(s.Body); i++ {
		if s.Body[0].Equals(s.Body[i]) {
			return true
		}
	}
	return false
}

func (s snake) Length() int { return len(s.Body) }
