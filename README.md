# Pomodoro Twitch Bot

This is a simple pomodoro bot that outputs the running pomos to a file which can be imported as an overlay.

# Instructions - Important
Download the archive from the releases on the right, unzip it, and run twitch_bot.exe.

You will need a second twitch account, this will be your bot. Head over to https://twitchapps.com/tmi/, connect with the bot account and copy the token you see, it should look like `oauth:xxxxxxxxxxxxxxxxxxxxxxxxxxxxx`.

Next, head over to the settings by clicking on the cog wheel at the top right of the app, in there, fill in the oauth token you just copied, as well as the prefix for the bot and your twitch channel (make sure to enter the channel name as it is found in the url, i.e. without any capital letters or weird symbols).

If you want to turn on the pomo board, make sure to press the "start board" button.

To add the board in OBS, create a new text source, then tick "read from file" and choose the .txt file created by the bot.

<img src="images\text_source.png" width="60%"/>
<img src="images\text_source_2.png" width="70%"/>

# Commands:
The bot will read chat and respond to commands. It also writes all the running pomos to a text file, the file is refreshed every 5 seconds.

- `[prefix]pomo [time] [task]`: starts a pomo for `[time]` minutes.
- `[prefix]pomo end/cancel/stop/finish`: ends the user's running pomo.
- `[prefix]pomo check`: gives the remaining time for the user's pomo.
- `[prefix]pomo chat/silent/silence/mod`: stops the bot from encouraging the user to focus.
- `[prefix]pomo add/plus/+ OR remove/minus/- [time]`: adds or removes time to the user's pomo.

When you see a `/` here it means that either words works and do the same thing, you can set `[prefix]` to whatever you want from the settings.

# Credits
This bot was made by ZeykaFX.
It was heavily inspired by [RumpleStudy](https://www.twitch.tv/rumplestudy)'s private pomo bot.

# Images

<img src="images\main_with_pomo.png" width="60%"/>
