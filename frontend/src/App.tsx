import { ThemeProvider } from '@mui/material/styles';
import theme from './theme';
import AuthPage from './components/AuthPage/AuthPage';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <AuthPage />
    </ThemeProvider>
  );
}

export default App;

