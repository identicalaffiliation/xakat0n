import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Paper,
  Container,
  CircularProgress,
  Divider,
  Grid,
  TextField,
  Chip,
  AppBar,
  Toolbar,
  IconButton,
} from '@mui/material';
import { Person, FavoriteBorder } from '@mui/icons-material';
import { useQueue } from '../../context/QueueContext';
import { products } from '../../data/products';

const QueuePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { state, leaveQueue, expireOffer } = useQueue();
  const product = products.find((p) => p.id === Number(id));
  const [timer, setTimer] = useState<number | null>(null);
  useEffect(() => {
    if (state.status === 'CHECKOUT' || state.status === 'OFFERED') {
      if (state.timeLeft === null) return;
      setTimer(state.timeLeft);
      const interval = setInterval(() => {
        setTimer((prev) => {
          if (prev === null || prev <= 1) {
            clearInterval(interval);
            if (state.status === 'OFFERED' || state.status === 'CHECKOUT') {
              expireOffer();
            }
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [state.status, state.timeLeft, expireOffer]);

  if (!product) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Typography variant="h5">Товар не найден</Typography>
        <Button variant="outlined" onClick={() => navigate('/products')}>
          Вернуться к каталогу
        </Button>
      </Container>
    );
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };
  const renderContent = () => {
    if (!state.status) {
      return (
        <Box>
          <Typography variant="h5">Статус не определён</Typography>
          <Typography variant="body1" color="text.secondary">
            Попробуйте вернуться к товару и нажать «Купить» снова.
          </Typography>
          <Button variant="contained" sx={{ mt: 2 }} onClick={() => navigate(`/product/${product.id}`)}>
            Назад к товару
          </Button>
        </Box>
      );
    }

    switch (state.status) {
      case 'CHECKOUT':
        return (
          <Box>
            {/* Верхняя панель (закреплённая) */}
            <AppBar position="sticky" color="default" elevation={0} sx={{ backgroundColor: '#fff', top: 0, zIndex: 1100 }}>
              <Toolbar sx={{ justifyContent: 'space-between', py: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', ml: -20 }}>
                  <img
                    src="/avito-logo.svg"
                    alt="Avito"
                    style={{ width: 48, height: 48, marginRight: 12 }}
                  />
                  <Typography variant="h5" sx={{ fontWeight: 700, color: '#000000', fontSize: '1.8rem'}}>
                    Avito
                  </Typography>
                </Box>
                <Box>
                  <IconButton color="inherit">
                    <FavoriteBorder />
                  </IconButton>
                  <IconButton color="inherit">
                    <Person />
                  </IconButton>
                </Box>
              </Toolbar>
            </AppBar>
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                Оформление заказа
              </Typography>

              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  {/* Карточка товара */}
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box
                        component="img"
                        src={product.image}
                        alt={product.title}
                        sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }}
                      />
                      <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                        <Typography variant="body2" color="text.secondary">
                          Категория: {product.category}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          В наличии: {product.stock} шт.
                        </Typography>
                        {product.is_limited && (
                          <Chip
                            label="Лимитированный"
                            size="small"
                            sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }}
                          />
                        )}
                      </Box>
                    </Box>
                  </Paper>

                  {/* Дополнительная плашка с описанием */}
                  <Paper elevation={0} sx={{ p: 2, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="body2">
                      <strong>Описание:</strong> {product.title} – отличный выбор. Состояние идеальное, полный комплект.
                    </Typography>
                  </Paper>

                  {/* Поля ФИО и телефон */}
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>
                      Данные получателя
                    </Typography>
                    <TextField
                      fullWidth
                      label="ФИО получателя"
                      variant="outlined"
                      size="small"
                      value="Иванов Иван Иванович"
                      sx={{ mb: 2 }}
                      slotProps={{ input: { readOnly: true } }}
                    />
                    <TextField
                      fullWidth
                      label="Контактный телефон"
                      variant="outlined"
                      size="small"
                      value="+7 (999) 123-45-67"
                      slotProps={{ input: { readOnly: true } }}
                    />
                  </Paper>

                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => navigate(`/product/${product.id}`)}
                    >
                    Назад
                    </Button>
                </Grid>

                {/* Правая колонка */}
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>
                      Детали заказа
                    </Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                      <Typography>Товар</Typography>
                      <Typography sx={{ fontWeight: 500 }}>{product.title}</Typography>
                    </Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                      <Typography>Количество</Typography>
                      <Typography sx={{ fontWeight: 500 }}>1 шт.</Typography>
                    </Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                      <Typography>Стоимость</Typography>
                      <Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>
                        {product.price.toLocaleString()} ₽
                      </Typography>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    <Box sx={{ textAlign: 'center', mb: 2 }}>
                      <Typography variant="body2" color="text.secondary">
                        Осталось времени для оплаты
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 2, mt: 1 }}>
                        <CircularProgress
                          variant="determinate"
                          value={((timer || 0) / 120) * 100}
                          size={64}
                        />
                        <Typography variant="h2" sx={{ fontWeight: 700, fontSize: '2.5rem' }}>
                          {formatTime(timer || 0)}
                        </Typography>
                      </Box>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    {/* Промокод и скидка в виде карточек */}
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, backgroundColor: '#fff' }}>
                        <Typography variant="caption" color="text.secondary">Промокод</Typography>
                        <Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography>
                      </Paper>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, backgroundColor: '#fff' }}>
                        <Typography variant="caption" color="text.secondary">Скидка по карте</Typography>
                        <Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography>
                      </Paper>
                    </Box>

                    <Button
                      variant="contained"
                      fullWidth
                      sx={{ backgroundColor: '#A169F7', py: 1.5, mb: 1, '&:hover': { backgroundColor: '#A169F7' } }}
                      onClick={() => alert('Переход к оплате')}
                    >
                      Оплатить
                    </Button>
                    <Button
                        variant="outlined"
                        color="secondary"
                        fullWidth
                        sx={{
                            borderWidth: '3px',
                        }}
                        onClick={() => {
                            leaveQueue();
                            navigate(`/product/${product.id}`);
                        }}
                        >
                        Отменить заказ
                        </Button>
                  </Paper>
                </Grid>
              </Grid>

              {/*соглашение, доставка, возврат */}
              <Box sx={{ mt: 6, pt: 3, borderTop: '1px solid #e0e0e0', textAlign: 'center' }}>
                <Typography variant="body2" color="text.secondary">
                  Нажимая «Оплатить», вы соглашаетесь с условиями оферты и политикой обработки данных.
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  Доставка осуществляется в течение 2–3 рабочих дней. Возврат товара возможен в течение 14 дней.
                </Typography>
              </Box>
            </Container>
          </Box>
        );
      default:
        return (
          <Box>
            <Typography variant="h5">Неизвестный статус: {state.status}</Typography>
            <Button variant="outlined" onClick={() => navigate('/products')}>На главную</Button>
          </Box>
        );
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 0 }}>
      {renderContent()}
    </Container>
  );
};

export default QueuePage;