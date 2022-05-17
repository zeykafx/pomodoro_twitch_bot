import {Button, Card, Grid, Tooltip, Typography} from '@mui/material'
import PomoGrid from "../Components/PomoGrid";
import {useEffect, useState} from "react";
import {Pomo} from "./App";

export declare const astilectron: any;

interface IDashBoardProps {
    astielectronReady: boolean
}

export default function DashBoard(props: IDashBoardProps) {
    const [botStatus, setBotStatus] = useState<boolean>(false); // false == off, true == on
    const [pomos, setPomos] = useState<Array<Pomo>>([]);
    const [pomoNbr, setPomosNbr] = useState<number>(0);
    const [boardStarted, setBoardStarted] = useState<boolean>(false);

    useEffect(() => {
        // get the status on page load
        if (props.astielectronReady) {
            getBotStatus();
            getPomos();
        }

        // then refresh the data every 5 seconds
        setInterval(() => {
            getBotStatus();
            getPomos();
        }, 5000);
    }, []);


    // get pomo info ----------
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
                let jsonPomos = JSON.parse(message["payload"])
                if (jsonPomos !== null) {
                    jsonPomos.map((pomo: Pomo) => {
                        pomo.time_left = Math.round(pomo.time_left)
                    })
                    setPomos(jsonPomos);
                    setPomosNbr(jsonPomos.length);
                } else {
                    setPomos([])
                    setPomosNbr(0);
                }
            }
        );
    };

    // start and stop writing to file ------
    const startWritingFile = () => {
        astilectron.sendMessage(
            {name: "START_FILE", payload: "start writing file "},
            function (message: any) {
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
                if (message["payload"] === "Stopped writing to file!") {
                    setBoardStarted(false);
                }
            }
        );
    };

    return (
        <div className={"MainContent"}>
            <Card variant="outlined" className="info-card">
                <Grid
                    container
                    spacing={1}
                    direction="row"
                    justifyContent="space-around"
                    alignItems="center"
                >
                    <Grid item xs>
                        <Tooltip title={"If 'ON' then the bot is connected to the chat"} enterDelay={500}>
                            <Typography variant="h6">
                                Bot status: <span className={botStatus ? "ON" : "OFF"}>{botStatus ? "ON" : "OFF"}</span>
                            </Typography>
                        </Tooltip>
                    </Grid>

                    <Grid item xs>
                        <Tooltip title={"If 'ON' then the bot is writing the text file with all the pomos"} enterDelay={500}>
                            <Typography variant={"h6"}>
                                Board Status: <span className={boardStarted ? "ON" : "OFF"}>{boardStarted ? "ON" : "OFF"}</span>
                            </Typography>
                        </Tooltip>
                    </Grid>

                    <Grid item xs>
                        <Tooltip title={"Number of active pomos"} enterDelay={500}>
                            <Typography variant={"h6"}>
                                Number of pomos: {pomoNbr}
                            </Typography>
                        </Tooltip>
                    </Grid>

                </Grid>
            </Card>

            <Card variant="outlined" className="info-card">
                <Typography variant={"h6"}>Pomos running:</Typography>
                <div className={"pomo-list"}>
                    {pomos.length > 0 ?
                        <PomoGrid pomos={pomos} setPomos={setPomos}/>
                        : "No pomos running"
                    }

                </div>
            </Card>

            <Card variant="outlined" className="info-card">
                <Typography variant={"h6"}>Controls:</Typography>
                <Grid
                    container
                    spacing={1}
                    direction="row"
                    justifyContent="space-around"
                    alignItems="center"
                >
                    <Grid item xs>
                        <Tooltip title={"Click to start or stop the board"} enterDelay={500}>
                            <Button
                                variant={"contained"}
                                onClick={() => {
                                    !boardStarted ? startWritingFile() : stopWritingFile()
                                }}
                            >
                                {boardStarted ? "Stop the board" : "Start the board"}
                            </Button>
                        </Tooltip>
                    </Grid>

                </Grid>
            </Card>

        </div>
    )
}
