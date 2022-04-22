import React, {useEffect, useState} from "react";
import "./App.css";
import {
    Alert,
    Button,
    Card,
    Grid,
    IconButton, Snackbar,
    Tooltip,
    Typography,
} from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import DashboardIcon from '@mui/icons-material/Dashboard';
import DashBoard from "./DashBoard";
import Settings from "./Settings";

export declare const astilectron: any;
let astielectron_ready = false;

function App() {

    // This will wait for the astilectron namespace to be ready
    document.addEventListener("astilectron-ready", function () {
        // This will listen to messages sent by GO
        astilectron.onMessage(function (message: any) {
            console.log(message)
            if (message["name"] === "connected") {
                console.log("Connected to go backend!");
                return "Connected";
            }
            else if (message["name"] === "NO SETTINGS") {
                setSettingsAlert(true);
                return "showed alert";
            }

        });
        astielectron_ready = true;
    });

    const [botStatus, setBotStatus] = useState<boolean>(false); // false == off, true == on
    const [pomos, setPomos] = useState<string>("");
    const [pomoNbr, setPomosNbr] = useState<number>(0);
    const [boardStarted, setBoardStarted] = useState<boolean>(false);

    const [settingsAlert, setSettingsAlert] = useState<boolean>(false);

    const [currentPage, setCurrentPage] = useState<number>(0); // 0 = dashboard, 1 = settings

    // get pomo info

    const getBotStatus = () => {
        astilectron.sendMessage(
            {name: "STATUS", payload: "STATUS"},
            function (message: any) {
                if (message["payload"] === "on") {
                    setBotStatus(true);
                } else {
                    setBotStatus(false);
                }
            }
        );
    };

    const getPomos = () => {
        astilectron.sendMessage(
            {name: "RUNNING_POMOS", payload: "RUNNING_POMOS"},
            function (message: any) {
                setPomos(message["payload"]);
                setPomosNbr(message["payload"].split("\n").length - 1);
            }
        );
    };

    // start and stop writing to file ------

    const startWritingFile = () => {
        astilectron.sendMessage(
            {name: "START_FILE", payload: "start writing file "},
            function (message: any) {
                console.log(message["payload"]);
                if (message["payload"] === "Started writing to file!") {
                    setBoardStarted(true);
                }
            }
        );
    };

    const stopWritingFile = () => {
        astilectron.sendMessage(
            {name: "STOP_FILE", payload: "stop writing file "},
            function (message: any) {
                console.log(message["payload"]);
                if (message["payload"] === "Stopped writing to file!") {
                    setBoardStarted(false);
                }
            }
        );
    };

    useEffect(() => {
        // get the status on page load
        if (astielectron_ready) {
            getBotStatus();
            getPomos();
        }

        // then refresh the data every 5 seconds
        setInterval(() => {
            getBotStatus();
            getPomos();
        }, 5000);
    }, []);

    return (
        <>
            {/* ERROR ALERT */}
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                open={settingsAlert}
                autoHideDuration={100000} // 100 seconds to make sure the users sees this
                onClose={() => setSettingsAlert(false)}
                sx={{
                    width: "100%",
                }}
            >
                <Alert
                    onClose={() => setSettingsAlert(false)}
                    severity="error"
                    sx={{
                        width: "100%", fontSize: "42px"
                    }}
                >
                    <Typography variant={"h6"} component="div">
                        You don't seem to have any settings configured, head to the settings (cog wheel top right) to setup the bot.
                    </Typography>
                </Alert>
            </Snackbar>

            <div className="App">
                <div className={"Header"}>
                    <Grid
                        container
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                    >
                        <Grid item>
                            <Typography
                                variant="h6"
                                component="div"
                                sx={{flexGrow: 1, color: "white", cursor: "pointer"}}
                                onClick={() => setCurrentPage(currentPage === 0 ? 1 : 0)}
                            >
                                Pomodoro Twitch Bot
                            </Typography>
                        </Grid>
                        <Grid item>
                            <Tooltip title={currentPage === 0 ? "Settings" : "Dashboard"} enterDelay={500}>
                                <IconButton aria-label="Settings" onClick={() => {
                                    setCurrentPage(currentPage === 0 ? 1 : 0)
                                }}>
                                    {currentPage === 0 ? <SettingsIcon/> : <DashboardIcon/>}
                                </IconButton>
                            </Tooltip>
                        </Grid>
                    </Grid>
                </div>

                {currentPage === 0 ? (
                    <DashBoard
                        botStatus={botStatus}
                        boardStarted={boardStarted}
                        pomoNbr={pomoNbr}
                        pomos={pomos}
                        startWritingFile={startWritingFile}
                        stopWritingFile={stopWritingFile}
                    />
                ) : (
                    <Settings />
                )}
            </div>
        </>

    );
}

export default App;
