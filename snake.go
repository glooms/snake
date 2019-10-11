package main

import (
//  "fmt"
  "time"

  "github.com/rivo/tview"
  "github.com/gdamore/tcell"
)

type Snake struct{
  app *tview.Application
  box *tview.Box

  tick time.Duration
  pauseChan chan bool
  clockChan chan bool

  started bool
  paused bool
  over bool

  m *model
  food *point

  grid [][]byte
  w, h int

  score int

  snakeStyle tcell.Style
  snakeFullStyle tcell.Style
  foodStyle tcell.Style
  sfStyle tcell.Style
  fsStyle tcell.Style
}

const (
  bu = '\u2580'
  bf = '\u2588'
  bl = '\u2584'
)

func NewSnake() *Snake {
  snk := &Snake{}
  //Setup styles
  csnk := tcell.ColorTeal
  cfod := tcell.ColorFuchsia
  snk.snakeStyle = tcell.StyleDefault.Foreground(csnk)
  snk.snakeFullStyle = snk.snakeStyle.Background(csnk)
  snk.foodStyle = tcell.StyleDefault.Foreground(cfod)
  snk.sfStyle = snk.snakeStyle.Background(cfod)
  snk.fsStyle = snk.foodStyle.Background(csnk)

  //Setup box
  b := tview.NewBox().SetBorder(true)
  b.SetInputCapture(snk.capture)
  b.SetDrawFunc(snk.draw)
  snk.box = b

  //Setup app
  snk.app = tview.NewApplication().SetRoot(b, true)

  //Setup threading and channels
  dur, err := time.ParseDuration("100ms")
  if err != nil {
    exit(err)
  }
  snk.tick = dur
  snk.pauseChan = make(chan bool, 10)
  snk.clockChan = make(chan bool, 10)

  //Setup the game things
  snk.reset()

  return snk
}

func (snk *Snake) Run() {
  if err := snk.app.Run(); err != nil {
    panic(err)
  }
}

func (snk *Snake) start() {
  if !snk.started {
//    fmt.Println("start")
    snk.pauseChan<-true
    snk.started = true
  }
}

func (snk *Snake) pause() {
  snk.pauseChan<-true
  snk.paused = !snk.paused
  if !snk.paused {
    snk.clock()
  }
}

func (snk *Snake) clock() {
  go func() {
    <-snk.pauseChan
    for {
      select {
      case <-snk.pauseChan:
        return
      default:
        snk.clockChan <- true
        time.Sleep(snk.tick)
      }
    }
  }()
}

func (snk *Snake) reset() {
  snk.clock()
  go snk.updateThread()
  snk.m = newModel(6)
  snk.food = newFood(snk.w, snk.h, snk.m.get())
  snk.over = false
  snk.started = false
  snk.paused = false
}

func (snk *Snake) capture(event *tcell.EventKey) *tcell.EventKey {
  if !snk.started {
    snk.start()
    return nil
  }
  if snk.over {
    snk.reset()
  }
  switch k := event.Key(); k {
  case tcell.KeyUp:
    snk.m.turn(UP)
  case tcell.KeyDown:
    snk.m.turn(DOWN)
  case tcell.KeyRight:
    snk.m.turn(RIGHT)
  case tcell.KeyLeft:
    snk.m.turn(LEFT)
  }
  switch r := event.Rune(); r {
  case 'p':
    snk.pause()
  default:
    if snk.paused {
      return nil
    }
    switch r {
    case 'h', 'a':
      snk.m.turn(LEFT)
    case 'j', 's':
      snk.m.turn(DOWN)
    case 'k', 'w':
      snk.m.turn(UP)
    case 'l', 'd':
      snk.m.turn(RIGHT)
    }
  }
  return nil
}

func (snk *Snake) updateThread() {
  for !snk.over {
    <-snk.clockChan
    snk.m.move()
    snk.app.Draw()
  }
}

func (snk *Snake) newGrid(w, h int) {
  snk.w, snk.h = w, h
  snk.grid = make([][]byte, snk.w)
  for i := range snk.grid {
    snk.grid[i] = make([]byte, snk.h * 2)
  }
}

func (snk *Snake) eat() {
  snk.m.grow()
  snk.food = newFood(snk.w, snk.h, snk.m.get())
  snk.score++
}

func (snk *Snake) draw(screen tcell.Screen, x, y, w, h int) (xn, yn, wn, hn int) {
  if snk.over {
    p("Score: ", snk.score)
    return
  }
  if snk.food == nil {
    snk.food = newFood(w, h, snk.m.get())
  }
  paint := func(px, py int, r rune, st tcell.Style) {
    screen.SetContent(px, py, r, []rune{}, st)
  }
  m := point{w / 2, h}
  head := snk.m.head()
  if head.x + m.x == m.x * 2 - 1 || head.x + m.x == 0 || head.y + m.y  == (m.y - 1) * 2 || head.y + m.y == 1 {
    if snk.started {
      snk.gameOver()
    }
  }
  snk.newGrid(w, h)
  for _, p := range snk.m.get() {
    x, y := p.x + m.x, p.y + m.y
    snk.grid[x][y] = 1
  }
  snk.grid[snk.food.x + m.x][snk.food.y + m.y] = 2
  for i, row := range snk.grid {
    for j, b := range row {
      even := j & 1 == 0
      switch b {
      case 1:
        switch {
        case even && snk.grid[i][j + 1] == 1:
          snk.grid[i][j + 1] = 0
          paint(i, j / 2, bf, snk.snakeStyle)
        case even:
          if snk.grid[i][j + 1] == 2 {
            snk.grid[i][j + 1] = 0
            paint(i, j / 2, bu, snk.fsStyle)
          } else {
            paint(i, j / 2, bu, snk.snakeStyle)
          }
        default:
          if snk.grid[i][j - 1] == 2 {
            snk.grid[i][j - 1] = 0
            paint(i, j / 2, bl, snk.sfStyle)
          } else {
            paint(i, j / 2, bl, snk.snakeStyle)
          }
        }
      case 2:
        switch {
        case even:
          paint(i, j / 2, bu, snk.foodStyle)
        default:
          paint(i, j / 2, bl, snk.foodStyle)
        }
      }
    }
  }
  if head == *snk.food {
    snk.eat()
  }
  screen.HideCursor()
  return
}

func (snk *Snake) gameOver() {
  snk.over = true
  snk.pause()
}
