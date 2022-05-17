import React, {useState} from "react";
import "./App.css";
import {
    Alert,
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
let astielectronReady = false;

export interface Pomo {
    id: number,
    username: string,
    task: string,
    end_timestamp: string,
    pomoDuration: number,
    silent: boolean,
    time_left: number
}

function App() {

    // This will wait for the astilectron namespace to be ready
    document.addEventListener("astilectron-ready", function () {
        // This will listen to messages sent by GO
        astilectron.onMessage(function (message: any) {
            console.log(message)
            if (message["name"] === "connected") {
                console.log("Connected to go backend!");
                return "Connected";
            } else if (message["name"] === "NO SETTINGS") {
                setSettingsAlert(true);
                return "showed alert";
            }

        });
        astielectronReady = true;
    });


    const [settingsAlert, setSettingsAlert] = useState<boolean>(false);
    const [currentPage, setCurrentPage] = useState<number>(0); // 0 = dashboard, 1 = settings


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
                        width: "80%", fontSize: "42px"
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
                        astielectronReady={astielectronReady}
                    />
                ) : (
                    <Settings/>
                )}
            </div>
        </>

    );
}

export default App;
