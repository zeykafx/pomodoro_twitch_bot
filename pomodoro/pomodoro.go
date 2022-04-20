package pomodoro

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/pomobot"
	"pomodoro_twitch_bot/twitch_api_wrapper"
	"strconv"
	"strings"
	"time"
)

// commands:
// - [prefix]pomo [time] [task]
// - [prefix]pomo end
// - [prefix]pomo check
// - [prefix]pomo chat
// - [prefix]pomo add/remove [time]

type Pomo struct {
	Username     string `json:"username"`
	Duration     int    `json:"pomoDuration"`
	Task         string `json:"task"`
	EndTimestamp string `json:"end_timestamp"`
	Silent       bool   `json:"silent"`
}

func HandlePomoCommand(bot *twitch_api_wrapper.Bot, message *twitch_api_wrapper.Message) {

	splitCommand := strings.Split(message.Message, " ")
	for _, elem := range splitCommand {
		elem = strings.ToLower(elem)
	}
	//fmt.Println("[" + strings.Join(splitCommand, ",") + "]")

	if len(splitCommand) == 1 {
		err := message.Reply(fmt.Sprintf("@%s To start your own pomodoro session use \"%spomo [duration] [task - optional]\"", message.User.Name, consts.Prefix))
		pomobot.CheckErr(err)
		return
	}

	// start a pomo
	if _, err := strconv.Atoi(splitCommand[1]); len(splitCommand) > 1 && err == nil {
		var pomo Pomo

		// parse pomoDuration
		pomoDuration, err := strconv.Atoi(splitCommand[1])
		if err != nil {
			err = message.Reply(fmt.Sprintf("Invalid format, the correct format is \"%spomo [duration] [task - optional]\"", consts.Prefix))
			pomobot.CheckErr(err)
			return
		}
		pomo.Duration = pomoDuration

		var task string
		// parse the task
		if len(splitCommand) > 2 {
			task = strings.Join(splitCommand[2:], " ") // get all the text
		} else {
			task = "Work/Study"
		}

		pomo.Task = task

		// set username to message author
		pomo.Username = message.User.Name

		// make the pomo silent for the broadcaster but not for other users
		if strings.ToLower(pomo.Username) == consts.Channel {
			pomo.Silent = true
		} else {
			pomo.Silent = false
		}

		// get the timestamp of now + pomoDuration aka the pomo end time
		endTime := time.Now().Add(time.Minute * time.Duration(pomoDuration))
		pomo.EndTimestamp = strconv.FormatInt(endTime.Unix(), 10)

		pomoAlreadyExists := insertIntoDb(&pomo, bot)
		if pomoAlreadyExists != nil {
			return
		}

		err = message.Reply(fmt.Sprintf("@%s Started pomodoro session on \"%s\" for %v minutes, good luck!", pomo.Username, pomo.Task, pomo.Duration))
		pomobot.CheckErr(err)
	} else {
		// if its like !pomo [command]

		if splitCommand[1] == "end" || splitCommand[1] == "cancel" || splitCommand[1] == "stop" || splitCommand[1] == "finish" { // ENDS A POMO
			terminatePomo(message.User.Name, bot)

		} else if splitCommand[1] == "check" { // SENDS THE TIME LEFT IN THE CURRENT POMO
			pomo := getPomoFromDb(message.User.Name)

			if pomo.EndTimestamp != "" {
				endTimestampInt, _ := strconv.ParseInt(pomo.EndTimestamp, 10, 64)
				timeLeft := time.Until(time.Unix(endTimestampInt, 0)).Minutes()

				msg := fmt.Sprintf("@%s you have %v minutes left on your pomodoro session", message.User.Name, math.Round(timeLeft))
				err := bot.Send(consts.Channel, msg)
				pomobot.CheckErr(err)
			} else {
				// no pomos found
				msg := fmt.Sprintf("@%s Couldn't find any pomodoro session running", message.User.Name)
				err := bot.Send(consts.Channel, msg)
				pomobot.CheckErr(err)
			}
		} else if splitCommand[1] == "chat" || splitCommand[1] == "silent" || splitCommand[1] == "silence" || splitCommand[1] == "mod" { // SILENT POMO MODE
			updateSilentStatusPomoDb(true, message.User.Name)
			err = message.Reply("You won't be reminded to focus for this pomo.")
			pomobot.CheckErr(err)

		} else if splitCommand[1] == "add" || splitCommand[1] == "plus" { // ADD TIME TO POMO
			if timeToAdd, err := strconv.Atoi(splitCommand[2]); len(splitCommand) >= 3 && err == nil {

				changePomoTime(message.User.Name, timeToAdd, bot)

			} else {
				msg := fmt.Sprintf("@%s invalid format, use \"%spomo add time\"", message.User.Name, consts.Prefix)
				err := bot.Send(consts.Channel, msg)
				pomobot.CheckErr(err)
			}
		} else if splitCommand[1] == "remove" || splitCommand[1] == "minus" { // REMOVES TIME TO POMO
			if timeToAdd, err := strconv.Atoi(splitCommand[2]); len(splitCommand) >= 3 && err == nil {

				changePomoTime(message.User.Name, -timeToAdd, bot)

			} else {
				msg := fmt.Sprintf("@%s invalid format, use \"%spomo remove time\"", message.User.Name, consts.Prefix)
				err := bot.Send(consts.Channel, msg)
				pomobot.CheckErr(err)
			}
		}
	}
}

