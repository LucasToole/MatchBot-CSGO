package main

import (
	"fmt"
)

func Analyze_Commands() {
	var pass = PassedLogs{}

	fmt.Println(logq.String())
	pass, _ = logq.PopFront().(PassedLogs)
	fmt.Println(pass.who)
}
