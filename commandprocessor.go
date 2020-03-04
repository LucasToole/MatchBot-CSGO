/* Copyright (c) 2020 Lucas Toole. See LICENSE for details */

package main

import (
	"fmt"
	"regexp"
	"net"
)

/* Set up regexes globally */
type RegexPatterns struct {
	ready *regexp.Regexp
	team *regexp.Regexp
}

var patterns *RegexPatterns // Too wasteful to pass to every function

func Analyze_Logs(logChan <-chan *PassedLogs, cmdQ chan<- *CommandInfo) {
	var pass *PassedLogs
	patterns = Init_Regex()

	for {
		pass = <- logChan

		fmt.Println(pass.validLog)

		pass.conn, pass.index = Determine_Server(pass)
		
		switch pass.typ {
		case COMMAND:
			Handle_Commands(pass, cmdQ)
		}
	}
}

func Init_Regex() *RegexPatterns {
	return &RegexPatterns {
		ready: regexp.MustCompile(`(?m)(\.|\!|\/)(ready|r|unpause)`),
		team: regexp.MustCompile(`(?m)<(CT|TERRORIST|Spectator)>`),
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
		if Ready_Up(Check_Team(command.validLog)) {
			cmdQ <- &CommandInfo{
				conn: command.conn,
				cmd: "say READY!",
			}
		}
		if Match_Ready() {
			cmdQ <- &CommandInfo{
				conn: command.conn,
				cmd: "say Match Starting!",
			}
			Start_Match()
		}
		return
	}
}

func Ready_Up(readyTeam int8) (bool) { // Should take index of server in future
	if (serv.isWarmup || serv.isPaused) && serv.isInit {
		if readyTeam != serv.ready {
			serv.ready += readyTeam
		}

		return true
	}
	return false
}

func Match_Ready() (bool) { // Needs server index
	if serv.ready == 3 {
		return true
	}
	return false
}

func Check_Team(logLine string) int8 {
	team := patterns.team.FindString(logLine)

	if team == "<CT>" {
		return 1
	}
	if team == "<TERRORIST>" {
		return 2
	}
	if team == "<Spectator>" {
		return 0
	}
	return -1 /* Somethings wrong, shouldn't reach this */
}

func Start_Match() {
	fmt.Println("This func will soon start a match")
	return
}