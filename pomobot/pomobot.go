package pomobot

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"twitch_bot/consts"
	"twitch_bot/logger"
	"twitch_bot/pomodoro_utils"
	"twitch_bot/twitch_api_wrapper"
)

var (
	Bot            *twitch_api_wrapper.Bot
	commandHandler CommandHandler = CommandHandler{commandHandlers: map[string]func(*twitch_api_wrapper.Bot, *twitch_api_wrapper.Message){}} // init empty map
)

type CommandHandler struct {
	commandHandlers map[string]func(*twitch_api_wrapper.Bot, *twitch_api_wrapper.Message)
}

// commands:
// - [prefix]pomo [time] [task]
// - [prefix]pomo end
// - [prefix]pomo check
// - [prefix]pomo chat
// - [prefix]pomo add/remove [time]

func InitBot() {
	err := godotenv.Load(".env")
	CheckErr(err)
	TOKEN := os.Getenv("TOKEN")

	// creating the bot
	Bot = twitch_api_wrapper.NewBot(TOKEN, "pomobot", []string{consts.Channel})

	Bot.OnLogin(func(bot *twitch_api_wrapper.Bot) {
		logger.Log("Logged in!")
		go pomodoro_utils.PomoLoop(Bot)
	})

	// reacts to all the commands and log them
	Bot.OnMessage(func(bot *twitch_api_wrapper.Bot, message *twitch_api_wrapper.Message) {
		logger.Log(fmt.Sprintf("'%s' sent '%s'", message.User.Name, message.Message))

		// if the message is not a command
		if !strings.HasPrefix(message.Message, consts.Prefix) {
			// checks if the user has a pomo running, if so we warn them in chat
			pomodoro_utils.CheckUserPomoDb(message)
		}

		// iterate over the command handlers and call them if it matches the message
		for commandName, handler := range commandHandler.commandHandlers {
			if strings.Split(message.Message, " ")[0] == prefixedCommand(commandName) {
				go handler(bot, message)
			}
		}
	})
	// run the bot
	Bot.Run()
}

// Accept is used to add a function as a command handler
func Accept(commandFunc func(*twitch_api_wrapper.Bot, *twitch_api_wrapper.Message), commandName string) {
	commandHandler.commandHandlers[commandName] = commandFunc
}

// returns the prefixed version of the command
func prefixedCommand(command string) string {
	return fmt.Sprintf("%s%s", consts.Prefix, command)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
