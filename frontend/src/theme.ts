import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  shape: {
    borderRadius: 6,
  },
  typography: {
    fontFamily: [
      'Apple System',
      'Segoe UI',
      'Verdana',
      'Arial',
      'sans-serif',
    ].join(','),
  },
  components: {
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 6,
          },
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 6,
          textTransform: 'none',
          fontWeight: 600,
        },
      },
    },
  },
});

export default theme;