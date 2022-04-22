# Pomodoro Twitch Bot

This is a simple pomodoro bot that outputs the running pomos to a file, making it easy to add as an overlay in OBS.


# Instructions - READ THIS
Download the archive from the releases on the right, unzip it, and run twitch_bot.exe.

You will need a second twitch account, this will be your bot. Head over to https://twitchapps.com/tmi/, connect with the bot account and copy the token you see, it should look like `oauth:xxxxxxxxxxxx`.

Next, run the executable, and fill in the settings (click on the cog wheel, top right of the gui).

If you want to turn on the pomo board, make sure to press the "start board" button.

Your bot should be working, enjoy!

# Commands:
When running the bot locally, it will read chat and respond to commands. It also writes all the running pomos to a text file which you have to import into OBS if you want to have a pomoboard overlay, the file is refreshed every 5 seconds.


- `[prefix]pomo [time] [task]`: starts a pomo for `[time]`.
- `[prefix]pomo end/cancel/stop/finish`: ends the user's running pomo.
- `[prefix]pomo check`: gives the remaining time for the user's pomo.
- `[prefix]pomo chat/silent/silence/mod`: stops the bot from discouraging the user to focus.
- `[prefix]pomo add/plus OR remove/minus [time]`: adds or removes time to the user's pomo.
When you see a `/` here it means that either works and do the same thing.

# Images

<img src="images\main_with_pomo.png" width="60%"/>

# Credits
This bot was made by ZeykaFX, originally made to be used by [ellenbearli](https://www.twitch.tv/ellenbearli). 

It was heavily inspired by [RumpleStudy](https://www.twitch.tv/rumplestudy)'s private pomo bot.
