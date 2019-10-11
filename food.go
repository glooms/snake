package main

import (
  "math/rand"
)

func newFood(w, h int, not []point) *point {
  if w < 2 || h < 2 {
    return nil
  }
  x := rand.Intn(w - 2) - w/2
  y := rand.Intn(h - 2) * 2 - h
  var p point
  for {
    found := false
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
