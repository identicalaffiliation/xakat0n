import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box, Typography, Button, Paper, Container, CircularProgress, Divider, Grid,
  TextField, AppBar, Toolbar, IconButton
} from '@mui/material'; 
import { Person, FavoriteBorder } from '@mui/icons-material';
import { getQueueStatus, cancelQueue, startCheckout, type Ticket } from '../../api/queue';
import { getSimilar, type Item } from '../../api/items';

const QueuePage: React.FC = () => {
  const { id } = useParams<{ id: string }>(); 
  const navigate = useNavigate();
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [similar, setSimilar] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [timer, setTimer] = useState<number | null>(null);

  useEffect(() => {
    if (!id) return;
    let interval: ReturnType<typeof setInterval>;
    const fetchStatus = async () => {
      try {
        const data = await getQueueStatus(id);
        setTicket(data);
        if (data.expires_in_seconds !== null) setTimer(data.expires_in_seconds);
        if (['PURCHASED', 'EXPIRED', 'SOLD_OUT', 'CANCELLED'].includes(data.status)) {
          clearInterval(interval);
        }
      } catch (error) {
        setTicket(null);
      } finally {
        setLoading(false);
      }
    };
    fetchStatus();
    interval = setInterval(fetchStatus, 2000);
    return () => clearInterval(interval);
  }, [id]);

  useEffect(() => {
    if (ticket?.expires_in_seconds === null || ticket?.expires_in_seconds === undefined) {
      setTimer(null);
      return;
    }
    setTimer(ticket.expires_in_seconds);
    const interval = setInterval(() => {
      setTimer(prev => {
        if (prev === null || prev <= 1) {
          clearInterval(interval);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [ticket?.expires_in_seconds]);
  useEffect(() => {
    if (!id) return;
    if (ticket && ['QUEUED', 'EXPIRED', 'CANCELLED'].includes(ticket.status)) {
      getSimilar(id, 3).then(setSimilar).catch(console.error);
    }
  }, [id, ticket?.status]);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const renderHeader = () => (
    <AppBar position="sticky" color="default" elevation={0} sx={{ bgcolor: '#fff', top: 0, zIndex: 1100 }}>
      <Toolbar sx={{ justifyContent: 'space-between', py: 1 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', ml: -20 }}>
          <img src="/avito-logo.svg" alt="Avito" style={{ width: 48, height: 48, marginRight: 12 }} />
          <Typography variant="h5" sx={{ fontWeight: 700, color: '#000', fontSize: '1.8rem' }}>Avito</Typography>
        </Box>
        <Box>
          <IconButton color="inherit"><FavoriteBorder /></IconButton>
          <IconButton color="inherit"><Person /></IconButton>
        </Box>
      </Toolbar>
    </AppBar>
  );

  if (loading) return <Typography>Загрузка...</Typography>;
  if (!ticket) return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h5">Нет активной заявки</Typography>
      <Button variant="outlined" onClick={() => navigate('/products')}>В каталог</Button>
    </Container>
  );

  const productImage = `https://picsum.photos/seed/${ticket.item_id}/200/200`;

  const renderContent = () => {
    switch (ticket.status) {
      case 'CHECKOUT':
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Оформление заказа</Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                      <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                        <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 2, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="body2"><strong>Описание:</strong> Товар зарезервирован для вас.</Typography>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>Данные получателя</Typography>
                    <TextField fullWidth label="ФИО получателя" variant="outlined" size="small" value="Иванов Иван Иванович" sx={{ mb: 2 }} slotProps={{ input: { readOnly: true } }} />
                    <TextField fullWidth label="Контактный телефон" variant="outlined" size="small" value="+7 (999) 123-45-67" slotProps={{ input: { readOnly: true } }} />
                  </Paper>
                  <Button variant="contained" color="primary" onClick={() => navigate(`/product/${ticket.item_id}`)}>Назад</Button>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>Детали заказа</Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Товар</Typography><Typography sx={{ fontWeight: 500 }}>ID {ticket.item_id}</Typography></Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Количество</Typography><Typography sx={{ fontWeight: 500 }}>1 шт.</Typography></Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}><Typography>Стоимость</Typography><Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>—</Typography></Box>
                    {timer !== null && (
                      <>
                        <Divider sx={{ my: 2 }} />
                        <Box sx={{ textAlign: 'center', mb: 2 }}>
                          <Typography variant="body2" color="text.secondary">Осталось времени для оплаты</Typography>
                          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 2, mt: 1 }}>
                            <CircularProgress variant="determinate" value={((timer || 0) / 120) * 100} size={64} />
                            <Typography variant="h2" sx={{ fontWeight: 700, fontSize: '2.5rem' }}>{formatTime(timer || 0)}</Typography>
                          </Box>
                        </Box>
                      </>
                    )}
                    <Divider sx={{ my: 2 }} />
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Промокод</Typography><Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography></Paper>
                      <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Скидка по карте</Typography><Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography></Paper>
                    </Box>
                    <Button variant="contained" fullWidth sx={{ bgcolor: '#A169F7', py: 1.5, mb: 1, '&:hover': { bgcolor: '#A169F7' } }} onClick={() => alert('Переход к оплате')}>Оплатить</Button>
                    <Button variant="outlined" color="secondary" fullWidth sx={{ borderWidth: '3px' }} onClick={async () => { await cancelQueue(ticket.item_id); navigate(`/product/${ticket.item_id}`); }}>Отменить заказ</Button>
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
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Очередь ожидания на лимитированный товар</Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                      <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                        <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 6, mb: 2 }}>
                      <Box>
                        <Typography variant="body2" color="text.secondary">Ваше место</Typography>
                        <Typography variant="h5" sx={{ fontWeight: 700 }}>{ticket.position ?? '?'}</Typography>
                      </Box>
                      <Box>
                        <Typography variant="body2" color="text.secondary">Макс. время ожидания</Typography>
                        <Typography variant="h5" sx={{ fontWeight: 700 }}>
                          {ticket.next_slot_free_in_seconds ? `${Math.floor(ticket.next_slot_free_in_seconds / 60)} мин` : '~8 мин'}
                        </Typography>
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
                    <Button variant="contained" color="primary" onClick={() => navigate('/products')}>Перейти в каталог</Button>
                    <Button variant="outlined" color="secondary" onClick={async () => { await cancelQueue(ticket.item_id); navigate('/products'); }}>Выйти из очереди</Button>
                  </Box>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус ожидания</Typography>
                    <Box sx={{ textAlign: 'center', mb: 2 }}>
                      <CircularProgress variant="determinate" value={((timer || 0) / 480) * 100} size={80} sx={{ color: 'primary.main' }} />
                      <Typography variant="h4" sx={{ fontWeight: 700, mt: 1 }}>{formatTime(timer || 0)}</Typography>
                      <Typography variant="body2" color="text.secondary">Осталось до получения права</Typography>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="body1" gutterBottom><strong>Краткая информация</strong></Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>Когда подойдёт ваша очередь, вы получите уведомление и сможете перейти к оформлению.</Typography>
                    <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${ticket.item_id}`)}>Назад к товару</Button>
                  </Paper>
                </Grid>
              </Grid>
              <Box sx={{ mt: 4 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
                <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                  {similar.map(p => (
                    <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
                      <img src={`https://picsum.photos/seed/${p.item_id}/200/200`} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Очередь ожидания на лимитированный товар</Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                      <Box><Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography><Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography></Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#E8F5E9', border: '2px solid #00C853', mb: 3 }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, color: '#00C853', mb: 1 }}>Товар освободился!</Typography>
                    <Typography variant="body1" sx={{ mb: 2 }}>Вы можете перейти к оформлению заказа. У вас есть {formatTime(timer || 60)} на оплату.</Typography>
                    <Button variant="contained" color="primary" size="large" onClick={async () => { await startCheckout(ticket.item_id); navigate(`/product/${ticket.item_id}/queue`); }}>Перейти к оформлению</Button>
                    <Button variant="outlined" color="secondary" sx={{ ml: 2 }} onClick={async () => { await cancelQueue(ticket.item_id); navigate('/products'); }}>Отказаться</Button>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 6, mb: 2 }}>
                      <Box><Typography variant="body2" color="text.secondary">Макс. время ожидания</Typography><Typography variant="h5" sx={{ fontWeight: 700 }}>~8 мин</Typography></Box>
                    </Box>
                    <Typography variant="body1" color="text.secondary" sx={{ mb: 1 }}>Краткое описание</Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>Товар освободился! У вас есть {formatTime(timer || 60)} на оформление.</Typography>
                  </Paper>
                  <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
                    <Button variant="contained" color="primary" onClick={() => navigate('/products')}>Перейти в каталог</Button>
                    <Button variant="outlined" color="secondary" onClick={async () => { await cancelQueue(ticket.item_id); navigate('/products'); }}>Выйти из очереди</Button>
                  </Box>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                    <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус ожидания</Typography>
                    <Box sx={{ textAlign: 'center', mb: 2 }}>
                      <CircularProgress variant="determinate" value={((timer || 0) / 60) * 100} size={80} sx={{ color: 'primary.main' }} />
                      <Typography variant="h4" sx={{ fontWeight: 700, mt: 1 }}>{formatTime(timer || 0)}</Typography>
                      <Typography variant="body2" color="text.secondary">Осталось для оформления</Typography>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="body1" gutterBottom><strong>Краткая информация</strong></Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>Вы получили право на покупку! Оплатите в течение {formatTime(timer || 60)}.</Typography>
                    <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${ticket.item_id}`)}>Назад к товару</Button>
                  </Paper>
                </Grid>
              </Grid>
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
                    <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                      <Box sx={{ display: 'flex', gap: 3 }}>
                        <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                        <Box>
                          <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                          <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                        </Box>
                      </Box>
                    </Paper>
                    <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                      <Typography variant="body1" sx={{ mb: 1 }}>
                        <strong>Время на оплату истекло</strong>
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Вы не оплатили товар вовремя. Товар вернулся в продажу и достался следующему в очереди.
                      </Typography>
                      <Box sx={{ display: 'flex', gap: 4, mt: 2 }}>
                        <Box>
                          <Typography variant="body2" color="text.secondary">Ваше место</Typography>
                          <Typography variant="h6" sx={{ fontWeight: 700 }}>—</Typography>
                        </Box>
                        <Box>
                          <Typography variant="body2" color="text.secondary">Макс. время ожидания</Typography>
                          <Typography variant="h6" sx={{ fontWeight: 700 }}>~8 мин</Typography>
                        </Box>
                      </Box>
                    </Paper>
                  </Grid>
                  <Grid size={{ xs: 12, md: 5 }}>
                    <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                      <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                      <Box sx={{ textAlign: 'center', py: 2 }}>
                        <Typography variant="h6" sx={{ fontWeight: 700, color: '#FF6163', mb: 1 }}>
                          Вы не оплатили товар вовремя
                        </Typography>
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
                          navigate(`/product/${ticket.item_id}`);
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
                    {similar.map(p => (
                      <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
                        <img src={`https://picsum.photos/seed/${p.item_id}/200/200`} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
                      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                        <Box sx={{ display: 'flex', gap: 3 }}>
                          <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                          <Box>
                            <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                            <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                          </Box>
                        </Box>
                      </Paper>
                      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                        <Typography variant="body1" sx={{ mb: 1 }}>
                          <strong>Вы вышли из очереди</strong>
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Ваше место передано следующему участнику. Вы можете встать в очередь заново (в конец).
                        </Typography>
                      </Paper>
                    </Grid>
                    <Grid size={{ xs: 12, md: 5 }}>
                      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                        <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                        <Box sx={{ textAlign: 'center', py: 3 }}>
                          <Typography variant="h2" sx={{ fontSize: '4rem' }}></Typography>
                          <Typography variant="h5" sx={{ fontWeight: 700, color: '#FF6163' }}>Выход из очереди</Typography>
                        </Box>
                        <Button
                          variant="contained"
                          color="primary"
                          fullWidth
                          sx={{ py: 1.5 }}
                          onClick={() => navigate(`/product/${ticket.item_id}`)}
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
                      {similar.map(p => (
                        <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
                          <img src={`https://picsum.photos/seed/${p.item_id}/200/200`} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
                        <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                          <Box sx={{ display: 'flex', gap: 3 }}>
                            <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                            <Box>
                              <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                              <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                            </Box>
                          </Box>
                        </Paper>
                        <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                          <Typography variant="body1" sx={{ mb: 1 }}>
                            <strong>Заказ передан в доставку</strong>
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
                        <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                          <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                          <Box sx={{ textAlign: 'center', py: 3 }}>
                            <Typography variant="h2" sx={{ fontSize: '4rem' }}></Typography>
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
                        {similar.map(p => (
                          <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
                            <img src={`https://picsum.photos/seed/${p.item_id}/200/200`} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
                            <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
                            <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
                          </Paper>
                        ))}
                      </Box>
                    </Box>
                  </Container>
                </Box>
              );
              case 'SOLD_OUT':
                return (
                  <Box>
                    {renderHeader()}
                    <Container maxWidth="lg" sx={{ py: 4 }}>
                      <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                        Товар раскуплен
                      </Typography>
                      <Grid container spacing={4}>
                        <Grid size={{ xs: 12, md: 7 }}>
                          <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                            <Box sx={{ display: 'flex', gap: 3 }}>
                              <Box component="img" src={productImage} alt="Товар" sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                              <Box>
                                <Typography variant="h6" sx={{ fontWeight: 600 }}>Товар</Typography>
                                <Typography variant="body2" color="text.secondary">ID: {ticket.item_id}</Typography>
                              </Box>
                            </Box>
                          </Paper>
                          <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                            <Typography variant="body1" sx={{ mb: 1 }}>
                              <strong>Все единицы выкуплены</strong>
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              Товар полностью распродан. Вы можете посмотреть похожие товары у других продавцов.
                            </Typography>
                          </Paper>
                        </Grid>
                        <Grid size={{ xs: 12, md: 5 }}>
                          <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                            <Typography variant="h6" gutterBottom sx={{ fontWeight: 700 }}>Статус</Typography>
                            <Box sx={{ textAlign: 'center', py: 3 }}>
                              <Typography variant="h2" sx={{ fontSize: '4rem' }}></Typography>
                              <Typography variant="h5" sx={{ fontWeight: 700, color: '#FF6163' }}>Раскуплено</Typography>
                            </Box>
                            <Button variant="contained" fullWidth sx={{ py: 1.5 }} onClick={() => navigate('/products')}>
                              Перейти в каталог
                            </Button>
                          </Paper>
                        </Grid>
                      </Grid>
                      <Box sx={{ mt: 4 }}>
                        <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
                        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                          {similar.map(p => (
                            <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
                              <img src={`https://picsum.photos/seed/${p.item_id}/200/200`} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
              <Typography variant="h4">Статус: {ticket.status}</Typography>
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