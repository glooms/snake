package snake

import (
  "fmt"
  "testing"
)

func TestHead(t *testing.T) {
  m := newModel(2)
  p := point{0, 0}
  if h := m.head(); p != h {
    t.Errorf("Expected %s, was %s", p, h)
  }
}

func TestTurn(t *testing.T) {
  m := newModel(1)
  if !m.turn(DOWN) {
    t.Errorf("Could not turn from %s to %s", m.d, DOWN)
  }
  if m.d != DOWN {
    t.Errorf("Expected dir: %s, was: %s", DOWN, m.d)
  }
  if !m.turn(LEFT) {
    t.Errorf("Could not turn from %s to %s", m.d, LEFT)
  }
  if !m.turn(UP) {
    t.Errorf("Could not turn from %s to %s", m.d, UP)
  }
  if !m.turn(RIGHT) {
    t.Errorf("Could not turn from %s to %s", m.d, RIGHT)
  }
  if m.d != RIGHT {
    t.Errorf("Expected dir: %s, was: %s", RIGHT, m.d)
  }
}

func TestMove1(t *testing.T) {
  m := newModel(2)
  m.move()
  m.move()
  pts := m.get()
  if l := len(pts); l != 2 {
    t.Errorf("Expected length: %d, was instead %d", m.length, l)
  }
  p := point{2, 0}
  if pts[1] != p {
    t.Errorf("Expected %s, was %s", pts[1].String(), p.String())
  }
  fmt.Println()
}

func TestMove2(t *testing.T) {
  m := newModel(10)
  for i := 0; i < 5; i++ {
    m.move()
  }
  if l := len(m.get()); l != 6 {
    t.Errorf("Expected length: %d, was instead %d", 6, l)
  }
  m.turn(DOWN)
  for i := 0; i < 5; i++ {
    m.move()
  }
  if l := len(m.get()); l != m.length {
    t.Errorf("Expected length: %d, was instead %d", m.length, l)
  }
  p := point{5, 5}
  if h := m.head(); h != p {
    t.Errorf("Expected %s, was %s", p, h)
  }
  fmt.Println()
}

func TestMove3(t *testing.T) {
  m := newModel(10)
  steps := 150
  for i := 0; i < steps; i++ {
    m.move()
  }
  pts := m.get()
  if l := len(pts); l != m.length {
    t.Errorf("Expected length: %d, was instead %d", l, m.length)
  }
  if pts[0].x != steps - 10 + 1 {
    t.Errorf("Expected x: %d, was: %d", steps - 10, pts[0].x)
  }
  p := point{150, 0}
  if h := m.head(); p != h {
    t.Errorf("Expected %s, was %s", p, h)
  }
  fmt.Println()
  fmt.Println(m.body)
}

func TestGrow1(t *testing.T) {
  l := 10
  m := &model{
    body: make([]point, 100),
    length: l,
  }
  for i := 0; i < 20; i++ {
    m.move()
  }
  m.grow()
  if m.length != l + 1 {
    t.Errorf("Expected length: %d, was instead %d", l + 1, m.length)
  }
  if l := len(m.get()); l != m.length {
    t.Errorf("Expected length: %d, was instead %d", m.length, l)
  }
}
