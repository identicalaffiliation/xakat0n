import { createTheme } from '@mui/material/styles';

const colors = {
  primary: '#00AAFF',
  secondary: '#FF6163',
  success: '#97CF26',
  purple: '#A169F7',
  black: '#000000',
  white: '#FFFFFF',
};

const theme = createTheme({
  shape: {
    borderRadius: 12,
  },
  palette: {
    primary: {
      main: colors.primary,
    },
    secondary: {
      main: colors.secondary,
    },
    success: {
      main: colors.success,
    },
    info: {
      main: colors.purple,
    },
    background: {
      default: colors.white,
      paper: colors.white,
    },
    text: {
      primary: colors.black,
    },
  },
  components: {
    MuiButton: {
      variants: [
        {
          props: { variant: 'contained', color: 'primary' },
          style: {
            backgroundColor: colors.primary,
            color: colors.white,
            '&:hover': {
              backgroundColor: '#0088cc',
            },
          },
        },
        {
          props: { variant: 'contained', color: 'secondary' },
          style: {
            backgroundColor: colors.secondary,
            color: colors.white,
          },
        },
        {
          props: { variant: 'contained', color: 'success' },
          style: {
            backgroundColor: colors.success,
            color: colors.white,
          },
        },
        {
          props: { variant: 'contained', color: 'info' },
          style: {
            backgroundColor: colors.purple,
            color: colors.white,
          },
        },
        {
          props: { variant: 'outlined', color: 'primary' },
          style: {
            borderColor: colors.primary,
            color: colors.primary,
          },
        },
        {
          props: { variant: 'outlined', color: 'secondary' },
          style: {
            borderColor: colors.secondary,
            color: colors.secondary,
          },
        },
        {
          props: { variant: 'text', color: 'primary' },
          style: {
            color: colors.primary,
          },
        },
      ],
      styleOverrides: {
        root: {
          borderRadius: 12,
          textTransform: 'none',
          fontWeight: 600,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 12,
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          borderRadius: 0,
        },
      },
    },
  },
});

export default theme;