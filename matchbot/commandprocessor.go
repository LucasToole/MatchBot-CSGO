package main

import (
	"fmt"
)

func Analyze_Commands(logChan <-chan PassedLogs) {
	var pass = PassedLogs{}

	for {
		pass = <- logChan
		fmt.Println(pass.who)
	}
} 
