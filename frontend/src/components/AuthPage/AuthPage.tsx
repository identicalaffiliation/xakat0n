import { useState } from 'react';
import { Paper, TextField, Button, Typography, Box } from '@mui/material';

const AuthPage = () => {
  const [phone, setPhone] = useState('');

  const handleLogin = () => {
    if (!phone.trim()) {
      alert('Введите номер телефона');
      return;
    }
    window.location.href = '/products';
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#f5f5f5',
      }}
    >
      <Paper
        elevation={0}
        sx={{
          p: 4,
          width: '100%',
          maxWidth: 400,
          borderRadius: 2,
          backgroundColor: '#fff',
          border: 'none',
          boxShadow: 'none',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', mb: 3 }}>
          <img
            src="./avito-logo.svg"
            alt="Авито"
            style={{ width: 32, height: 32, marginRight: 8 }}
          />
          <Typography variant="h4" sx={{ fontWeight: 700, color: '#000000' }}>
            Авито
          </Typography>
        </Box>

        <TextField
          fullWidth
          label="Телефон или email"
          variant="outlined"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          sx={{
            mb: 2,
            '& .MuiInputBase-root': { height: 48 },
            '& .MuiOutlinedInput-root': { borderRadius: 2 },
          }}
        />
        <Button
          fullWidth
          variant="contained"
          sx={{
            backgroundColor: '#00A2E3',
            color: '#fff',
            height: 48,
            fontSize: 16,
            fontWeight: 600,
            textTransform: 'none',
            borderRadius: 2,
            '&:hover': { backgroundColor: '#0088c7' },
          }}
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