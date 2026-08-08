import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Paper, TextField, Button, Typography, Box } from '@mui/material';
import { login } from '../../api/auth';

const AuthPage = () => {
  const [username, setUsername] = useState('');
  const navigate = useNavigate();

  const handleLogin = async () => {
    if (!username.trim()) return alert('Введите имя');
    try {
      const data = await login(username);
      sessionStorage.setItem('sessionToken', data.token);
      navigate('/products');
    } catch (error) {
      console.error(error);
      alert('Ошибка входа');
    }
  };

  return (
    <Box sx={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: '#fff' }}>
      <Paper elevation={0} sx={{ p: 4, width: '100%', maxWidth: 400, borderRadius: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', mb: 3 }}>
          <img src="/avito-logo.svg" alt="Avito" style={{ width: 32, height: 32, marginRight: 8 }} />
          <Typography variant="h4" sx={{ fontWeight: 700, color: '#000' }}>Avito</Typography>
        </Box>
        <TextField
          fullWidth
          label="Имя"
          variant="outlined"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          sx={{ mb: 2, '& .MuiInputBase-root': { height: 48 }, '& .MuiOutlinedInput-root': { borderRadius: 2 } }}
        />
        <Button
          fullWidth
          variant="contained"
          sx={{ bgcolor: '#00AAFF', color: '#fff', height: 48, borderRadius: 2, textTransform: 'none', fontWeight: 600, '&:hover': { bgcolor: '#0088cc' } }}
          onClick={handleLogin}
        >
          Войти
        </Button>
        <Typography variant="body2" color="textSecondary" sx={{ textAlign: 'center', mt: 2 }}>
          Войдите, чтобы продолжить покупки
        </Typography>
      </Paper>
    </Box>
  );
};

export default AuthPage;