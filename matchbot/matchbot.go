package main

import (
	"net"
	"fmt"

	"github.com/pelletier/go-toml"
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

type BotInfo struct {
	addr string
	name string
}

type ServerInfo struct {
	conn net.Conn
	addr string
	rconPass string
	ready int8 /* 0 - none, 1 - CT, 2 - T, 3 - Both */
	isInit bool
	isWarmup bool
	isPaused bool
}

var serv *ServerInfo = new(ServerInfo)
var bot *BotInfo = new(BotInfo)

func Init_Server() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		fmt.Println(err)
	}
	bot.name = config.Get("Bot.name").(string)
	bot.addr = (config.Get("Bot.ip").(string) + ":" + config.Get("Bot.port").(string))

	serv.addr = (config.Get("Server.ip").(string) + ":" + config.Get("Server.port").(string))
	serv.rconPass = config.Get("Server.rconpass").(string)

	fmt.Println("Bot name: " + bot.name)
	fmt.Println("Bot address: " + bot.addr)

	fmt.Println("Server address: " + serv.addr)
	fmt.Println("Server Rcon Pass: " + serv.rconPass)
	
}

func main() {
	logChan := make(chan *PassedLogs, 10)
	cmdQ := make(chan *CommandInfo, 10)

	Init_Server()
	
	go Run_server(logChan)
	go Analyze_Logs(logChan)

	var rcmd *CommandInfo
	
	for {
		rcmd = <-cmdQ
		fmt.Println(rcmd.cmd)
	}
}

