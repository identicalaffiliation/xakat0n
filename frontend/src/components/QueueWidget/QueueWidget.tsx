import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Paper, Typography, Button } from '@mui/material';
import { useQueue } from '../../context/QueueContext';

export const QueueWidget: React.FC = () => {
  const navigate = useNavigate();
  const { state, startCheckout, leaveQueue } = useQueue();
  if (!state.status || state.status === 'PURCHASED' || state.status === 'CANCELLED' || state.status === 'EXPIRED') {
    return null;
  }

  const handleGoToQueue = () => {
    if (state.productId) {
      navigate(`/product/${state.productId}/queue`);
    }
  };

  const handleGoToCheckout = () => {
    if (state.productId) {
      startCheckout(state.productId);
      navigate(`/product/${state.productId}/queue`);
    }
  };

  const getStatusText = () => {
    switch (state.status) {
      case 'QUEUED':
        return `В очереди (место ${state.queuePosition || '?'})`;
      case 'OFFERED':
        return 'Товар доступен!';
      case 'CHECKOUT':
        return 'Оформление заказа';
      default:
        return 'Статус';
    }
  };

  const getActionButton = () => {
    if (state.status === 'QUEUED') {
      return (
        <Button variant="outlined" size="small" onClick={handleGoToQueue}>
          Перейти в очередь
        </Button>
      );
    }
    if (state.status === 'OFFERED') {
      return (
        <Button variant="contained" size="small" color="primary" onClick={handleGoToCheckout}>
          Перейти к оформлению
        </Button>
      );
    }
    if (state.status === 'CHECKOUT') {
      return (
        <Button variant="outlined" size="small" onClick={handleGoToQueue}>
          Продолжить оформление
        </Button>
      );
    }
    return null;
  };

  return (
    <Paper
      elevation={3}
      sx={{
        position: 'fixed',
        bottom: 20,
        right: 20,
        p: 2,
        zIndex: 9999,
        maxWidth: 280,
        borderRadius: 3,
        backgroundColor: '#fff',
        border: '1px solid #e0e0e0',
      }}
    >
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
        <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
          {getStatusText()}
        </Typography>
        {state.status === 'QUEUED' && state.timeLeft && (
          <Typography variant="caption" color="text.secondary">
            Осталось: ~{Math.floor(state.timeLeft / 60)} мин
          </Typography>
        )}
        <Box sx={{ display: 'flex', gap: 1, mt: 1, flexWrap: 'wrap' }}>
          {getActionButton()}
          {state.status === 'QUEUED' && (
            <Button
              variant="text"
              size="small"
              color="error"
              onClick={() => {
                if (state.productId) {
                  leaveQueue();
                  navigate('/products');
                }
              }}
            >
              Выйти
            </Button>
          )}
        </Box>
      </Box>
    </Paper>
  );
};