package snake

import (
	"fmt"
	"math/rand"
)

type model struct {
	cur, next dir
	body      []point
	length    int
	at        int
}

type point struct {
	x, y int
}

type dir struct {
	dx, dy int
}

var (
	NIL   = dir{0, 0}
	LEFT  = dir{-1, 0}
	DOWN  = dir{0, 1}
	UP    = dir{0, -1}
	RIGHT = dir{1, 0}
)

func newModel(length int) *model {
	capacity := 100
	if length > capacity/2 {
		capacity = length * 4
	}
	m := &model{
		cur:    RIGHT,
		body:   make([]point, capacity),
		length: length,
	}
	return m
}

func (m *model) move() {
	if m.at == len(m.body)-1 {
		for i, p := range m.get() {
			m.body[i] = p
		}
		m.at = m.length - 1
	}
	p := m.body[m.at]
	m.at++
	if m.next != NIL {
		m.cur = m.next
		m.next = NIL
	}
	m.body[m.at] = point{
		p.x + m.cur.dx,
		p.y + m.cur.dy,
	}
}

func (m *model) grow() {
	m.length++
	if m.length > len(m.body)/2 {
		nbody := make([]point, m.length*4)
		copy(nbody, m.body)
		m.body = nbody
	}
}

func (m *model) head() point {
	return m.body[m.at]
}

func (m *model) get() []point {
	from := m.at - m.length + 1
	if from < 0 {
		from = 0
	}
	to := m.at + 1
	return m.body[from:to]
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func (m *model) turn(d dir) bool {
	if m.cur.dx != d.dx && m.cur.dy != d.dy {
		m.next = d
		return true
	}
	return false
}

func newFood(w, h int, not []point) *point {
	if w < 2 || h < 2 {
		return nil
	}
	var p point
	for {
		found := false
		x := rand.Intn(w/4 - 1)
		y := rand.Intn(h/2 - 1)
		p = point{x, y}
		for _, n := range not {
			if p == n {
				found = true
				continue
			}
		}
		if !found {
			return &p
		}
	}
	return nil
}

func (d dir) String() string {
	switch d {
	case LEFT:
		return "LEFT"
	case DOWN:
		return "DOWN"
	case UP:
		return "UP"
	case RIGHT:
		return "RIGHT"
	}
	return ""
}
