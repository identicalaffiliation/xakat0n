import React, { useState, useMemo, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box, Typography, Button, Paper, Container, AppBar, Toolbar, IconButton,
  Chip, Autocomplete, TextField,
} from '@mui/material';
import { FavoriteBorder, Favorite, Person, Search } from '@mui/icons-material';
import { getItem, type Item } from '../../api/items';
import { enterQueue, startCheckout } from '../../api/queue';
// import { useQueue } from '../../context/QueueContext';
import { products, type Product } from '../../data/products';

const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Item | null>(null);
  const [isFavorite, setIsFavorite] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  // const { state, forceStatus } = useQueue();

  useEffect(() => {
    if (!id) return;
    getItem(id)
      .then(data => setProduct(data))
      .catch(console.error);
  }, [id]);

  const filteredProducts = useMemo(() => {
    if (!searchQuery.trim()) return [];
    return products.filter(p => p.title.toLowerCase().includes(searchQuery.toLowerCase()));
  }, [searchQuery]);

  const handleProductSelect = (_event: any, value: string | Product | null) => {
    if (value && typeof value !== 'string') navigate(`/product/${value.id}`);
  };
  const handleSearch = () => {
    if (filteredProducts.length > 0) navigate(`/product/${filteredProducts[0].id}`);
  };

  const handleBuy = async () => {
    if (!product) return;
    if (product.is_limited) {
      try {
        await enterQueue(product.item_id);
        navigate(`/product/${product.item_id}/queue`);
      } catch (error) {
        console.error('Ошибка при входе в очередь:', error);
        alert('Не удалось встать в очередь. Попробуйте позже.');
      }
    } else {
      try {
        await startCheckout(product.item_id);
        // navigate(`/product/${product.item_id}/queue`);
        navigate(`/product/${product.item_id}/checkout`);
      } catch (error) {
        console.error('Ошибка при переходе к оформлению:', error);
        alert('Не удалось перейти к оформлению.');
      }
    }
  };

  if (!product) return <Typography>Загрузка...</Typography>;

  const stockInfo = product.stock !== undefined && product.stock !== null
    ? { prefix: 'В наличии: ', suffix: product.stock === 0 ? 'нет в наличии' : product.stock === 1 ? 'один' : 'несколько' }
    : { prefix: '', suffix: '' };

  return (
    <Box sx={{ bgcolor: '#fff', minHeight: '100vh' }}>
      <AppBar position="sticky" color="default" elevation={0} sx={{ bgcolor: '#fff', top: 0, zIndex: 1100 }}>
        <Toolbar sx={{ justifyContent: 'space-between', py: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', ml: 4, flex: 1 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', cursor: 'pointer', mr: 2 }} onClick={() => navigate('/products')}>
              <img src="/avito-logo.svg" alt="Avito" style={{ width: 36, height: 36, marginRight: 10 }} />
              <Typography variant="h5" sx={{ fontWeight: 700, color: '#000', fontSize: '1.8rem' }}>Avito</Typography>
            </Box>
            <Button variant="contained" onClick={() => navigate('/products')} sx={{ bgcolor: '#00AAFF', color: '#fff', borderRadius: 4, textTransform: 'none', fontWeight: 400, px: 4, py: 1, boxShadow: 'none', '&:hover': { bgcolor: '#0088cc' } }}>
              Назад
            </Button>
            <Box sx={{ flex: 1, display: 'flex', justifyContent: 'center', mx: 2 }}>
              <Autocomplete
                options={filteredProducts}
                getOptionLabel={(option) => (typeof option === 'string' ? option : option.title)}
                value={null}
                onChange={handleProductSelect}
                inputValue={searchQuery}
                onInputChange={(_event, newValue) => setSearchQuery(newValue)}
                freeSolo
                renderInput={(params) => (
                  <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', maxWidth: 900, border: '2px solid #00AAFF', borderRadius: '20px', overflow: 'hidden', bgcolor: '#fff', height: 48 }}>
                    <Box sx={{ px: 1.5, color: '#00AAFF' }}><Search /></Box>
                    <TextField {...params} placeholder="Поиск товара..." variant="outlined" size="small" sx={{ flex: 1, '& .MuiOutlinedInput-root': { border: 'none', borderRadius: 0, height: 48, '& fieldset': { border: 'none' } }, '& .MuiInputBase-input': { padding: '8px 0', fontSize: '1rem' } }} onKeyDown={(e) => e.key === 'Enter' && handleSearch()} />
                    <Button onClick={handleSearch} sx={{ bgcolor: '#00AAFF', color: '#fff', borderRadius: 0, height: 48, px: 3, textTransform: 'none', fontWeight: 400, '&:hover': { bgcolor: '#0088cc' } }}>Найти</Button>
                  </Box>
                )}
                sx={{ width: '100%' }}
              />
            </Box>
          </Box>
          <Box sx={{ flexShrink: 0, mr: 2 }}>
            <IconButton color="inherit" size="large"><Person fontSize="large" /></IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Paper elevation={0} sx={{ p: 3, borderRadius: 4 }}>
          <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 4 }}>
            <Box sx={{ flex: '0 0 60%', position: 'relative' }}>
              <Box component="img" src={`https://picsum.photos/seed/${product.item_id}/600/400`} alt={product.title} sx={{ width: '100%', height: 'auto', borderRadius: 4, maxHeight: 600, objectFit: 'cover' }} />
              {product.is_limited && (
                <Chip label="Лимитированный" sx={{ position: 'absolute', top: 16, left: 16, bgcolor: '#FF6163', color: '#fff', fontWeight: 600, fontSize: '1rem', borderRadius: 3, px: 2, py: 1 }} />
              )}
            </Box>
            <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Typography variant="h4" sx={{ fontWeight: 400, mb: 1 }}>{product.title}</Typography>
                  <IconButton onClick={() => setIsFavorite(!isFavorite)} sx={{ color: isFavorite ? '#FF6163' : '#999' }}>
                    {isFavorite ? <Favorite /> : <FavoriteBorder />}
                  </IconButton>
                </Box>
                <Typography variant="h3" sx={{ fontWeight: 700, color: '#000', mb: 1 }}>{product.price.toLocaleString()} ₽</Typography>
                {stockInfo.prefix && (
                  <Typography variant="body1" sx={{ mb: 2, fontSize: '1.1rem' }}>
                    <Box component="span" sx={{ fontWeight: 500, color: '#000' }}>{stockInfo.prefix}</Box>
                    <Box component="span" sx={{ fontWeight: 600, color: '#00C853' }}>{stockInfo.suffix}</Box>
                  </Typography>
                )}
                <Typography variant="body1" color="text.secondary" sx={{ mb: 0.5 }}>Категория: {product.category || 'Без категории'}</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>Опубликовано: сегодня</Typography>
              </Box>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
                <Button variant="contained" size="large" sx={{ bgcolor: '#00AAFF', color: '#fff', borderRadius: 4, textTransform: 'none', fontWeight: 600, fontSize: '1.2rem', py: 1.5, '&:hover': { bgcolor: '#0088cc' } }} onClick={handleBuy}>
                  Купить
                </Button>
                <Button variant="outlined" size="large" sx={{ borderColor: '#00AAFF', color: '#00AAFF', borderRadius: 4, textTransform: 'none', fontWeight: 600, fontSize: '1.2rem', py: 1.5, '&:hover': { borderColor: '#0088cc', bgcolor: 'rgba(0,170,255,0.04)' } }}>
                  Добавить в корзину
                </Button>
              </Box>
            </Box>
          </Box>
          {/* {import.meta.env.DEV && (
            <Box sx={{ mt: 4, p: 2, border: '1px dashed #ccc', borderRadius: 2 }}>
              <Typography variant="subtitle2">Тестовый режим</Typography>
              <FormControl fullWidth size="small" sx={{ mt: 1 }}>
                <InputLabel>Статус товара</InputLabel>
                <Select
                  value={state.status || ''}d
                  onChange={(e) => {
                    const status = e.target.value as QueueStatus;
                    if (product) {
                      forceStatus(
                        Number(product.item_id), 
                        status as any,
                        status === 'QUEUED' ? 480 : status === 'CHECKOUT' ? 120 : status === 'OFFERED' ? 60 : undefined
                      );
                      navigate(`/product/${product.item_id}/queue`);
                    }
                  }}
                  label="Статус товара"
                >
                  <MenuItem value="CHECKOUT">Оформление (CHECKOUT)</MenuItem>
                  <MenuItem value="QUEUED">Очередь (QUEUED)</MenuItem>
                  <MenuItem value="OFFERED">Право выдано (OFFERED)</MenuItem>
                  <MenuItem value="EXPIRED">Время вышло (EXPIRED)</MenuItem>
                </Select>
              </FormControl>
            </Box>
          )} */}
        </Paper>
      </Container>
    </Box>
  );
};

export default ProductDetail;