package setup

import (
	"fmt"
	"os"
)

// FirstTimeSetup is not used anymore, the setup is done in the gui
func FirstTimeSetup() {
	var token, prefix, channel string

	fmt.Println("----- Pomoboard bot first time setup, Welcome! ----- ")
	fmt.Println("Thanks for installing this bot, to complete the setup you will need to enter 3 things:")
	fmt.Println("\t - 1: The bot token, this is can be found if you login with the bot at this address: https://twitchapps.com/tmi/ . The token should look like 'oauth:xxxxxxxxxxxxxxxxxxx'")
	fmt.Println("\t - 2: The bot prefix, you can choose any command prefix such as '!' or '?'")
	fmt.Println("\t - 3: Your twitch channel, just enter the name like 'zeykafx' or 'hasanabi',...")

	fmt.Println("Enter the bot's token (should be something like 'oauth:xxxxxxxxxxxxxxxxx'):")
	_, err := fmt.Scanln(&token)
	if err != nil {
		panic(err)
	}

	fmt.Println("Enter the bot's prefix:")
	_, err = fmt.Scanln(&prefix)
	if err != nil {
		panic(err)
	}

	fmt.Println("Enter your Twitch channel:")
	_, err = fmt.Scanln(&channel)
	if err != nil {
		panic(err)
	}

	err = WriteSettingsToFile(token, prefix, channel)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setup complete! The bot will begin writing the running pomos to \"pomoboard.txt\" which you have to import in OBS.")
	fmt.Println("Thanks for installing, enjoy! - Zeyka.")
}

// WriteSettingsToFile is used to set the settings by the gui
func WriteSettingsToFile(token string, prefix string, channel string) error {
	file, err := os.Create(".env")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString(fmt.Sprintf("TOKEN=\"%s\"\n", token))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("PREFIX=\"%s\"\n", prefix))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("CHANNEL=\"%s\"\n", channel))
	if err != nil {
		return err
	}

	return nil
}