func changePomoTime(username string, timeToAdd int, bot *twitch_api_wrapper.Bot) {
	pomoToUpdate := getPomoFromDb(username)
	endTimestampInt, _ := strconv.ParseInt(pomoToUpdate.EndTimestamp, 10, 64)
	timeLeft := time.Unix(endTimestampInt, 0).Add(time.Duration(timeToAdd) * time.Minute)
	newEndTimestamp := timeLeft.Unix()

	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	pomobot.CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("update pomodoros set end_timestamp=? where username=?")
	pomobot.CheckErr(err)

	res, err := stmt.Exec(newEndTimestamp, username)
	pomobot.CheckErr(err)

	affected, err := res.RowsAffected()
	if affected > 0 {
		var addOrRemovedStr string
		if timeToAdd >= 0 {
			addOrRemovedStr = "Added"
		} else {
			addOrRemovedStr = "Removed"
		}
		msg := fmt.Sprintf("@%s %s %v minutes to your pomodoro", username, addOrRemovedStr, timeToAdd)
		err := bot.Send(consts.Channel, msg)
		pomobot.CheckErr(err)

	} else {
		msg := fmt.Sprintf("@%s Couldn't find any pomodoro session to end", username)
		err = bot.Send(consts.Channel, msg)
		pomobot.CheckErr(err)
	}
}

func updateSilentStatusPomoDb(silent bool, username string) {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	pomobot.CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("update pomodoros set silent=? where username=?")
	pomobot.CheckErr(err)

	_, err = stmt.Exec(silent, username)
	pomobot.CheckErr(err)
}

func getPomoFromDb(username string) Pomo {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	pomobot.CheckErr(err)

	defer db.Close()
	stmt := fmt.Sprintf("SELECT * FROM pomodoros WHERE username='%s'", username)
	rows, err := db.Query(
		stmt,
	)
	pomobot.CheckErr(err)

	var currentPomo Pomo
	fmt.Println(currentPomo.EndTimestamp)

	for rows.Next() {
		err := rows.Scan(
			&currentPomo.Username,
			&currentPomo.Task,
			&currentPomo.EndTimestamp,
			&currentPomo.Duration,
			&currentPomo.Silent,
		)
		pomobot.CheckErr(err)
	}
	return currentPomo
}

func insertIntoDb(pomo *Pomo, bot *twitch_api_wrapper.Bot) error {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	pomobot.CheckErr(err)
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO pomodoros (username, task, end_timestamp, duration, silent) VALUES ($1, $2, $3, $4, $5)",
		&pomo.Username, &pomo.Task, &pomo.EndTimestamp, &pomo.Duration, &pomo.Silent,
	)
	if err != nil {
		// there is already a pomo running for that user
		msg := fmt.Sprintf("@%s You already have a pomodoro session running", pomo.Username)
		sendErr := bot.Send(consts.Channel, msg)
		pomobot.CheckErr(sendErr)
		return err
	}
	return nil
}

func terminatePomo(username string, bot *twitch_api_wrapper.Bot) {
	db, err := sql.Open("sqlite3", "./twitch_bot.sqlite")
	pomobot.CheckErr(err)

	defer db.Close()
	statement, err := db.Prepare("DELETE FROM pomodoros WHERE username=?")
	pomobot.CheckErr(err)

	res, err := statement.Exec(username)
	pomobot.CheckErr(err)

	affected, err := res.RowsAffected()
	if affected > 0 {
		msg := fmt.Sprintf("@%s Ended running pomodoro session", username)
		err = bot.Send(consts.Channel, msg)
		pomobot.CheckErr(err)

	} else {
		msg := fmt.Sprintf("@%s Couldn't find any pomodoro session to end", username)
		err = bot.Send(consts.Channel, msg)
		pomobot.CheckErr(err)
	}
}
