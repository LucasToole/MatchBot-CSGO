package main

import (
	"time"
	
	"github.com/phf/go-queue/queue" /* may replace this with a DIY one*/
)

type PassedLogs struct {
	who string
	goodLog string // Let's change that variable name :)
	typ int8 /* Shouldn't need more than this */
}

var logq = queue.New()

func main() {
	go Run_server()
	time.Sleep(10 * time.Second)
	Analyze_Commands()
	
}

