/* Copyright (c) 2020 Lucas Toole. See LICENSE for details */

package main

import (
	//"fmt"
	"regexp"
	"net"
)

/* Set up regexes globally */
type RegexPatterns struct {
	ready *regexp.Regexp
}

var patterns *RegexPatterns // Too wasteful to pass to every function

func Analyze_Logs(logChan <-chan *PassedLogs, cmdQ chan<- *CommandInfo) {
	var pass *PassedLogs
	patterns = Init_Regex()

	for {
		pass = <- logChan

		pass.conn, pass.index = Determine_Server(pass)
		
		switch pass.typ {
		case COMMAND:
			Handle_Commands(pass, cmdQ)
		}
	}
}

func Init_Regex() *RegexPatterns {
	return &RegexPatterns {
		ready: regexp.MustCompile(`(?m).*(\.|\!|\/)(ready|r|unpause)"`),
	}
}

func Determine_Server(pass *PassedLogs) (net.Conn, int8) {
	if pass.who == serv.fullAddr {
		return serv.conn, 1 // TEMP!!!
	}
	// Figure out later !
	return serv.conn, 1 // Should return the index of server for easy future access
}

func Handle_Commands(command *PassedLogs, cmdQ chan<- *CommandInfo) {
	if patterns.ready.MatchString(command.validLog) {
		if Ready_Up() {
			cmdQ <- &CommandInfo{
				conn: command.conn,
				cmd: "say READY!",
			}
		}
		return
	}
}

func Ready_Up() (bool) { // Should take index of server in future
	if (serv.isWarmup || serv.isPaused) && serv.isInit {
		return true
	}
	return false
}

func Check_Team(logLine string) int8 {
	return 3
}
