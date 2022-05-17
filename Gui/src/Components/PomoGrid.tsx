import * as React from 'react';
import {DataGrid, GridActionsCellItem, GridColumns, GridRowId} from '@mui/x-data-grid';
import {Pomo} from "../Views/App";
import {Alert, Snackbar, Typography} from "@mui/material";
import DeleteIcon from '@mui/icons-material/DeleteOutlined';

export declare const astilectron: any;

interface IPomoGridProps {
    pomos: Array<Pomo>,
    setPomos: any,
}

export default function PomoGrid(props: IPomoGridProps) {
    const [successAlertOpen, setSuccessAlertOpen] = React.useState(false);
    const [errorAlertOpen, setErrorAlertOpen] = React.useState(false);

    const columns: GridColumns = [
        {field: 'id', headerName: "ID", hide: true},
        {field: 'username', headerName: 'Username', width: 130},
        {
            field: 'task', headerName: 'Task - editable', width: 130, flex: 1, editable: true, valueSetter: params => {
                editPomoTaskFromUsername(params.row.username, params.value, params.row.id);
                return {...params.row, task: params.value};
            }
        },
        {field: 'pomoDuration', headerName: 'Duration', width: 70},
        {field: 'silent', headerName: 'Silent', width: 70},
        {field: 'time_left', headerName: 'Time Left', width: 100},
        {
            field: "actions",
            type: 'actions',
            headerName: "Actions",
            width: 70,
            cellClassName: 'actions',
            getActions: ({id}) => {
                return [
                    <GridActionsCellItem
                        icon={<DeleteIcon/>}
                        label="Delete"
                        onClick={handleDeleteClick(id)}
                        color="inherit"
                        showInMenu={false}/>
                ]
            }
        }
    ];


    let editPomoTaskFromUsername = (username: string, newTask: string, id: number): void => {
        let jsonRepr = JSON.stringify({username: username, new_task: newTask})
        astilectron.sendMessage(
            {name: "EDIT_TASK", payload: jsonRepr},
            function (message: any) {
                if (message["payload"] === "OK") {
                    setSuccessAlertOpen(true);
                } else {
                    setErrorAlertOpen(true);
                }
            }
        );
    }

    const handleDeleteClick = (id: GridRowId) => () => {
        // @ts-ignore
        let pomoToDelete: Pomo = props.pomos[id];
        props.setPomos(props.pomos.filter((pomo) => pomo.id !== id));
        astilectron.sendMessage(
            {name: "DELETE_POMO", payload: pomoToDelete.username},
            function (message: any) {
                if (message["payload"] === "OK") {
                    setSuccessAlertOpen(true);
                } else {
                    setErrorAlertOpen(true);
                }
            }
        );
    };

    return (
        <>
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                open={successAlertOpen}
                autoHideDuration={1000}
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
                        Success!
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
                        Error, couldn't edit/delete that user's pomo, its most likely a bot error, try again or restart the app if this persists.
                    </Typography>
                </Alert>
            </Snackbar>

            <div style={{height: 330, width: '100%'}}>
                <DataGrid
                    rows={props.pomos}
                    columns={columns}
                    pageSize={4}
                    rowsPerPageOptions={[4]}
                />
            </div>
        </>
    );
}
