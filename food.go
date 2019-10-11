package main

import (
  "math/rand"
)

func newFood(w, h int, not []point) *point {
  if w < 2 || h < 2 {
    return nil
  }
  var p point
  for {
    found := false
    x := rand.Intn(w / 4 - 1)
    y := rand.Intn(h / 2 - 1)
    p = point{x, y}
    for _, n := range not {
      if p == n {
        found = true
        continue
      }
    }
    if !found {
      print("food: ", p)
      return &p
    }
  }
  return nil
}
