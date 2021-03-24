/* Copyright (c) 2020 Lucas Toole. See LICENSE for details */

package main

import (
	"fmt"
	"regexp"
)

/* Set up regexes globally */
type RegexPatterns struct {
	ready *regexp.Regexp
	pause *regexp.Regexp
	unpause *regexp.Regexp
	team *regexp.Regexp
	start *regexp.Regexp
	force *regexp.Regexp
	end *regexp.Regexp
	leave *regexp.Regexp
	score *regexp.Regexp
}

var patterns *RegexPatterns // Too wasteful to pass to every function

func Analyze_Logs(logChan <-chan *PassedLogs, cmdQ chan<- *CommandInfo) {
	var pass *PassedLogs
	patterns = Init_Regex()

	for {
		pass = <- logChan

		fmt.Println(pass.validLog)

		pass.serv_index = Determine_Server(pass)

		switch pass.typ {
		case COMMAND:
			Handle_Commands(pass, cmdQ)
		}
	}
}

func Init_Regex() *RegexPatterns {
	return &RegexPatterns {
		ready: regexp.MustCompile(`(?m)(\.|\!|\/)(ready|r)`),
		pause: regexp.MustCompile(`(?m)(\.|\!|\/)(pause)`),
		unpause: regexp.MustCompile(`(?m)(\.|\!|\/)(unpause)`),
		team: regexp.MustCompile(`(?m)<(CT|TERRORIST|Spectator)>`),
		start: regexp.MustCompile(`(?m)(\.|\!|\/)(start|map|maps)`),
		force: regexp.MustCompile(`(?m)(\.|\!|\/)(force)`),
		end: regexp.MustCompile(`(?m)(\.|\!|\/)(end|stop)`),
		leave: regexp.MustCompile(`(?m)(\.|\!|\/)(leave|exit)`),
		score: regexp.MustCompile(`(?m)(\.|\!|\/)(score)`),
	}
}

func Determine_Server(pass *PassedLogs) int {
	for  i := range server {
		if pass.who == (server[i].addr + ":" + server[i].port) { // TODO: Not IPv6 safe
			return i
		}
	}
	return -1 // TODO: Make this a real error?
}

func Handle_Commands(command *PassedLogs, cmdQ chan<- *CommandInfo) {
	if patterns.ready.MatchString(command.validLog) {
		if Ready_Up(Check_Team(command.validLog), command.serv_index) {
			cmdQ <- &CommandInfo{
				conn: server[command.serv_index].conn,
				cmd: "say READY!",
			}
		}
		if Match_Ready(command.serv_index) {
			cmdQ <- &CommandInfo{
				conn: server[command.serv_index].conn,
				cmd: "say Match Starting!",
			}
			Start_Match(command.serv_index)
		}
		return
	}

	if patterns.pause.MatchString(command.validLog) {
		if Pause_Match(Check_Team(command.validLog), command.serv_index) {
			cmdQ <- &CommandInfo{
				conn: server[command.serv_index].conn,
				cmd: "mp_pause_match",
			}
		}
		return
	}

	if patterns.unpause.MatchString(command.validLog) {
		if Unpause_Match(Check_Team(command.validLog), command.serv_index) {
			cmdQ <- &CommandInfo{
				conn: server[command.serv_index].conn,
				cmd: "mp_unpause_match",
			}
		}
		return
	}

}

func Ready_Up(readyTeam int8, serv_index int) (bool) {
	if (server[serv_index].isWarmup || server[serv_index].isPaused) && server[serv_index].isInit {
		if readyTeam != server[serv_index].ready {
			server[serv_index].ready += readyTeam
		}

		return true
	}
	return false
}

func Match_Ready(serv_index int) (bool) {
	if server[serv_index].ready == 3 {
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

func Start_Match(serv_index int) {
	server[serv_index].isWarmup = false
	server[serv_index].isPaused = false
	server[serv_index].ready = 0
	fmt.Println("This func will soon start a match")
	return
}

func Pause_Match(team int8, serv_index int) (bool) {
	if !server[serv_index].isPaused && team > 0 {
		server[serv_index].isPaused = true
		server[serv_index].ready = 0
		return true
	}
	return false
}

func Unpause_Match(team int8, serv_index int) (bool) {
	if server[serv_index].isPaused {
		Ready_Up(team, serv_index)
		if Match_Ready(serv_index) {
			server[serv_index].isPaused = false
			return true
		}
	}
	return false
}
