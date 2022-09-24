package pomobot

import (
	"fmt"
	"log"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/logger"
	"pomodoro_twitch_bot/pomodoro_utils"
	"pomodoro_twitch_bot/twitch_api_wrapper"
	"strings"
)

var (
	//Bot            *twitch_api_wrapper.Bot
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
	//err := godotenv.Load(".env")
	//CheckErr(err)
	//TOKEN := os.Getenv("TOKEN")

	// creating the bot
	consts.Bot = twitch_api_wrapper.NewBot(consts.Token, "pomobot", []string{consts.Channel})

	consts.Bot.OnLogin(func(bot *twitch_api_wrapper.Bot) {
		logger.Log("Logged in!")
		go pomodoro_utils.PomoLoop(consts.Bot)
	})

	// reacts to all the commands and log them
	consts.Bot.OnMessage(func(bot *twitch_api_wrapper.Bot, message *twitch_api_wrapper.Message) {
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
	consts.Bot.Run()
}

func StopBot() {
	if consts.Bot != nil {
		consts.Bot.Client.Close()
	}
	logger.Log("closed bot connection.")
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
