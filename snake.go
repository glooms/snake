package main

import (
  "fmt"
  "time"

  "github.com/rivo/tview"
  "github.com/gdamore/tcell"
)

type Snake struct{
  app *tview.Application
  box *tview.Box

  tick time.Duration
  dstep time.Duration
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
  speed int

  snakeStyle tcell.Style
  foodStyle tcell.Style
  hurtStyle tcell.Style
  textStyle tcell.Style
}

const (
  block = '\u2588'
)

func NewSnake() *Snake {
  snk := &Snake{}
  //Setup styles
  csnk := tcell.ColorForestGreen
  cfod := tcell.ColorDarkViolet
  chrt := tcell.ColorDarkGreen
  snk.snakeStyle = tcell.StyleDefault.Background(csnk)
  snk.foodStyle = tcell.StyleDefault.Background(cfod)
  snk.hurtStyle = tcell.StyleDefault.Background(chrt)

  //Setup box
  b := tview.NewBox().SetBorder(true)
  b.SetInputCapture(snk.capture)
  b.SetDrawFunc(snk.draw)
  snk.box = b

  //Setup app
  snk.app = tview.NewApplication().SetRoot(b, true)

  //Setup threading and channels
  dur, err := time.ParseDuration("10ms")
  if err != nil {
    exit(err)
  }
  snk.tick = dur * 10
  snk.dstep = dur
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
  snk.score = 0
  snk.speed = 1
  snk.tick = snk.dstep * 10
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
  snk.grid = make([][]byte, snk.w / 2 + 1)
  for i := range snk.grid {
    snk.grid[i] = make([]byte, snk.h)
  }
}

func (snk *Snake) eat() {
  snk.m.grow()
  snk.food = newFood(snk.w, snk.h, snk.m.get())
  snk.score++
  if snk.speed * 5 < snk.score && snk.speed < 9 {
    snk.speed++
    snk.tick -= snk.dstep
  }
}

func (snk *Snake) draw(screen tcell.Screen, x, y, w, h int) (xn, yn, wn, hn int) {
  if snk.food == nil {
    snk.food = newFood(w, h, snk.m.get())
  }
  fill := func(px, py int, st tcell.Style) {
    r, cmb, _, _ := screen.GetContent(px, py)
    screen.SetContent(px, py, r, cmb, st)
    r, cmb, _, _ = screen.GetContent(px + 1, py)
    screen.SetContent(px + 1, py, r, cmb, st)
  }
  m := point{w / 4, h / 2}
  str := fmt.Sprintf("SCORE: %d\t SPEED: %d", snk.score, snk.speed)
  for i, r := range str {
    screen.SetContent(3 * w / 4 + i, 1, r, []rune{}, snk.textStyle)
  }
  head := snk.m.head()
  if head.x + m.x == m.x * 2 || head.x + m.x == 0 || head.y + m.y == m.y * 2 || head.y + m.y == 0 {
    if snk.started {
      snk.gameOver()
    }
  }
  snk.newGrid(w, h)
  for _, p := range snk.m.get() {
    x, y := p.x + m.x, p.y + m.y
    print("snk: ", x, " ", y)
    if snk.grid[x][y] == 1 {
      snk.gameOver()
    } else {
      snk.grid[x][y] = 1
    }
  }

  if snk.food != nil {
    x, y := snk.food.x + m.x, snk.food.y + m.y
    print("food: ", x, " ", y)
    if snk.grid[x][y] == 1 {
      snk.eat()
    } else {
      snk.grid[x][y] = 2
    }
  }

  if snk.over {
    x, y := head.x + m.x, head.y + m.y
    snk.grid[x][y] = 3
  }

  for i, row := range snk.grid {
    for j, b := range row {
      switch b {
      case 1:
        fill(i * 2, j, snk.snakeStyle)
      case 2:
        fill(i * 2, j, snk.foodStyle)
      case 3:
        fill(i * 2, j, snk.hurtStyle)
      }
    }
  }
  screen.HideCursor()
  if snk.over {
    str := fmt.Sprintf("GAME OVER")
    for i, r := range str {
      screen.SetContent(w/2 - len(str) / 2 + i, h/2 - 1, r, []rune{}, snk.textStyle)
    }
    str = fmt.Sprintf("Press any key to try again")
    for i, r := range str {
      screen.SetContent(w/2 - len(str) / 2 + i, h/2 + 1, r, []rune{}, snk.textStyle)
    }
    str = fmt.Sprintf("SCORE: %d", snk.score)
    for i, r := range str {
      screen.SetContent(w/2 - len(str) / 2 + i, h/2 + 3, r, []rune{}, snk.textStyle)
    }
    var adj string
    switch {
    case snk.score < 5:
      adj = "poor"
    case snk.score < 10:
      adj = "mediochre"
    case snk.score < 20:
      adj = "decent"
    case snk.score < 30:
      adj = "good"
    case snk.score < 40:
      adj = "respectable"
    case snk.score < 50:
      adj = "crazy"
    case snk.score < 100:
      adj = "INSANE"
    case snk.score > 100:
      adj = "GODLIKE"
    }
    if snk.score > 100 {
      str = "You are: sNaKE_g0d"
    } else {
      str = fmt.Sprintf("It's %s.", adj)
    }
    for i, r := range str {
      screen.SetContent(w/2 - len(str) / 2 + i, h/2 + 4, r, []rune{}, snk.textStyle)
    }
  }
  return
}

func (snk *Snake) gameOver() {
  snk.over = true
  snk.pause()
}
