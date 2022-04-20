package consts

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"pomodoro_twitch_bot/logger"
	"pomodoro_twitch_bot/setup"
)

var (
	Prefix  string = "!"
	Channel string = "zeykafx"
)

func LoadPrefix() {
	err := godotenv.Load(".env")
	if err != nil {
		setup.FirstTimeSetup()
	}
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	Prefix = os.Getenv("PREFIX")
	logger.Log(fmt.Sprintf("Using %s as prefix", Prefix))
	Channel = os.Getenv("CHANNEL")
}
