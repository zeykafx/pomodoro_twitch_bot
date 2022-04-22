package main

import (
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/pomodoro_utils"
	"pomodoro_twitch_bot/setup"
	"runtime"
	"strconv"
	"time"
)

var shouldWriteToFile bool = false

// handleMessages handles messages sent to and from the React app
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {

	case "START_FILE":
		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			return nil, err
		}
		fmt.Println(payload)
		if !shouldWriteToFile {
			shouldWriteToFile = true
			go writeToFile()

		} else {
			// already writing to file
			// do nothing
		}
		return "Started writing to file!", nil

	case "STOP_FILE":
		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			return nil, err
		}
		fmt.Println(payload)
		if shouldWriteToFile {
			shouldWriteToFile = false
		}

		return "Stopped writing to file!", nil

	case "RUNNING_POMOS":
		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			panic(err)
		}
		fmt.Println(payload)
		stringRepr := getRunningPomos() // get all the pomos as a string and return that to the app
		return stringRepr, nil

	case "STATUS":
		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			panic(err)
		}
		fmt.Println(payload)
		var status string
		// dammit I want my ternary operator
		if pomodoro_utils.Running {
			status = "on"
		} else {
			status = "off"
		}
		return status, nil

	case "SET_SETTINGS":

		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			return "didnt' save", err
		}
		var settings settingsJson
		err = json.Unmarshal([]byte(payload), &settings)
		if err != nil {
			return "didnt' save", err
		}
		err = setup.WriteSettingsToFile(settings.Token, settings.Prefix, settings.Channel)
		if err != nil {
			return "didnt' save", err
		}
		return "saved settings", nil

	case "GET_SETTINGS":
		var settings settingsJson
		settings.Token = os.Getenv("TOKEN")
		settings.Prefix = consts.Prefix
		settings.Channel = consts.Channel
		settingsString, err := json.Marshal(&settings)
		if err != nil {
			return nil, err
		}
		return string(settingsString), nil

	case "URL":
		var payload string
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			panic(err)
		}
		openbrowser(payload)
		return nil, nil
	}

	return
}

// openbrowser opens the default browser on all platforms, this is used instead of opening links in electron
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

type settingsJson struct {
	Token   string `json:"token"`
	Prefix  string `json:"prefix"`
	Channel string `json:"channel"`
}

type PomoWithTimeLeft struct {
	pomodoro_utils.Pomo
	TimeLeft float64 `json:"time_left"`
}

func getRunningPomos() string {
	allPomos := pomodoro_utils.FetchDbPomos()
	var allPomosWithTimeLeft []PomoWithTimeLeft
	for _, pomo := range allPomos {
		endTimestampInt, _ := strconv.ParseInt(pomo.EndTimestamp, 10, 64)
		timeLeft := time.Until(time.Unix(endTimestampInt, 0)).Minutes()
		currPomo := PomoWithTimeLeft{
			pomo, timeLeft,
		}
		allPomosWithTimeLeft = append(allPomosWithTimeLeft, currPomo)
	}

	var stringRepr string
	for _, pomoWithTimeLeft := range allPomosWithTimeLeft {
		stringRepr += fmt.Sprintf("%s: \"%s\" | %v minutes left\n", pomoWithTimeLeft.Username, pomoWithTimeLeft.Task, math.Round(pomoWithTimeLeft.TimeLeft))
	}
	return stringRepr
}

func writeToFile() {
	var fileName string = "pomoboard.txt"

	for shouldWriteToFile {
		var stringRepr string = getRunningPomos()

		err := ioutil.WriteFile(fileName, []byte(stringRepr), 0644)
		if err != nil {
			panic(err)
		}
		// overwrite the whole file every 5 seconds
		time.Sleep(time.Second * 5)
	}

}
