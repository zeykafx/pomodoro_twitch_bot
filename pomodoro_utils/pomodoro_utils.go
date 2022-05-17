package pomodoro_utils

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/twitch_api_wrapper"
	"strconv"
	"time"
)

var Running bool = true

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
		err := message.Reply(fmt.Sprintf("You can do it @%s! You have %v minutes left on your pomodoro session! Use %spomo chat to stop those reminders", message.User.Name, math.Round(timeLeft), consts.Prefix))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func PomoLoop(bot *twitch_api_wrapper.Bot) {
	for Running {
		allPomos := FetchDbPomos()
		for _, currentPomo := range allPomos {

			endTimestampInt, _ := strconv.ParseInt(currentPomo.EndTimestamp, 10, 64)
			timeLeft := time.Until(time.Unix(endTimestampInt, 0)).Minutes()

			if timeLeft <= 0 {
				go func() {
					err := TerminatePomo(currentPomo.Username, false)
					if err != nil {
						panic(err)
					}
				}()
				msg := fmt.Sprintf("@%s the time is up on your pomodoro session", currentPomo.Username)
				err := bot.Send(consts.Channel, msg)
				if err != nil {
					panic(err)
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

func EditUserPomo(username string, newTask string, shouldSendMessage bool) error {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	statement, err := db.Prepare("UPDATE pomodoros SET task=? WHERE username=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newTask, username)
	if err != nil {
		return err
	}

	if shouldSendMessage {
		msg := fmt.Sprintf("@%s your task was edited to \"%s\"", username, newTask)
		err := consts.Bot.Send(consts.Channel, msg)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func TerminatePomo(username string, shouldSendMessage bool) error {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	statement, err := db.Prepare("DELETE FROM pomodoros WHERE username=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username)
	if err != nil {
		return err
	}

	if shouldSendMessage {
		msg := fmt.Sprintf("@%s your pomodoro session was cancelled", username)
		err := consts.Bot.Send(consts.Channel, msg)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
