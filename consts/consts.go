package consts

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/joho/godotenv"
	"os"
	"pomodoro_twitch_bot/logger"
	"pomodoro_twitch_bot/twitch_api_wrapper"
	"time"
)

var (
	Bot     *twitch_api_wrapper.Bot
	Prefix  string = "!"
	Channel string = ""
)

func LoadPrefix(w *astilectron.Window) {
	err := godotenv.Load(".env")
	if err != nil {
		if err := bootstrap.SendMessage(w, "NO SETTINGS", "NO SETTINGS SET"); err != nil {
			panic(err)
		}
	}
	err = godotenv.Load(".env")
	for err != nil { // while the error isn't nil
		time.Sleep(time.Second * 10) // sleep 10 secs until the user inputs the settings in the ui
		err = godotenv.Load(".env")
	}

	Prefix = os.Getenv("PREFIX")
	logger.Log(fmt.Sprintf("Using %s as prefix", Prefix))
	Channel = os.Getenv("CHANNEL")
}
