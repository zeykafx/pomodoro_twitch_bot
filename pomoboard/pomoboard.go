package pomoboard

import (
	"fmt"
	"io/ioutil"
	"math"
	"pomodoro_twitch_bot/pomodoro_utils"
	"strconv"
	"time"
)

var Running bool = true

func StartPomoBoard() {
	type PomoWithTimeLeft struct {
		pomodoro_utils.Pomo
		TimeLeft float64 `json:"time_left"`
	}

	var fileName string = "pomoboard.txt"

	for Running {
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

		err := ioutil.WriteFile(fileName, []byte(stringRepr), 0644)
		if err != nil {
			panic(err)
		} // overwrite the whole file every 5 seconds
		time.Sleep(time.Second * 5)
	}

}
