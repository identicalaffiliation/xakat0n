import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box, Typography, Button, Paper, Container, Grid, TextField,
  AppBar, Toolbar, IconButton, Divider
} from '@mui/material';
import { Person, FavoriteBorder } from '@mui/icons-material';
import { getItem, getSimilar, type Item } from '../../api/items';
import { getProductImage } from '../../utils/imageUtils';
import { getQueueStatus, cancelQueue, paymentCallback, type Ticket } from '../../api/queue';

const CheckoutPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Item | null>(null);
  const [loading, setLoading] = useState(true);
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [ticketChecked, setTicketChecked] = useState(false);
  const [timer, setTimer] = useState<number | null>(null);
  const [similar, setSimilar] = useState<Item[]>([]);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getItem(id)
      .then(setProduct)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  useEffect(() => {
    if (!id || !product?.is_limited) return;
    let interval: ReturnType<typeof setInterval>;
    const fetchStatus = async () => {
      try {
        const data = await getQueueStatus(id);
        setTicket(data);
        if (data.expiresInSeconds !== null) setTimer(data.expiresInSeconds);
        if (['PURCHASED', 'EXPIRED', 'SOLD_OUT', 'CANCELLED'].includes(data.status)) {
          clearInterval(interval);
        }
      } catch (error: any) {
        if (error?.response?.status === 404) {
          setTicket(null);
        } else {
          console.error('Ошибка обновления статуса заказа:', error);
        }
      } finally {
        setTicketChecked(true);
      }
    };
    fetchStatus();
    interval = setInterval(fetchStatus, 2000);
    return () => clearInterval(interval);
  }, [id, product?.is_limited]);

  useEffect(() => {
    if (ticket?.expiresInSeconds === null || ticket?.expiresInSeconds === undefined) {
      setTimer(null);
      return;
    }
    setTimer(ticket.expiresInSeconds);
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
  }, [ticket?.expiresInSeconds]);

  useEffect(() => {
    if (!id || !ticket) return;
    if (ticket.status === 'EXPIRED') {
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
  if (!product) return <Typography>Товар не найден</Typography>;

  const productImage = product.image_url || getProductImage(product.item_id);

  const renderProductCard = () => (
    <Grid size={{ xs: 12, md: 7 }}>
      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9' }}>
        <Box sx={{ display: 'flex', gap: 3 }}>
          <Box component="img" src={productImage} alt={product.title} onError={(e) => {
            (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
          }} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
          <Box>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
            <Typography variant="body2" color="text.secondary">Категория: {product.category || 'Без категории'}</Typography>
          </Box>
        </Box>
      </Paper>
    </Grid>
  );

  const renderSimilar = () => (
    similar.length > 0 && (
      <Box sx={{ mt: 4 }}>
        <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>Похожие товары</Typography>
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
          {similar.map(p => (
            <Paper key={p.item_id} sx={{ p: 1, width: 140, cursor: 'pointer' }} onClick={() => navigate(`/product/${p.item_id}`)}>
              <img src={p.image_url || getProductImage(p.item_id)} alt={p.title} onError={(e) => {
                (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
              }} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
              <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>{p.title}</Typography>
              <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block' }}>{p.price.toLocaleString()} ₽</Typography>
            </Paper>
          ))}
        </Box>
      </Box>
    )
  );

  const handleCancel = async () => {
    try {
      await cancelQueue(product.item_id);
    } catch (error) {
      console.error('Ошибка отмены заказа:', error);
    }
    navigate(`/product/${product.item_id}`);
  };

  // Обычный (не лимитированный) товар — тикета никогда не было и не будет, чекаут без таймера
  if (!product.is_limited) {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Оформление заказа</Typography>
          <Grid container spacing={4}>
            <Grid size={{ xs: 12, md: 7 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                <Box sx={{ display: 'flex', gap: 3 }}>
                  <Box component="img" src={productImage} alt={product.title} onError={(e) => {
                    (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
                  }} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                  <Box>
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                    <Typography variant="body2" color="text.secondary">Категория: {product.category || 'Без категории'}</Typography>
                    <Typography variant="body2" color="text.secondary">В наличии: {product.stock ?? '?'} шт.</Typography>
                  </Box>
                </Box>
              </Paper>
              <Paper elevation={0} sx={{ p: 2, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                <Typography variant="body2"><strong>Описание:</strong> {product.description || 'Товар без описания'}</Typography>
              </Paper>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>Данные получателя</Typography>
                <TextField fullWidth label="ФИО получателя" variant="outlined" size="small" value="Иванов Иван Иванович" sx={{ mb: 2 }} slotProps={{ input: { readOnly: true } }} />
                <TextField fullWidth label="Контактный телефон" variant="outlined" size="small" value="+7 (999) 123-45-67" slotProps={{ input: { readOnly: true } }} />
              </Paper>
              <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.item_id}`)}>Назад</Button>
            </Grid>
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>Детали заказа</Typography>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Товар</Typography><Typography sx={{ fontWeight: 500 }}>{product.title}</Typography></Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Количество</Typography><Typography sx={{ fontWeight: 500 }}>1 шт.</Typography></Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}><Typography>Стоимость</Typography><Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{product.price.toLocaleString()} ₽</Typography></Box>
                <Divider sx={{ my: 2 }} />
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                  <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Промокод</Typography><Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography></Paper>
                  <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Скидка по карте</Typography><Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography></Paper>
                </Box>
                <Button
                  variant="contained"
                  fullWidth
                  sx={{ bgcolor: '#A169F7', py: 1.5, mb: 1 }}
                  onClick={() => {
                    alert('Заказ оформлен! Оплата не требуется.');
                    navigate('/products');
                  }}
                >
                  Оплатить
                </Button>
                <Button variant="outlined" color="secondary" fullWidth sx={{ borderWidth: '3px' }} onClick={() => navigate('/products')}>Отменить заказ</Button>
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
  }

  // Лимитированный товар, тикет ещё не подгрузился первым поллингом
  if (!ticketChecked) {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography>Загрузка...</Typography>
        </Container>
      </Box>
    );
  }

  // Лимитированный товар, но активного оформления нет (прямой заход по ссылке без права)
  if (!ticket || ticket.status === 'QUEUED' || ticket.status === 'OFFERED' || ticket.status === 'SOLD_OUT') {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Оформление недоступно</Typography>
          <Grid container spacing={4}>
            {renderProductCard()}
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                <Typography variant="body1" sx={{ mb: 2 }}>
                  У вас нет активного оформления по этому товару.
                </Typography>
                <Button variant="contained" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${product.item_id}/queue`)}>
                  К очереди
                </Button>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>
    );
  }

  if (ticket.status === 'EXPIRED') {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Время вышло</Typography>
          <Grid container spacing={4}>
            {renderProductCard()}
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#FF6163' }}>Время вышло</Typography>
                <Typography variant="body1" sx={{ mb: 2 }}>
                  Право на покупку истекло, товар вернулся в продажу и достался следующему в очереди.
                </Typography>
                <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${product.item_id}`)}>
                  Встать в очередь заново
                </Button>
              </Paper>
            </Grid>
          </Grid>
          {renderSimilar()}
        </Container>
      </Box>
    );
  }

  if (ticket.status === 'CANCELLED') {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Вы отказались от покупки</Typography>
          <Grid container spacing={4}>
            {renderProductCard()}
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                <Typography variant="body1" sx={{ mb: 2 }}>
                  Место передано следующему участнику. Можно встать заново — в конец очереди.
                </Typography>
                <Button variant="contained" color="primary" fullWidth sx={{ py: 1.5 }} onClick={() => navigate(`/product/${product.item_id}`)}>
                  Встать в очередь заново
                </Button>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>
    );
  }

  if (ticket.status === 'PURCHASED') {
    return (
      <Box>
        {renderHeader()}
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Покупка оформлена</Typography>
          <Grid container spacing={4}>
            {renderProductCard()}
            <Grid size={{ xs: 12, md: 5 }}>
              <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                <Typography variant="body1" sx={{ mb: 2 }}>
                  Заказ передан в доставку. Спасибо за покупку!
                </Typography>
                <Button variant="contained" fullWidth sx={{ bgcolor: '#04e162', py: 1.5, mb: 1 }} onClick={() => alert('Перейти к заказу')}>
                  Перейти к заказу
                </Button>
                <Button variant="outlined" fullWidth sx={{ py: 1.5 }} onClick={() => navigate('/products')}>
                  Перейти к каталогу
                </Button>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>
    );
  }

  // ticket.status === 'CHECKOUT' — активное оформление лимитированного товара
  return (
    <Box>
      {renderHeader()}
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
          Оформление заказа
        </Typography>
        <Grid container spacing={4}>
          <Grid size={{ xs: 12, md: 7 }}>
            <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
              <Box sx={{ display: 'flex', gap: 3 }}>
                <Box component="img" src={productImage} alt={product.title} onError={(e) => {
                  (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
                }} sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }} />
                <Box>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>{product.title}</Typography>
                  <Typography variant="body2" color="text.secondary">Категория: {product.category || 'Без категории'}</Typography>
                  <Typography variant="body2" color="text.secondary">В наличии: {product.stock ?? '?'} шт.</Typography>
                </Box>
              </Box>
            </Paper>
            <Paper elevation={0} sx={{ p: 2, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
              <Typography variant="body2"><strong>Описание:</strong> {product.description || 'Товар без описания'}</Typography>
            </Paper>
            <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>Данные получателя</Typography>
              <TextField fullWidth label="ФИО получателя" variant="outlined" size="small" value="Иванов Иван Иванович" sx={{ mb: 2 }} slotProps={{ input: { readOnly: true } }} />
              <TextField fullWidth label="Контактный телефон" variant="outlined" size="small" value="+7 (999) 123-45-67" slotProps={{ input: { readOnly: true } }} />
            </Paper>
            <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.item_id}/queue`)}>Назад</Button>
          </Grid>
          <Grid size={{ xs: 12, md: 5 }}>
            <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>Детали заказа</Typography>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Товар</Typography><Typography sx={{ fontWeight: 500 }}>{product.title}</Typography></Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Количество</Typography><Typography sx={{ fontWeight: 500 }}>1 шт.</Typography></Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}><Typography>Стоимость</Typography><Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{product.price.toLocaleString()} ₽</Typography></Box>
              <Divider sx={{ my: 2 }} />
              <Box sx={{ textAlign: 'center', mb: 2 }}>
                <Typography variant="h5" sx={{ fontWeight: 700 }}>Осталось: {formatTime(timer ?? 0)}</Typography>
                <Typography variant="body2" color="text.secondary">Товар держится за вами, пока идёт время</Typography>
              </Box>
              <Divider sx={{ my: 2 }} />
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Промокод</Typography><Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography></Paper>
                <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Скидка по карте</Typography><Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography></Paper>
              </Box>
              <Button
                variant="contained"
                fullWidth
                sx={{ bgcolor: '#A169F7', py: 1.5, mb: 1 }}
                onClick={async () => {
                  try {
                    const updated = await paymentCallback(product.item_id, ticket.ticketId, 'paid');
                    setTicket(updated);
                  } catch (error) {
                    console.error('Ошибка оплаты:', error);
                    alert('Ошибка при оплате. Попробуйте ещё раз.');
                  }
                }}
              >
                Оплатить
              </Button>
              <Button variant="outlined" color="secondary" fullWidth sx={{ borderWidth: '3px' }} onClick={handleCancel}>
                Отменить заказ
              </Button>
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
};

export default CheckoutPage;
