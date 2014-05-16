package main

import (
  "time"
  "log"
  "strconv"
)

func main () {
  var y int
  var m time.Month
  var d int

  foo := "1"
  if y,m,d = time.Now().Date(); y > 1 {
    log.Println(y)
  }
  log.Println(y,m,d)
  log.Println(strconv.Atoi(foo))
}
