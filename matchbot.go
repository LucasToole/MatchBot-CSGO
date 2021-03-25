/* Copyright (c) 2020 Lucas Toole. See LICENSE for details */

package main

import (
	"net"
	"fmt"
	"os"
	"flag"

	"github.com/pelletier/go-toml"
	"github.com/LucasToole/SourceRcon-go/rcon"
)

type PassedLogs struct {
	conn net.Conn
	serv_index int
	who string
	validLog string
	typ int8 /* Shouldn't need more than this */
}

type CommandInfo struct {
	conn net.Conn
	cmd string
}

type BotInfo struct {
	fullAddr string
	port string
	name string
}

type ServerInfo struct {
	conn net.Conn
	est bool
	addr string
	port string
	rconPass string
	ready int8 /* 0 - none, 1 - CT, 2 - T, 3 - Both */
	isInit bool
	isWarmup bool
	isPaused bool
}

var bot *BotInfo = new(BotInfo)
var server = make([]ServerInfo, 32) /* Who needs more than 32 Servers TODO: Uncap? Custom value? */

func Init_Server(configfile *string) int {
	config, err := toml.LoadFile(*configfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	bot.name = config.Get("Bot.name").(string)
	// TODO: When CS:GO supports IPv6 the format for v6 addresses will probably need to be -> [2001:db8:beef::e]:13337
	bot.fullAddr = (config.Get("Bot.ip").(string) + ":" + config.Get("Bot.port").(string))
	bot.port = config.Get("Bot.port").(string)

	for i := 0; i < 1; i++ { // TODO: Get the number of servers in the config file
		server[i].addr = config.Get("Server.ip").(string)
		server[i].port = config.Get("Server.port").(string)
		server[i].rconPass = config.Get("Server.rconpass").(string)
		server[i].isInit = true // Temp
		server[i].isWarmup = true // Temp
		server[i].isPaused = false
		server[i].est = false
	}

	/* Initiate Rcon Conenctions */
	var connectionCount int = 0
	for i := range server { // TODO: Range loops through all 32 elements
		server[i].conn, _, err = rcon.RconInitConnection(server[i].addr, server[i].port, server[i].rconPass)
		if err != nil {
			continue
		}
		rcon.RconSend(server[i].conn, 2, "logaddress_add " + bot.fullAddr)
		rcon.RconSend(server[i].conn, 2, "say The match is being managed by " + bot.name + "!")
		rcon.RconSend(server[i].conn, 2, "say Admin: Type .start to initiate a game")
		server[i].est = true
		connectionCount++
	}
	if connectionCount == 0 {
		fmt.Println("Error: Could not connect to any servers")
		os.Exit(2)
	}
	return connectionCount
}

func main() {
	logChan := make(chan *PassedLogs, 10)
	cmdQ := make(chan *CommandInfo, 10)

	configfile := flag.String("cfg", "config.toml", "The configuration file to use")
	flag.Parse()

	fmt.Println("Running CS:GO MatchBot")
	fmt.Println("Initiating connection to server...")
	connectionCount := Init_Server(configfile)
	fmt.Println("RCON Connection established for: ")
	for i := range server {
		if server[i].est == true {
			fmt.Println(server[i].addr + ":" + server[i].port)
		}
	}

	fmt.Println("Starting Log Listener...")
	go Run_server(logChan)

	fmt.Println("Starting Log Analyzer...")
	go Analyze_Logs(logChan, cmdQ)

	fmt.Printf("Done! Matchbot running for %d servers\n", connectionCount)

	var rcmd *CommandInfo

	for {
		rcmd = <-cmdQ
		rcon.RconSend(rcmd.conn, 2, rcmd.cmd)
	}
}
