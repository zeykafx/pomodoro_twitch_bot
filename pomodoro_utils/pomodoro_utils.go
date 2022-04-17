package pomodoro_utils

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	"twitch_bot/consts"
	"twitch_bot/twitch_api_wrapper"
)

type Pomo struct {
	Username     string `json:"username"`
	Duration     int    `json:"pomoDuration"`
	Task         string `json:"task"`
	EndTimestamp string `json:"end_timestamp"`
	Silent       bool   `json:"silent"`
}

func CheckUserPomoDb(message *twitch_api_wrapper.Message) {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	selectStmt := fmt.Sprintf("SELECT * FROM pomodoros WHERE username = '%s'", message.User.Name)
	rows, err := db.Query(
		selectStmt,
	)
	if err != nil {
		log.Fatalln(err)
	}
	var pomo Pomo

	for rows.Next() {
		err := rows.Scan(
			&pomo.Username,
			&pomo.Task,
			&pomo.EndTimestamp,
			&pomo.Duration,
			&pomo.Silent,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}

	endTimestampInt, err := strconv.ParseInt(pomo.EndTimestamp, 10, 64)
	timeLeft := time.Until(time.Unix(endTimestampInt, 0)).Minutes()

	if timeLeft > 0 && !pomo.Silent {
		err := message.Reply(fmt.Sprintf("You can do it @%s! You have %v minutes left on your pomodoro session", message.User.Name, math.Round(timeLeft)))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func PomoLoop(bot *twitch_api_wrapper.Bot) {
	for {
		allPomos := FetchDbPomos()
		for _, currentPomo := range allPomos {

			endTimestampInt, _ := strconv.ParseInt(currentPomo.EndTimestamp, 10, 64)
			timeLeft := time.Until(time.Unix(endTimestampInt, 0)).Minutes()

			if timeLeft <= 0 {
				go terminatePomo(currentPomo.Username)
				msg := fmt.Sprintf("@%s the time is up on your pomodoro session", currentPomo.Username)
				err := bot.Send(consts.Channel, msg)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}

		time.Sleep(time.Second * 5) // sleep for 5 seconds
	}
}

func FetchDbPomos() []Pomo {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM pomodoros",
	)
	if err != nil {
		log.Fatalln(err)
	}

	var allPomos []Pomo

	for rows.Next() {
		var currentPomo Pomo
		err := rows.Scan(
			&currentPomo.Username,
			&currentPomo.Task,
			&currentPomo.EndTimestamp,
			&currentPomo.Duration,
			&currentPomo.Silent,
		)
		if err != nil {
			log.Fatalln(err)
		}
		allPomos = append(allPomos, currentPomo)
	}
	return allPomos
}

func terminatePomo(username string) {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	statement, err := db.Prepare("DELETE FROM pomodoros WHERE username=?")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(username)
	if err != nil {
		log.Fatalln(err)
	}
}
