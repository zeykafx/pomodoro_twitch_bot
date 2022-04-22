import {
    Card,
    Grid,
    Tooltip,
    Typography,
    TextField,
    Button,
    Fade, Alert, Snackbar,
} from "@mui/material";
import React, {useEffect, useState} from "react";
import {IMaskInput} from "react-imask";
import "./Settings.css";

export declare const astilectron: any;

interface IMaskProps {
    onChange: (event: { target: { name: string; value: string } }) => void;
    name: string;
}

const TextMaskAuthToken = React.forwardRef<HTMLElement, IMaskProps>(
    (props, ref) => {
        const {onChange, ...other} = props;
        return (
            <IMaskInput
                {...other}
                //@ts-ignore
                mask={"o\\auth:000000000000000000000000000000"}
                definitions={{
                    "0": /[0-9, a-z, A-Z]/,
                }}
                lazy={true}
                // fuck ts errors, that shit makes no sense
                //@ts-ignore
                inputRef={ref}
                onAccept={(value: any) =>
                    onChange({target: {name: props.name, value}})
                }
                overwrite
            />
        );
    }
);


export default function Settings() {

    const [token, setToken] = useState<string>("");
    const [prefix, setPrefix] = useState<string>("");
    const [channel, setChannel] = useState<string>("");

    const [saveChanges, setSaveChanges] = useState<boolean>(false);
    const [successAlertOpen, setSuccessAlertOpen] = React.useState(false);
    const [errorAlertOpen, setErrorAlertOpen] = React.useState(false);

    useEffect(() => {
        if (saveChanges) {
            // send the new data once the user saves the changes
            let json_payload = JSON.stringify({token: token, prefix: prefix, channel: channel});

            // @ts-ignore
            astilectron.sendMessage(
                {name: "SET_SETTINGS", payload: json_payload},
                function (message: any) {
                    if (message["payload"] === "saved settings") {
                        // correctly saved the settings
                        setSuccessAlertOpen(true)
                    } else {
                        // couldn't save the settings
                        setErrorAlertOpen(true);
                    }
                }
            );
        }
    }, [saveChanges]);

    useEffect(() => {
        astilectron.sendMessage({name: "GET_SETTINGS", payload: ""}, (message: any) => {
            console.log(message["payload"])
            let jsonMessage = JSON.parse(message["payload"])
            console.log(jsonMessage)
            setToken(jsonMessage["token"])
            setPrefix(jsonMessage["prefix"])
            setChannel(jsonMessage["channel"])
        })
    }, []);


    let openExternalBrowser = (e: any) => {
        e.preventDefault()
        // @ts-ignore
        astilectron.sendMessage(
            {name: "URL", payload: e.target.href}, (message: any) => {
                console.log(message["payload"]);
            }
        )
    }

    return (
        <>
            {/* SUCCESS ALERT */}
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                open={successAlertOpen}
                autoHideDuration={10000}
                onClose={() => setSuccessAlertOpen(false)}
                sx={{
                    width: "75%",
                }}
            >
                <Alert
                    onClose={() => setSuccessAlertOpen(false)}
                    severity="info"
                    sx={{
                        width: "100%", fontSize: "42px"
                    }}
                >
                    <Typography variant={"h6"} component="div">
                        Successfully saved! Please restart the app.
                    </Typography>
                </Alert>
            </Snackbar>

            {/* ERROR ALERT */}
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                open={errorAlertOpen}
                autoHideDuration={10000}
                onClose={() => setErrorAlertOpen(false)}
                sx={{
                    width: "100%",
                }}
            >
                <Alert
                    onClose={() => setErrorAlertOpen(false)}
                    severity="error"
                    sx={{
                        width: "100%", fontSize: "42px"
                    }}
                >
                    <Typography variant={"h6"} component="div">
                        Error, couldn't save! Restart the app, and maybe re-download it. If that fails, contact Zeyka#8688 on discord.
                    </Typography>
                </Alert>
            </Snackbar>

            <Card variant={"outlined"} className={"settings-card"}>
                <Grid
                    container
                    direction="column"
                    justifyContent="center"
                    alignItems="flex-start"
                    spacing={2}
                >
                    {/* OAUTH TOKEN ------------------------*/}
                    <Grid item>
                        <Typography variant="h6">Auth token: </Typography>
                    </Grid>

                    <Grid item>
                        <Grid
                            container
                            direction="row"
                            justifyContent="flex-start"
                            alignItems="center"
                            spacing={2}
                        >
                            <Grid item>
                                <Tooltip
                                    title={"OAUTH Token used by the bot to connect to chat"}
                                    enterDelay={500}
                                >
                                    <TextField
                                        variant={"outlined"}
                                        id="oauth-token-input"
                                        label="Oauth token"
                                        value={token}
                                        InputProps={{
                                            inputComponent: TextMaskAuthToken as any,
                                        }}
                                        onChange={(newVal) => setToken(newVal.target.value)}
                                    />
                                </Tooltip>
                            </Grid>

                            <Grid item>
                                <Typography variant="body1">
                                    Enter the token that you got from{" "}
                                    <a onClick={openExternalBrowser} href={"https://twitchapps.com/tmi/"}>
                                        https://twitchapps.com/tmi/
                                    </a>
                                </Typography>
                            </Grid>
                        </Grid>
                    </Grid>

                    {/* PREFIX ------------------------ */}
                    <Grid item>
                        <Typography variant="h6">Command Prefix: </Typography>
                    </Grid>

                    <Grid item>
                        <Grid
                            container
                            direction="row"
                            justifyContent="flex-start"
                            alignItems="center"
                            spacing={2}
                        >
                            <Grid item>
                                <Tooltip
                                    title={"Command bot prefix, used to start all the commands"}
                                    enterDelay={500}
                                >
                                    <TextField
                                        variant={"outlined"}
                                        id="command-prefix-input"
                                        label="Command prefix"
                                        value={prefix}
                                        onChange={(newVal) => setPrefix(newVal.target.value)}
                                    />
                                </Tooltip>
                            </Grid>

                            <Grid item>
                                <Typography variant="body1">
                                    Enter the command prefix (e.g.: "!", "?",...)
                                </Typography>
                            </Grid>
                        </Grid>
                    </Grid>

                    {/* Channel ------------------------ */}
                    <Grid item>
                        <Typography variant="h6">Channel: </Typography>
                    </Grid>

                    <Grid item>
                        <Grid
                            container
                            direction="row"
                            justifyContent="flex-start"
                            alignItems="center"
                            spacing={2}
                        >
                            <Grid item>
                                <Tooltip title={"Twitch channel name"} enterDelay={500}>
                                    <TextField
                                        variant={"outlined"}
                                        id="channel-input"
                                        label="Channel name"
                                        value={channel}
                                        onChange={(newVal) => setChannel(newVal.target.value)}
                                    />
                                </Tooltip>
                            </Grid>

                            <Grid item>
                                <Typography variant="body1">
                                    Enter your twitch channel name, <b>MUST BE LOWER CASE!</b>
                                </Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Card>

            <Fade in={token !== "" || prefix !== "" || channel !== ""} className={"save-settings-button"}>
                <Grid item>
                    <Tooltip title={"Save the current changes"} enterDelay={500}>
                        <Button
                            variant={"contained"}
                            onClick={() => {
                                setSaveChanges(true);
                            }}
                        >
                            Save Changes
                        </Button>
                    </Tooltip>
                </Grid>
            </Fade>
        </>
    );
}
