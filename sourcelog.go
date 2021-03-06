/* Copyright (c) 2020 Lucas Toole. See LICENSE for details */

package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

const (
	COMMAND int8 = 1
	SCOREUPDATE int8 = 2
)

func Run_server(logChan chan<- *PassedLogs) {
	packet, err := net.ListenPacket("udp", ":" + bot.port)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	/* Lazy Processing here just to not flood the analyzer */
	command := regexp.MustCompile(`(?m).*".+<[0-9]+><STEAM_.+><(CT|TERRORIST|Spectator)>" say(_team)? "(\.|\!|\/).+"`)
	score := regexp.MustCompile(`(?m).+ triggered "SFUI_Notice_.+" \(CT "[0-9]+"\) \(T "[0-9]+"\)`)
	
	for {
		buf := make([]byte, 1024)
		_, addrr, err := packet.ReadFrom(buf)
		if err != nil {
			continue
		}
		
		if command.MatchString(string(buf[30:])) == true {
			logChan <- &PassedLogs{
			who: addrr.String(),
			validLog: string(buf[30:]),
			typ: COMMAND,
			}
			continue
		}

		if score.MatchString(string(buf[30:])) == true {
			logChan <- &PassedLogs{
				who: addrr.String(),
				validLog: string(buf[30:]),
				typ: SCOREUPDATE,
			}
			continue
		}
	}
}
