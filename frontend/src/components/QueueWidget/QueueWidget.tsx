

import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Box, Paper, Typography, Button } from '@mui/material';
import { useQueue } from '../../context/QueueContext';

export const QueueWidget: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { state, startCheckout, leaveQueue } = useQueue();
  const [timer, setTimer] = useState<number | null>(null);
  useEffect(() => {
    if (!state.status || state.timeLeft === null) {
      setTimer(null);
      return;
    }

    setTimer(state.timeLeft);
    const interval = setInterval(() => {
      setTimer((prev) => {
        if (prev === null || prev <= 1) {
          clearInterval(interval);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [state.status, state.timeLeft]);
  const isQueuePage = location.pathname.includes('/queue');
  const isCheckoutPage = location.pathname.includes('/checkout');
  if (isQueuePage || isCheckoutPage) {
    return null;
  }
  if (!state.status || state.status === 'PURCHASED' || state.status === 'CANCELLED' || state.status === 'EXPIRED') {
    return null;
  }
  if (location.pathname.match(/^\/product\/\d+\/queue$/)) {
    return null;
  }

  const handleGoToQueue = () => {
    if (state.itemId) {
      navigate(`/product/${state.itemId}/queue`);
    }
  };

  const handleGoToCheckout = () => {
    if (state.itemId) {
      startCheckout(state.itemId);
      navigate(`/product/${state.itemId}/queue`);
    }
  };

  const getStatusText = () => {
    switch (state.status) {
      case 'QUEUED':
        return `В очереди`;
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

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
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
        {timer !== null && (
          <Typography variant="caption" color="text.secondary">
            Осталось: {formatTime(timer)}
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
                if (state.itemId) {
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