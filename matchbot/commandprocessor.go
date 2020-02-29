package main

import (
	"fmt"
	"regexp"
)

/* Set up regexes globally */
type RegexPatterns struct {
	ready *regexp.Regexp
}

var patterns *RegexPatterns // Too wasteful to pass to every function

func Analyze_Logs(logChan <-chan *PassedLogs) {
	var pass *PassedLogs
	patterns = Init_Regex()

	for {
		pass = <- logChan

		switch pass.typ {
		case COMMAND:
			Handle_Commands(pass)
		}
	}
}

func Init_Regex() *RegexPatterns {
	return &RegexPatterns {
		ready: regexp.MustCompile(`(?m).*(\.|\!|\/)(ready|r|unpause)"`),
	}
}

func Determine_Server() {
	
}

func Handle_Commands(command *PassedLogs) {
	if patterns.ready.MatchString(command.validLog) {
		fmt.Println("True")
		return
	}
}

func Ready_Up() (bool) {
	// if (isWarmup || isPaused) && isInit {
		return true
	// }
}

func Check_Team(logLine string) int8 {
	return 3
}
