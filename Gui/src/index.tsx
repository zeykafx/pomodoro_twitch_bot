import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {createTheme, CssBaseline, darkScrollbar, ThemeProvider} from "@mui/material";

const themeMode = 'dark';

const theme = createTheme({
    palette: {
        mode: themeMode,
        primary: {
            main: "#2C3E50",
        },
        secondary: {
            main: "#34495E",
        },
    },
});

// @ts-ignore
theme.components.MuiCssBaseline ={
    styleOverrides: {
        body: theme.palette.mode === 'dark' ? darkScrollbar() : null,
    },
}

ReactDOM.render(
    <ThemeProvider theme={theme}>
        <App />
        <CssBaseline />
    </ThemeProvider>,
    document.getElementById("root")
);
