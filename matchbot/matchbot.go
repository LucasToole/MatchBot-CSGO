package main

import (
	"net"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/LucasToole/SourceRcon-go/rcon"
)

type PassedLogs struct {
	conn net.Conn
	index int8
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
	fullAddr string
	addr string
	port string
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
	bot.fullAddr = (config.Get("Bot.ip").(string) + ":" + config.Get("Bot.port").(string))
	bot.port = config.Get("Bot.port").(string)

	serv.addr = config.Get("Server.ip").(string)
	serv.port = config.Get("Server.port").(string)
	serv.fullAddr = serv.addr + ":" + serv.port
	serv.rconPass = config.Get("Server.rconpass").(string)
	serv.isInit = true // Temp
	serv.isWarmup = true // Temp
	serv.isPaused = false

	/* Initiate Rcon Conenctions */
	serv.conn, _, err = rcon.RconInitConnection(serv.addr, serv.port, serv.rconPass)
	if err != nil {
		fmt.Println(err)
		os.Exit(2) // For now stop the program. When multiple servers becomes real just move on.
	}
	rcon.RconSend(serv.conn, 2, "logaddress_add " + bot.fullAddr)
	rcon.RconSend(serv.conn, 2, "say The match is being managed by " + bot.name + "!")
	rcon.RconSend(serv.conn, 2, "say Admin: Type .start to initiate a game")
}

func main() {
	logChan := make(chan *PassedLogs, 10)
	cmdQ := make(chan *CommandInfo, 10)

	Init_Server()
	
	go Run_server(logChan)
	go Analyze_Logs(logChan, cmdQ)

	var rcmd *CommandInfo
	
	for {
		rcmd = <-cmdQ
		rcon.RconSend(rcmd.conn, 2, rcmd.cmd)
	}
}

