package main

import (
  "fmt"
  "os"
  "log"
)

var logger *log.Logger

func main() {
  f, err := os.Create("./log")
  if err != nil {
    exit(err)
  }
  defer f.Close()
  logger = log.New(f, "", 0)

  snk := NewSnake()
  snk.Run()
}

func exit(e error) {
  fmt.Println(e)
  os.Exit(1)
}

func p(v ...interface{}) {
  logger.Print(v...)
}
