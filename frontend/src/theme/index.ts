import { createTheme } from '@mui/material/styles';

const theme = createTheme({
    palette: {
        primary: {
            main: '#000000', // Preto
        },
        secondary: {
            main: '#D4AF37', // Dourado
        },
        background: {
            default: '#FFFFFF', // Branco
        },
    },
    typography: {
        fontFamily: [
            'Roboto',
            'Arial',
            'sans-serif',
        ].join(','),
    },
});

export default theme;