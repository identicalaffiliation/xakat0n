import React, { useEffect, useState } from 'react';
import { useNavigate, useParams} from 'react-router-dom';
import {
  Box, Typography, Button, Paper, Container, CircularProgress, Grid,
  AppBar, Toolbar, IconButton
} from '@mui/material';
import { Person, FavoriteBorder } from '@mui/icons-material';
import { getQueueStatus, cancelQueue, startCheckout, type Ticket } from '../../api/queue';
import { getItem, getSimilar, type Item } from '../../api/items';
import { getProductImage } from '../../utils/imageUtils';

const QueuePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [product, setProduct] = useState<Item | null>(null);
  const [similar, setSimilar] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [timer, setTimer] = useState<number | null>(null);

  // const [searchParams] = useSearchParams();
  // const mockStatus = searchParams.get('mockStatus') as Ticket['status'] | null;
  useEffect(() => {
    if (!id) return;
    getItem(id).then(setProduct).catch(console.error);
  }, [id]);

  useEffect(() => {
    if (!id) return;
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
          console.error('Ошибка обновления статуса очереди:', error);
        }
      } finally {
        setLoading(false);
      }
    };
    fetchStatus();
    interval = setInterval(fetchStatus, 2000);
    return () => clearInterval(interval);
  }, [id]);
  // useEffect(() => {
  //   if (!id) return;
  //   if (mockStatus) {
  //     //Создаём мок
  //     const mockTicket: Ticket = {
  //       ticketId: 'mock',
  //       itemId: id,
  //       status: mockStatus,
  //       position: mockStatus === 'QUEUED' ? 3 : null,
  //       nextSlotFreeInSeconds: mockStatus === 'QUEUED' ? 120 : null,
  //       expiresInSeconds: mockStatus === 'OFFERED' || mockStatus === 'CHECKOUT' ? 120 : null,
  //       expiresAt: null,
  //       createdAt: new Date().toISOString(),
  //       serverTime: new Date().toISOString(),
  //     };
  //     setTicket(mockTicket);
  //     setLoading(false);
  //     return;
  //   }
  //   let interval: ReturnType<typeof setInterval>;
  //   const fetchStatus = async () => {
  //     try {
  //       const data = await getQueueStatus(id);
  //       setTicket(data);
  //       if (data.expiresInSeconds !== null) setTimer(data.expiresInSeconds);
  //       if (['PURCHASED', 'EXPIRED', 'SOLD_OUT', 'CANCELLED'].includes(data.status)) {
  //         clearInterval(interval);
  //       }
  //     } catch (error) {
  //       setTicket(null);
  //     } finally {
  //       setLoading(false);
  //     }
  //   };
  //   fetchStatus();
  //   interval = setInterval(fetchStatus, 2000);
  //   return () => clearInterval(interval);
  // }, [id, mockStatus]);

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
    if (!id) return;
    if (ticket && ['QUEUED', 'EXPIRED', 'CANCELLED', 'PURCHASED', 'SOLD_OUT'].includes(ticket.status)) {
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

  const productImage = product?.image_url || getProductImage(ticket.itemId);
  const renderProductCard = () => (
    <Grid size={{ xs: 12, md: 7 }}>
      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9' }}>
        <Box sx={{ display: 'flex', gap: 3 }}>
          <Box
            component="img"
            src={productImage}
            alt="Товар"
            onError={(e) => (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg'}
            sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }}
          />
          <Box>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              {product ? product.title : 'Загрузка...'}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Категория: {product?.category || 'Без категории'}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Цена: {product ? product.price.toLocaleString() : '—'} ₽
            </Typography>
          </Box>
        </Box>
      </Paper>
    </Grid>
  );
  const renderContent = () => {
    switch (ticket.status) {
      case 'QUEUED':
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>Очередь ожидания на лимитированный товар</Typography>
              <Grid container spacing={4}>
                {renderProductCard()}
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                    <Typography variant="h5" gutterBottom sx={{ fontWeight: 700 }}>
                      Вы {ticket.position ? `${ticket.position}-й в` : 'в'} очереди
                    </Typography>
                    <Typography variant="body1" sx={{ mb: 2 }}>
                      Ближайшее место по этому товару освободится не позже чем через{' '}
                      {ticket.nextSlotFreeInSeconds
                        ? formatTime(ticket.nextSlotFreeInSeconds)
                        : 'несколько секунд'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                      Как только товар освободится, вы получите право на покупку.
                    </Typography>
                    
                    <Button
                      variant="contained"
                      color="primary"
                      fullWidth
                      sx={{ py: 1.5, mb: 1 }}
                      onClick={() => navigate('/products')}
                    >
                      Перейти в каталог
                    </Button>
                    
                    <Button
                      variant="outlined"
                      color="secondary"
                      fullWidth
                      sx={{ py: 1.5 }}
                      onClick={async () => {
                        await cancelQueue(ticket.itemId);
                        navigate('/products');
                      }}
                    >
                      Выйти из очереди
                    </Button>
                  </Paper>
                </Grid>
              </Grid>
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
            </Container>
          </Box>
        );

      case 'OFFERED':
      case 'CHECKOUT': {
        const isCheckoutStarted = ticket.status === 'CHECKOUT';
        return (
          <Box>
            {renderHeader()}
            <Container maxWidth="lg" sx={{ py: 4 }}>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                Очередь ожидания на лимитированный товар
              </Typography>
              <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9', mb: 3 }}>
                    <Box sx={{ display: 'flex', gap: 3 }}>
                      <Box
                        component="img"
                        src={productImage}
                        alt="Товар"
                        onError={(e) => (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg'}
                        sx={{ width: 120, height: 120, borderRadius: 2, objectFit: 'cover' }}
                      />
                      <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>
                          {product ? product.title : 'Загрузка...'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Категория: {product?.category || 'Без категории'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Цена: {product ? product.price.toLocaleString() : '—'} ₽
                        </Typography>
                      </Box>
                    </Box>
                  </Paper>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f9f9f9' }}>
                    <Typography variant="body2" color="text.secondary">
                      <strong>Будьте внимательны:</strong> вы можете приобрести товар только в указанное время, иначе право будет передано следующему в очереди.
                    </Typography>
                  </Paper>
                </Grid>
                {/* Правая колонка */}
                <Grid size={{ xs: 12, md: 5 }}>
                  <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                    <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#00C853' }}>
                      {isCheckoutStarted ? 'Оформление начато' : 'Товар ваш!'}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', mb: 2 }}>
                      <CircularProgress
                        variant="determinate"
                        value={((timer || 0) / 60) * 100}
                        size={80}
                        sx={{ color: 'primary.main' }}
                      />
                      <Typography variant="h4" sx={{ fontWeight: 700, ml: 2 }}>
                        {formatTime(timer || 0)}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', mb: 2 }}>
                      {isCheckoutStarted
                        ? 'Товар держится за вами, пока идёт время.'
                        : 'У вас есть время, чтобы перейти к оформлению и оплатить заказ.'}
                    </Typography>
                    <Button
                      variant="contained"
                      color="primary"
                      fullWidth
                      sx={{ py: 1.5, mb: 1 }}
                      onClick={async () => {
                        if (isCheckoutStarted) {
                          navigate(`/product/${ticket.itemId}/checkout`);
                          return;
                        }
                        try {
                          await startCheckout(ticket.itemId);
                          navigate(`/product/${ticket.itemId}/checkout`);
                        } catch (error) {
                          console.error('Ошибка при переходе к оформлению:', error);
                          alert('Не удалось перейти к оформлению. Попробуйте ещё раз.');
                        }
                      }}
                    >
                      {isCheckoutStarted ? 'Продолжить оформление' : 'Перейти к оформлению'}
                    </Button>
                    <Button
                      variant="outlined"
                      color="secondary"
                      fullWidth
                      onClick={async () => {
                        await cancelQueue(ticket.itemId);
                        navigate('/products');
                      }}
                    >
                      Отказаться
                    </Button>
                  </Paper>
                </Grid>
              </Grid>
              <Box sx={{ mt: 4, textAlign: 'center' }}>
                <Typography variant="body2" color="text.secondary">
                  Похожие товары скрыты, чтобы вы не потеряли эту возможность.
                </Typography>
              </Box>
            </Container>
          </Box>
        );
      }
        case 'EXPIRED':
          return (
            <Box>
              {renderHeader()}
              <Container maxWidth="lg" sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom sx={{ fontWeight: 600, mb: 4 }}>
                  Время вышло
                </Typography>
                <Grid container spacing={4}>
                  {renderProductCard()}
                  <Grid size={{ xs: 12, md: 5 }}>
                    <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                      <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#FF6163' }}>
                        Время вышло
                      </Typography>
                      <Typography variant="body1" sx={{ mb: 2 }}>
                        Право на покупку истекло, товар вернулся в продажу и достался следующему в очереди.
                      </Typography>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        Вы можете встать в очередь заново, если товар ещё доступен.
                      </Typography>
                      <Button
                        variant="contained"
                        color="primary"
                        fullWidth
                        sx={{ py: 1.5 }}
                        onClick={() => navigate(`/product/${ticket.itemId}`)}
                      >
                        Встать в очередь заново
                      </Button>
                    </Paper>
                  </Grid>
                </Grid>
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
                    {renderProductCard()}
                    <Grid size={{ xs: 12, md: 5 }}>
                      <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                        <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#FF6163' }}>
                          Вы вышли из очереди
                        </Typography>
                        <Typography variant="body1" sx={{ mb: 2 }}>
                          Место передано следующему участнику. Можно встать заново — в конец очереди.
                        </Typography>
                        <Button
                          variant="contained"
                          color="primary"
                          fullWidth
                          sx={{ py: 1.5 }}
                          onClick={() => navigate(`/product/${ticket.itemId}`)}
                        >
                          Встать в очередь заново
                        </Button>
                      </Paper>
                    </Grid>
                  </Grid>
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
                      {renderProductCard()}
                      <Grid size={{ xs: 12, md: 5 }}>
                        <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                          <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#000000' }}>
                            Товар ваш!
                          </Typography>
                          <Typography variant="body1" sx={{ mb: 2 }}>
                            Заказ передан в доставку. Спасибо за покупку!
                          </Typography>
                          <Button
                            variant="contained"
                            fullWidth
                            sx={{ bgcolor: '#04e162', py: 1.5, mb: 1 }}
                            onClick={() => alert('Перейти к заказу')}
                          >
                            Перейти к заказу
                          </Button>
                          <Button
                            variant="outlined"
                            fullWidth
                            sx={{ py: 1.5 }}
                            onClick={() => navigate('/products')}
                          >
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
                            <img src={p.image_url || getProductImage(p.item_id)} alt={p.title} onError={(e) => {
                              (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
                            }} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
                        {renderProductCard()}
                        <Grid size={{ xs: 12, md: 5 }}>
                          <Paper elevation={0} sx={{ p: 3, borderRadius: 3, bgcolor: '#f5f5f5' }}>
                            <Typography variant="h5" gutterBottom sx={{ fontWeight: 700, color: '#FF6163' }}>
                              Товар раскуплен
                            </Typography>
                            <Typography variant="body1" sx={{ mb: 2 }}>
                              Все единицы выкуплены, пока вы ждали.
                            </Typography>
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                              Посмотрите похожие товары ниже.
                            </Typography>
                            <Button
                              variant="contained"
                              fullWidth
                              sx={{ bgcolor: '#00AAFF', py: 1.5 }}
                              onClick={() => navigate('/products')}
                            >
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
                              <img src={p.image_url || getProductImage(p.item_id)} alt={p.title} style={{ width: '100%', height: 100, objectFit: 'cover', borderRadius: 8 }} />
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
//src={`https://picsum.photos/seed/${p.item_id}/200/200`}