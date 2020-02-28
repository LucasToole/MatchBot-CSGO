package main

import (
	
)

type PassedLogs struct {
	who string
	validLog string
	typ int8 /* Shouldn't need more than this */
}

func main() {
	logChan := make(chan PassedLogs, 10)
	
	go Run_server(logChan)
	Analyze_Commands(logChan)
}

