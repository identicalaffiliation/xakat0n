import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box, Typography, Button, Paper, Container, Grid, TextField,
  AppBar, Toolbar, IconButton, Divider
} from '@mui/material';
import { Person, FavoriteBorder } from '@mui/icons-material';
import { getItem, type Item } from '../../api/items';
import { getProductImage } from '../../utils/imageUtils';
import { getQueueStatus, paymentCallback, type Ticket } from '../../api/queue';
const CheckoutPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Item | null>(null);
  const [loading, setLoading] = useState(true);
  const [ticket, setTicket] = useState<Ticket | null>(null);
  useEffect(() => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      getItem(id),
      getQueueStatus(id).catch(() => null)
    ])
      .then(([productData, ticketData]) => {
        setProduct(productData);
        setTicket(ticketData);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);
  // useEffect(() => {
  //   if (!id) return;
  //   getItem(id)
  //     .then(data => setProduct(data))
  //     .catch(console.error)
  //     .finally(() => setLoading(false));
  // }, [id]);

  if (loading) return <Typography>Загрузка...</Typography>;
  if (!product) return <Typography>Товар не найден</Typography>;

  // const productImage = `https://picsum.photos/seed/${product.item_id}/200/200`;
  const productImage = getProductImage(product.item_id);
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
            <Button variant="contained" color="primary" onClick={() => navigate(`/product/${product.item_id}`)}>Назад</Button>
          </Grid>
          <Grid size={{ xs: 12, md: 5 }}>
            <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5', mt: -2 }}>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 700, fontSize: '1.25rem' }}>Детали заказа</Typography>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Товар</Typography><Typography sx={{ fontWeight: 500 }}>{product.title}</Typography></Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}><Typography>Количество</Typography><Typography sx={{ fontWeight: 500 }}>1 шт.</Typography></Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}><Typography>Стоимость</Typography><Typography sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{product.price.toLocaleString()} ₽</Typography></Box>
              <Divider sx={{ my: 2 }} />
              {/* Нет таймера */}
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Промокод</Typography><Typography sx={{ fontWeight: 500 }}>AVITO10 – скидка 10%</Typography></Paper>
                <Paper variant="outlined" sx={{ p: 1.5, borderRadius: 2, bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">Скидка по карте</Typography><Typography sx={{ fontWeight: 500 }}>5% (накопительная)</Typography></Paper>
              </Box>
                {ticket ? (
                  <Button
                    variant="contained"
                    fullWidth
                    sx={{ bgcolor: '#A169F7', py: 1.5, mb: 1 }}
                    onClick={async () => {
                      try {
                        await paymentCallback(product.item_id, ticket.ticketId, 'paid');
                        alert('Оплата прошла успешно!');
                        navigate(`/product/${product.item_id}`);
                      } catch (error) {
                        console.error('Ошибка оплаты:', error);
                        alert('Ошибка при оплате');
                      }
                    }}
                  >
                    Оплатить
                  </Button>
                ) : (
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
                )}
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
};

export default CheckoutPage;