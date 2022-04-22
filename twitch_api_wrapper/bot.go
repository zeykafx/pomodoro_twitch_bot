package twitch_api_wrapper

import (
	"log"
)

type EventHandler struct {
	messageHandlers []func(*Bot, *Message)
}

type Bot struct {
	Client   *Client
	host     string
	onLogin  func(*Bot)
	events   EventHandler
	channels []string
}

type User struct {
	ID   string
	Name string
}

type Message struct {
	ID      string
	Channel string
	User    *User
	Message string
	Bot     *Bot
}

func ParseMessage(command *Command, bot *Bot) *Message {
	return &Message{
		ID:      command.Tags["id"],
		Channel: command.Args[0][1:],
		User:    &User{ID: command.Tags["user-id"], Name: command.Tags["display-name"]},
		Message: command.Suffix,
		Bot:     bot,
	}
}

func NewBot(token string, nick string, channels []string) *Bot {
	client := Client{Token: token, Nick: nick}
	return &Bot{Client: &client, host: "irc.chat.twitch.tv:6667", events: EventHandler{}, channels: channels}
}

func (message *Message) Reply(msg string) error {
	err := message.Bot.SendMessage(&Message{Message: msg, Channel: message.Channel})
	if err != nil {
		return err
	}
	return nil
}

func (message *Message) Delete() error {
	err := message.Bot.DeleteMessage(&Message{ID: message.ID, Channel: message.Channel})
	if err != nil {
		return err
	}
	return nil
}

func (message *Message) Ban() error {
	err := message.Bot.BanUser(message.Channel, message.User.Name)
	if err != nil {
		return err
	}
	return nil
}

func (event *EventHandler) configure(bot *Bot) {

	bot.Client.AddHandler("PRIVMSG", func(command *Command) bool {
		message := ParseMessage(command, bot)
		for _, handler := range event.messageHandlers {
			go handler(bot, message)
		}
		return true
	})

}

func (bot *Bot) OnLogin(f func(*Bot)) {
	bot.onLogin = f
}

func (bot *Bot) SendMessage(message *Message) error {
	err := bot.Client.Send(&Command{Command: "PRIVMSG", Args: []string{"#" + message.Channel}, Suffix: message.Message})
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) DeleteMessage(message *Message) error {
	err := bot.SendMessage(&Message{Channel: message.Channel, Message: "/delete " + message.ID})
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) Join(channel string) error {
	err := bot.Client.Join("#" + channel)
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) Send(channel string, content string) error {
	err := bot.SendMessage(&Message{Channel: channel, Message: content})
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) BanUser(channel string, user string) error {
	err := bot.SendMessage(&Message{Channel: channel, Message: "/ban " + user})
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) OnMessage(f func(*Bot, *Message)) {
	bot.events.messageHandlers = append(bot.events.messageHandlers, f)
}

func (bot *Bot) GetClient() *Client {
	return bot.Client
}

func (bot *Bot) Run() {
	for {
		err := bot.Start()
		if err != nil {
			log.Printf("Bot error: %s\n", err)
		}
	}
}

func (bot *Bot) Start() error {
	defer bot.Client.Close()
	err := bot.Client.Connect(bot.host)
	if err != nil {
		return err
	}

	err = bot.Client.Auth()
	if err != nil {
		return err
	}

	bot.events.configure(bot)

	bot.Client.AddHandler("PING", func(command *Command) bool {
		err := bot.Client.Send(&Command{Command: "PONG", Suffix: "tmi.twitch.tv"})
		if err != nil {
			return false
		}
		return true
	})

	bot.Client.AddHandler("376", func(command *Command) bool {
		err = bot.GetClient().CapReq("twitch.tv/tags twitch.tv/commands")
		if err != nil {
			return false
		}
		for _, channel := range bot.channels {
			err = bot.Join(channel)
			if err != nil {
				return false
			}
		}
		if bot.onLogin != nil {
			bot.onLogin(bot)
		}
		return true
	})

	err = bot.Client.Handle()
	if err != nil {
		return err
	}

	return nil
}
