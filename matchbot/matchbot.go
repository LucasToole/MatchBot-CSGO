package main

import (
	"net"
	"fmt"
)

type PassedLogs struct {
	conn net.Conn
	who string
	validLog string
	typ int8 /* Shouldn't need more than this */
}

type CommandInfo struct {
	conn net.Conn
	cmd string
}

type ServerInfo struct {
	conn net.Conn
	addr string
	ready int8 /* 0 - none, 1 - CT, 2 - T, 3 - Both */
	isInit bool
	isWarmup bool
	isPaused bool
}

func main() {
	logChan := make(chan *PassedLogs, 10)
	cmdQ := make(chan *CommandInfo, 10)
	
	go Run_server(logChan)
	go Analyze_Logs(logChan)

	var rcmd *CommandInfo
	
	for {
		rcmd = <-cmdQ
		fmt.Println(rcmd.cmd)
	}
}

