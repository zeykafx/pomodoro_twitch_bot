# Pomodoro Twitch Bot

This is a simple pomodoro bot that outputs the running pomos to a file, making it easy to add as an overlay in OBS.


# Commands:
When running the bot locally, it will read chat and respond to commands. It also writes all the running pomos to a text file which you have to import into OBS if you want to have a pomoboard overlay, the file is refreshed every 5 seconds.


- `[prefix]pomo [time] [task]`: starts a pomo for `[time]`.
- `[prefix]pomo end/cancel/stop/finish`: ends the user's running pomo.
- `[prefix]pomo check`: gives the remaining time for the user's pomo.
- `[prefix]pomo chat/silent/silence/mod`: stops the bot from discouraging the user to focus.
- `[prefix]pomo add/plus OR remove/minus [time]`: adds or removes time to the user's pomo.
When you see a `/` here it means that either works and do the same thing.


# Instructions
You will need a second twitch account, this will be your bot. Head over to https://twitchapps.com/tmi/, connect with the bot account and copy the token you see, it should look like `oauth:xxxxxxxxxxxx`.

Next, run the executable, and complete the first time setup.

Your bot should be working, enjoy!
