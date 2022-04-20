package main

import (
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/help_commands"
	"pomodoro_twitch_bot/pomoboard"
	"pomodoro_twitch_bot/pomobot"
	"pomodoro_twitch_bot/pomodoro"
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
