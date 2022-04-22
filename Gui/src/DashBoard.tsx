import {Button, Card, Grid, Tooltip, Typography} from '@mui/material'
import React, {Component} from 'react'

interface IDashBoardProps {
    botStatus: boolean,
    boardStarted: boolean,
    pomoNbr: number,
    pomos: string,
    startWritingFile: () => void,
    stopWritingFile: () => void,
}

export default function DashBoard(props: IDashBoardProps) {
    let {botStatus, boardStarted, pomoNbr, pomos, startWritingFile, stopWritingFile} = props;

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
                        <Tooltip title={"Number of active pomodoros"} enterDelay={500}>
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
                        pomos.split("\n").map((pomo) =>
                            <Typography variant={"body1"} key={pomo}>{pomo}</Typography>
                        )
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
