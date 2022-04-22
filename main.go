package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
	"os"
	"pomodoro_twitch_bot/consts"
	"pomodoro_twitch_bot/help_commands"
	"pomodoro_twitch_bot/pomobot"
	"pomodoro_twitch_bot/pomodoro"
	"time"
)

func initCommands() {
	// add the functions that handle commands here
	pomobot.Accept(pomodoro.HandlePomoCommand, "pomo")
	pomobot.Accept(help_commands.Help, "help")
}

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	fs    = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	debug = fs.Bool("d", false, "enables the debug mode")
	w     *astilectron.Window
)

func main() {
	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Parse flags
	fs.Parse(os.Args[1:])

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:         Asset, // ignore the errors, this is injected by the bindings (which you get by running astilectron-bundler)
		AssetDir:      AssetDir,
		RestoreAssets: RestoreAssets,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  *debug,
		Logger: l,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("Options"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Open dev tools"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						fmt.Println("Opened dev tools")
						err := w.OpenDevTools()
						if err != nil {
							return false
						}
						return
					},
				},
				{
					Label: astikit.StrPtr("Reload"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						fmt.Println("Reload")
						err := w.ExecuteJavaScript("location.reload()")
						if err != nil {
							return false
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]

			// --- START THE BOT ---
			go func() {
				consts.LoadPrefix(w)
				initCommands()
				pomobot.InitBot()
			}()

			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "connect", "connect"); err != nil {
					l.Println(fmt.Errorf("sending connect event failed: %w", err))
				}
			}()
			return nil
		},
		Windows: []*bootstrap.Window{{
			Homepage: "index.html", // this uses the built React app in ./resources/app/
			//Homepage:       "http://localhost:3000/", // this when running the React app in dev mode
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(600),
				Width:           astikit.IntPtr(800),
				AutoHideMenuBar: astikit.BoolPtr(true),
				Title:           astikit.StrPtr("Pomodoro Bot"),
			},
		}},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}

}
