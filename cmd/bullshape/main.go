package main

import (
	"bullshape/router"
	l "bullshape/utils/logger"
)

func main() {
	log := l.NewLogger("Initialise")
	log.Printf("Lets start")
	router.Run()
}
