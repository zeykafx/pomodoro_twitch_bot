package main

import (
	"twitch_bot/consts"
	"twitch_bot/help_commands"
	"twitch_bot/pomoboard"
	"twitch_bot/pomobot"
	"twitch_bot/pomodoro"
)

func initCommands() {
	// add the functions that handle commands here
	pomobot.Accept(pomodoro.HandlePomoCommand, "pomo")
	pomobot.Accept(help_commands.Help, "help")
	// starts the web server for the pomo board
	go pomoboard.StartPomoBoard()
}

func main() {
	consts.LoadPrefix()
	initCommands()
	pomobot.InitBot()
}
