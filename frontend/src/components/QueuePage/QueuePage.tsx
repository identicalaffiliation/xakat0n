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
  const { state, leaveQueue, expireOffer, forceStatus, joinQueue } = useQueue();
  const product = products.find((p) => p.id === Number(id));
  const [timer, setTimer] = useState<number | null>(null);
  useEffect(() => {
    if (!product) return;
    if (!state.status || state.timeLeft === null) return;
    setTimer(state.timeLeft);
    const interval = setInterval(() => {
      setTimer((prev) => {
        if (prev === null || prev <= 1) {
          clearInterval(interval);
          if (state.status === 'QUEUED') {
            forceStatus(product.id, 'OFFERED', 60);
          }
          if (state.status === 'CHECKOUT' || state.status === 'OFFERED') {
            expireOffer();
          }
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [state.status, state.timeLeft, expireOffer, forceStatus, product]);

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

  const renderHeader = () => (
    <AppBar position="sticky" color="default" elevation={0} sx={{ backgroundColor: '#fff', top: 0, zIndex: 1100 }}>
      <Toolbar sx={{ justifyContent: 'space-between', py: 1 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', ml: -20 }}>
          <img
            src="/avito-logo.svg"
            alt="Avito"
            style={{ width: 48, height: 48, marginRight: 12 }}
          />
          <Typography variant="h5" sx={{ fontWeight: 700, color: '#000000', fontSize: '1.8rem' }}>
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
  );

  const renderContent = () => {
    if (!product) return null;

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
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                Оформление заказа
              </Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
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
                        <Typography variant="body2" color="text.secondary">Категория: {product.category}</Typography>
                        <Typography variant="body2" color="text.secondary">В наличии: {product.stock} шт.</Typography>
                        {product.is_limited && <Chip label="Лимитированный" size="small" sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }} />}
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 2, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="body2"><strong>Описание:</strong> {product.title} – отличный выбор.</Typography>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>Данные получателя</Typography>
                    <TextField fullWidth label="ФИО получателя" variant="outlined" size="small" value="Иванов Иван Иванович" sx={{ mb: 2 }} slotProps={{ input: { readOnly: true } }} />
                    <TextField fullWidth label="Контактный телефон" variant="outlined" size="small" value="+7 (999) 123-45-67" slotProps={{ input: { readOnly: true } }} />
                  </Paper>
                  <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.id}`)}>Назад</Button>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>Детали заказа</Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Товар</Typography><Typography sx={{ fontWeight: 500 }}>{product.title}</Typography></Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Количество</Typography><Typography sx={{ fontWeight: 500 }}>1 шт.</Typography></Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}><Typography>Стоимость</Typography><Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{product.price.toLocaleString()} ₽</Typography></Box>
                    <Divider sx={{ my: 2 }} />
                    <Box sx={{ textAlign: 'center', mb: 2 }}>
                      <Typography variant="body2" color="text.secondary">Осталось времени для оплаты</Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 2, mt: 1 }}>
                        <CircularProgress variant="determinate" value={((timer || 0) / 120) * 100} size={64} />
                        <Typography variant="h2" sx={{ fontWeight: 700, fontSize: '2.5rem' }}>{formatTime(timer || 0)}</Typography>
                      </Box>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, backgroundColor: '#fff' }}><Typography variant="caption" color="text.secondary">Промокод</Typography><Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography></Paper>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, backgroundColor: '#fff' }}><Typography variant="caption" color="text.secondary">Скидка по карте</Typography><Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography></Paper>
                    </Box>
                    <Button variant="contained" fullWidth sx={{ backgroundColor: '#A169F7', py: 1.5, mb: 1, '&:hover': { backgroundColor: '#A169F7' } }} onClick={() => alert('Переход к оплате')}>Оплатить</Button>
                    <Button variant="outlined" color="secondary" fullWidth sx={{ borderWidth: '3px' }} onClick={() => { leaveQueue(); navigate(`/product/${product.id}`); }}>Отменить заказ</Button>
                  </Paper>
                </Grid>
              </Grid>
              <Box sx={{ mt: 6, pt: 3, borderTop: '1px solid #e0e0e0', textAlign: 'center' }}>
                <Typography variant="body2" color="text.secondary">Нажимая «Оплатить», вы соглашаетесь с условиями оферты и политикой обработки данных.</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>Доставка осуществляется в течение 2–3 рабочих дней. Возврат товара возможен в течение 14 дней.</Typography>
              </Box>
            </Container>
          </Box>
        );

      case 'QUEUED':
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                Очередь ожидания на лимитированный товар
              </Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
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
                        <Typography variant="body2" color="text.secondary">Категория: {product.category}</Typography>
                        <Typography variant="body2" color="text.secondary">В наличии: {product.stock} шт.</Typography>
                        {product.is_limited && <Chip label="Лимитированный" size="small" sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }} />}
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 6, mb: 2 }}>
                      <Box>
                        <Typography variant="body2" color="text.secondary">В очереди</Typography>
                        <Typography variant="h5" sx={{ fontWeight: 700 }}>{state.totalInQueue}</Typography>
                      </Box>

                      <Box>
                        <Typography variant="body2" color="text.secondary">Макс. время ожидания</Typography>
                        <Typography variant="h5" sx={{ fontWeight: 700 }}>~8 мин</Typography>
                      </Box>
                    </Box>
                    <Typography variant="body1" color="text.secondary" sx={{ mb: 1 }}>
                      Краткое описание в зависимости от статуса, чтобы дать пользователю понять на каком этапе он находится
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                      Вы в очереди. Как только товар освободится, у вас будет 1 минута на оформление.
                    </Typography>
                  </Paper>
                  <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
                    <Button variant="contained" color="primary" onClick={() => navigate('/products')}>
                      Перейти в каталог
                    </Button>
                    <Button
                      variant="outlined"
                      color="secondary"
                      onClick={() => {
                        leaveQueue();
                        navigate('/products');
                      }}
                    >
                      Выйти из очереди
                    </Button>
                  </Box>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус ожидания</Typography>
                    <Box sx={{ textAlign: 'center', mb: 2 }}>
                      <CircularProgress variant="determinate" value={((timer || 0) / 480) * 100} size={80} sx={{ color: 'primary.main' }} />
                      <Typography variant="h4" sx={{ fontWeight: 700, mt: 1 }}>{formatTime(timer || 0)}</Typography>
                      <Typography variant="body2" color="text.secondary">Осталось до получения права</Typography>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="body1" gutterBottom><strong>Краткая информация</strong></Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>Когда подойдёт ваша очередь, вы получите уведомление и сможете перейти к оформлению.</Typography>
                    <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${product.id}`)}>Назад к товару</Button>
                  </Paper>
                </Grid>
              </Grid>
              <Box sx={{ mt: 4 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
                <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                  {products.filter(p => p.id !== product.id).slice(0, 3).map(p => (
                    <Paper key={p.id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.id}`)}>
                      <img src={p.image} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
                      <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
                      <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
                    </Paper>
                  ))}
                </Box>
              </Box>
            </Container>
          </Box>
        );

      case 'OFFERED':
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom>Право выдано!</Typography>
              <Typography variant="body1" color="text.secondary">У вас есть {formatTime(timer || 60)} на оформление.</Typography>
              <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.id}/checkout`)}>Перейти к оформлению</Button>
              <Button variant="outlined" color="secondary" onClick={() => { leaveQueue(); navigate(`/product/${product.id}`); }}>Отменить</Button>
            </Container>
          </Box>
        );

      case 'EXPIRED':
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                Время вышло
              </Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box component="img" src={product.image} alt={product.title} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                      <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                        <Typography variant="body2" color="text.secondary">Категория: {product.category}</Typography>
                        <Typography variant="body2" color="text.secondary">В наличии: {product.stock} шт.</Typography>
                        {product.is_limited && <Chip label="Лимитированный" size="small" sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }} />}
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="body1" sx={{ mb: 1 }}>
                      <strong>Макс. время ожидания в очереди</strong>
                    </Typography>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>~8 мин</Typography>
                    <Typography variant="body2" color="text.secondary">
                      Оплата не прошла, вы можете встать в очередь и, если другой пользователь откажется от товара - мы обязательно сообщим вам о возможности покупки!
                    </Typography>
                  </Paper>
                  {/* Левая кнопка удалена */}
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                    <Box sx={{ textAlign: 'center', py: 2 }}>
                      <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'left' }}>
                        К сожалению, вы не успели завершить оплату в отведённое время. Товар вернулся в продажу и достался следующему в очереди.
                      </Typography>
                    </Box>
                    <Button
                      variant="contained"
                      color="primary"
                      fullWidth
                      sx={{ py: 1.5 }}
                      onClick={() => {
                        joinQueue(product.id);
                        navigate(`/product/${product.id}/queue`);
                      }}
                    >
                      Встать в очередь
                    </Button>
                    <Button variant="outlined" fullWidth sx={{ mt: 1 }} onClick={() => navigate('/products')}>
                      На главную
                    </Button>
                  </Paper>
                </Grid>
              </Grid>
              <Box sx={{ mt: 4 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
                <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                  {products.filter(p => p.id !== product.id).slice(0, 3).map(p => (
                    <Paper key={p.id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.id}`)}>
                      <img src={p.image} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
                      <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
                      <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
                    </Paper>
                  ))}
                </Box>
              </Box>
            </Container>
          </Box>
        );
  case 'CANCELLED':
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
            Вы вышли из очереди
          </Typography>
          <Grid container spacing={4}>
            <Grid size={{ xs: 12, md: 7 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                <Box sx={{ display: 'flex', gap: 3 }}>
                  <Box component="img" src={product.image} alt={product.title} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                  <Box>
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                    <Typography variant="body2" color="text.secondary">Категория: {product.category}</Typography>
                    <Typography variant="body2" color="text.secondary">В наличии: {product.stock} шт.</Typography>
                    {product.is_limited && <Chip label="Лимитированный" size="small" sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }} />}
                  </Box>
                </Box>
              </Paper>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                <Typography variant="body1" sx={{ mb: 1 }}>
                  <strong>🚪 Вы вышли из очереди</strong>
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Ваше место передано следующему участнику. Вы можете встать в очередь заново (в конец).
                </Typography>
              </Paper>
              <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.id}`)}>
                Встать в очередь заново
              </Button>
            </Grid>
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                <Box sx={{ textAlign: 'center', py: 3 }}>
                  <Typography variant="h2" sx={{ fontSize: '4rem' }}>🚪</Typography>
                  <Typography variant="h5" sx={{ fontWeight: 700, color: '#FF6163' }}>Выход из очереди</Typography>
                </Box>
                <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${product.id}`)}>
                  Встать в очередь
                </Button>
                <Button variant="outlined" fullWidth sx={{ mt: 1 }} onClick={() => navigate('/products')}>
                  На главную
                </Button>
              </Paper>
            </Grid>
          </Grid>
          <Box sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
              {products.filter(p => p.id !== product.id).slice(0, 3).map(p => (
                <Paper key={p.id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.id}`)}>
                  <img src={p.image} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
                  <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
                  <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
                </Paper>
              ))}
            </Box>
          </Box>
        </Container>
      </Box>
    );

  case 'PURCHASED':
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
            Покупка оформлена
          </Typography>
          <Grid container spacing={4}>
            <Grid size={{ xs: 12, md: 7 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                <Box sx={{ display: 'flex', gap: 3 }}>
                  <Box component="img" src={product.image} alt={product.title} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                  <Box>
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                    <Typography variant="body2" color="text.secondary">Категория: {product.category}</Typography>
                    <Typography variant="body2" color="text.secondary">В наличии: {product.stock} шт.</Typography>
                    {product.is_limited && <Chip label="Лимитированный" size="small" sx={{ mt: 1, backgroundColor: '#FF6163', color: '#fff' }} />}
                  </Box>
                </Box>
              </Paper>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f9f9f9', mb: 3 }}>
                <Typography variant="body1" sx={{ mb: 1 }}>
                  <strong>✅ Заказ передан в доставку</strong>
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Спасибо за покупку! Вы можете отслеживать статус заказа в личном кабинете.
                </Typography>
              </Paper>
              <Button variant="contained" color="primary" onClick={() => alert('Перейти к заказу')}>
                Перейти к заказу
              </Button>
            </Grid>
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, backgroundColor: '#f5f5f5', mt: -2 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                <Box sx={{ textAlign: 'center', py: 3 }}>
                  <Typography variant="h2" sx={{ fontSize: '4rem' }}>✅</Typography>
                  <Typography variant="h5" sx={{ fontWeight: 700, color: '#00C853' }}>Покупка оформлена</Typography>
                </Box>
                <Button variant="contained" color="success" fullWidth sx={{ py: 1.5 }} onClick={() => alert('Перейти к заказу')}>
                  Перейти к заказу
                </Button>
                <Button variant="outlined" fullWidth sx={{ mt: 1 }} onClick={() => navigate('/products')}>
                  Купить ещё
                </Button>
              </Paper>
            </Grid>
          </Grid>
          <Box sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
              {products.filter(p => p.id !== product.id).slice(0, 3).map(p => (
                <Paper key={p.id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.id}`)}>
                  <img src={p.image} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
                  <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
                  <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
                </Paper>
              ))}
            </Box>
          </Box>
        </Container>
      </Box>
    );

      default:
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h5">Неизвестный статус: {state.status}</Typography>
              <Button variant="outlined" onClick={() => navigate('/products')}>На главную</Button>
            </Container>
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