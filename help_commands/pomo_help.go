package help_commands

import (
	"fmt"
	"strings"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/pomobot"
	"pomodoro_twitch_bot/twitch_api_wrapper"
)

// commands:
// - [prefix]pomo [time] [task]
// - [prefix]pomo end
// - [prefix]pomo check
// - [prefix]pomo chat
// - [prefix]pomo add/remove [time]

func Help(bot *twitch_api_wrapper.Bot, message *twitch_api_wrapper.Message) {
	splitCommand := strings.Split(message.Message, " ")
	if len(splitCommand) == 2 {
		if splitCommand[1] == "pomo" {
			// send help message for the pomo command
			err := message.Reply(fmt.Sprintf("To start your own pomodoro timer use \"%spomo time task (optional)\", to end your timer early use \"%spomo end\", to check your timer use \"%spomo check\", if you don't want the bot to warn when talking with a timer running use \"%spomo chat\", and to change the timer use \"%spomo add/remove time\"", consts.Prefix, consts.Prefix, consts.Prefix, consts.Prefix, consts.Prefix))
			pomobot.CheckErr(err)
			return
		}
	} else { // just the !help
		err := message.Reply(fmt.Sprintf("Commands available:    - %spomo: create a custom pomodoro timer and appear on stream (use \"%shelp pomo\" to get help using the pomo command)", consts.Prefix, consts.Prefix))
		pomobot.CheckErr(err)
		return
	}
}
