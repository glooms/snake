package snake

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Snake struct {
	app *tview.Application
	box *tview.Box

	tick      time.Duration
	dstep     time.Duration
	pauseChan chan bool
	clockChan chan bool

	init bool

	started bool
	paused  bool
	over    bool

	m    *model
	food *point

	grid [][]byte
	w, h int

	score int
	speed int

	snakeStyle tcell.Style
	foodStyle  tcell.Style
	hurtStyle  tcell.Style
}

const (
	block = '\u2588'
)

func New() *Snake {
	snk := &Snake{}

	//Setup styles
	col1 := tcell.ColorForestGreen
	col2 := tcell.ColorDarkViolet
	col3 := tcell.ColorDarkGreen

	snk.snakeStyle = tcell.StyleDefault.Background(col1)
	snk.foodStyle = tcell.StyleDefault.Background(col2)
	snk.hurtStyle = tcell.StyleDefault.Background(col3)

	//Setup box
	b := tview.NewBox().SetBorder(true)
	b.SetInputCapture(snk.capture)
	b.SetDrawFunc(snk.draw)
	snk.box = b

	//Setup app
	snk.app = tview.NewApplication().SetRoot(b, false)
	snk.app.SetBeforeDrawFunc(func(sc tcell.Screen) bool {
		if !snk.init {
			w, h := sc.Size()
			snk.box.SetRect(0, 0, w, h)
			snk.init = true
			tview.Print(sc, "Hello!", 0, 0, w, 0, -1)
		}
		return false
	})

	//Setup threading and channels
	dur, _ := time.ParseDuration("10ms")
	snk.tick = dur * 10
	snk.dstep = dur
	snk.pauseChan = make(chan bool, 10)
	snk.clockChan = make(chan bool, 10)

	//Setup the game things
	snk.reset()

	return snk
}

func (snk *Snake) Run() {
	err := snk.app.Run()
	if err != nil {
		panic(err)
	}
}

func (snk *Snake) start() {
	if !snk.started {
		snk.pauseChan <- true
		snk.started = true
	}
}

func (snk *Snake) pause() {
	snk.pauseChan <- true
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
	snk.grid = make([][]byte, snk.w/2+1)
	for i := range snk.grid {
		snk.grid[i] = make([]byte, snk.h)
	}
}

func (snk *Snake) eat() {
	snk.m.grow()
	snk.food = newFood(snk.w, snk.h, snk.m.get())
	snk.score++
	if snk.speed*5 < snk.score && snk.speed < 9 {
		snk.speed++
		snk.tick -= snk.dstep
	}
}

func (snk *Snake) draw(screen tcell.Screen, x, y, w, h int) (xn, yn, wn, hn int) {
	if snk.food == nil {
		snk.food = newFood(w, h, snk.m.get())
		return
	}
	if !snk.started || snk.paused {
		var str string
		if snk.paused {
			str = "PRESS 'P' TO RESUME"
		} else {
			str = "PRESS ANY KEY TO START"
		}
		tview.Print(screen, str, 0, h/4, w, 1, tcell.ColorWhite)
		str = "CONTROLS:"
		tview.Print(screen, str, w/4, h/4+2, w/2, 0, tcell.ColorWhite)
		ctrls := []string{
			"PAUSE:          'P'",
			"",
			"\u2191   'UP'    'W' 'K'",
			"\u2193   'DOWN'  'S' 'J'",
			"\u2190   'LEFT'  'A' 'H'",
			"\u2192   'RIGHT' 'D' 'L'",
		}
		for i, c := range ctrls {
			tview.Print(screen, c, w/4, h/4+2+i, w/2, 2, tcell.ColorWhite)
		}

	}
	fill := func(px, py int, st tcell.Style) {
		r, cmb, _, _ := screen.GetContent(px, py)
		screen.SetContent(px, py, r, cmb, st)
		r, cmb, _, _ = screen.GetContent(px+1, py)
		screen.SetContent(px+1, py, r, cmb, st)
	}
	m := point{w / 4, h / 2}
	str := fmt.Sprintf("SCORE: %d\t SPEED: %d", snk.score, snk.speed)
	tview.Print(screen, str, 3*w/4, 1, w/4, 0, tcell.ColorWhite)
	head := snk.m.head()
	if head.x+m.x == m.x*2 || head.x+m.x == 0 || head.y+m.y == m.y*2 || head.y+m.y == 0 {
		if snk.started {
			snk.gameOver()
		}
	}
	snk.newGrid(w, h)
	for _, p := range snk.m.get() {
		x, y := p.x+m.x, p.y+m.y
		if snk.grid[x][y] == 1 {
			snk.gameOver()
		} else {
			snk.grid[x][y] = 1
		}
	}

	if snk.food != nil {
		x, y := snk.food.x+m.x, snk.food.y+m.y
		if snk.grid[x][y] == 1 {
			snk.eat()
		} else {
			snk.grid[x][y] = 2
		}
	}

	if snk.over {
		x, y := head.x+m.x, head.y+m.y
		snk.grid[x][y] = 3
	}

	for i, row := range snk.grid {
		for j, b := range row {
			switch b {
			case 1:
				fill(i*2, j, snk.snakeStyle)
			case 2:
				fill(i*2, j, snk.foodStyle)
			case 3:
				fill(i*2, j, snk.hurtStyle)
			}
		}
	}
	screen.HideCursor()
	if snk.over {
		str := "GAME OVER"
		tview.Print(screen, str, 0, h/2-1, w, 1, tcell.ColorWhite)
		str = "Press any key to try again"
		l := len(str)
		tview.Print(screen, str, 0, h/2+1, w, 1, tcell.ColorWhite)
		str = "SCORE:"
		tview.Print(screen, str, w/2-l/2, h/2+4, w, 0, tcell.ColorWhite)
		str = fmt.Sprintf("%d", snk.score)
		tview.Print(screen, str, w/2-l/2, h/2+4, l/2, 2, tcell.ColorRed)
		switch lvl := snk.score / 5; lvl {
		case 0:
			str = "Next time will go better!"
		case 1:
			str = "You can do better!"
		case 2:
			str = "It's a learning process."
		case 3:
			str = "You're getting there!"
		case 4:
			str = "Shoot for the stars!"
		case 5:
			str = "Well done!"
		case 6:
			str = "Good, really good!"
		case 7:
			str = "Ssssnake whisperer"
		case 8:
			str = "Were you born a reptile?"
		case 9:
			str = "WOW! Crazy good!"
		case 10:
			str = "Holy Snakes"
		case 11:
			str = "sNaK3_g0dz'"
		}
		tview.Print(screen, str, w/2-l/2, h/2+6, w, 0, tcell.ColorWhite)
	}
	return
}

func (snk *Snake) gameOver() {
	snk.over = true
	snk.pause()
}
