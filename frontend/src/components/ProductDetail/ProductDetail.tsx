import React, { useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Paper,
  Container,
  AppBar,
  Toolbar,
  IconButton,
  Chip,
  Autocomplete,
  TextField,
} from '@mui/material';
import { FavoriteBorder, Favorite, Person, Search } from '@mui/icons-material';
import { products, type Product } from '../../data/products';
import { useQueue } from '../../context/QueueContext';

const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const product = products.find((p) => p.id === Number(id));
  const [isFavorite, setIsFavorite] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState('');
  const { state, startCheckout, joinQueue, isProductOccupied } = useQueue();
  const filteredProducts = useMemo(() => {
    if (!searchQuery.trim()) return [];
    return products.filter((p) =>
      p.title.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [searchQuery]);

  if (!product) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography variant="h5">Товар не найден</Typography>
        <Button variant="outlined" onClick={() => navigate('/products')}>
          Вернуться к каталогу
        </Button>
      </Box>
    );
  }

  const toggleFavorite = () => {
    setIsFavorite(!isFavorite);
  };

  const getStockDisplay = (stock: number) => {
    if (stock === 0) {
      return { prefix: 'В наличии: ', suffix: 'нет в наличии' };
    }
    if (stock === 1) {
      return { prefix: 'В наличии: ', suffix: 'один' };
    }
    return { prefix: 'В наличии: ', suffix: 'несколько' };
  };

  const stockInfo = getStockDisplay(product.stock);
  const handleProductSelect = (_event: any, value: string | Product | null) => {
    if (value && typeof value !== 'string') {
      navigate(`/product/${value.id}`);
    }
  };
  const handleSearch = () => {
    if (filteredProducts.length > 0) {
      navigate(`/product/${filteredProducts[0].id}`);
    }
  };
  // const handleBuy = () => {
  //   if (!product) return;
  //   if (product.is_limited) {
  //     if (isProductOccupied(product.id)) {
  //       if (state.productId === product.id && state.status !== null && state.status !== 'EXPIRED' && state.status !== 'CANCELLED') {
  //         navigate(`/product/${product.id}/queue`);
  //         return;
  //       }
  //       joinQueue(product.id);
  //       navigate(`/product/${product.id}/queue`);
  //     } else {
  //       startCheckout(product.id);
  //       navigate(`/product/${product.id}/queue`);
  //     }
  //   } else {
  //     alert('Обычный товар – переход к оформлению');
  //   }
  // };
  const handleBuy = () => {
    if (!product) return;
    if (product.is_limited) {
      // Всегда встаём в очередь (для теста)
      joinQueue(product.id);
      navigate(`/product/${product.id}/queue`);
    } else {
      alert('Обычный товар – переход к оформлению');
    }
  };
  return (
    <Box sx={{ backgroundColor: '#ffffff', minHeight: '100vh' }}>
      <AppBar
        position="sticky"
        color="default"
        elevation={0}
        sx={{ backgroundColor: '#fff', top: 0, zIndex: 1100 }}
      >
        <Toolbar sx={{ justifyContent: 'space-between', py: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', ml: 4, flex: 1 }}>
            {/* Логотип + кнопка "Назад" */}
            <Box
              sx={{ display: 'flex', alignItems: 'center', cursor: 'pointer', mr: 2 }}
              onClick={() => navigate('/products')}
            >
              <img
                src="/avito-logo.svg"
                alt="Avito"
                style={{ width: 36, height: 36, marginRight: 10 }}
              />
              <Typography
                variant="h5"
                sx={{ fontWeight: 700, color: '#000000', fontSize: '1.8rem' }}
              >
                Avito
              </Typography>
            </Box>

            <Button
              variant="contained"
              onClick={() => navigate('/products')}
              sx={{
                backgroundColor: '#00AAFF',
                color: '#fff',
                borderRadius: 4,
                textTransform: 'none',
                fontWeight: 400,
                fontSize: '1rem',
                px: 4,
                py: 1,
                boxShadow: 'none',
                '&:hover': {
                  backgroundColor: '#0088cc',
                  boxShadow: 'none',
                },
              }}
            >
              Назад
            </Button>

            {/* Поисковик */}
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
                  <Box
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      width: '100%',
                      maxWidth: 900,
                      border: '2px solid #00AAFF',
                      borderRadius: '20px',
                      overflow: 'hidden',
                      backgroundColor: '#fff',
                      height: 48,
                    }}
                  >
                    <Box sx={{ px: 1.5, color: '#00AAFF' }}>
                      <Search />
                    </Box>
                    <TextField
                      {...params}
                      placeholder="Поиск товара..."
                      variant="outlined"
                      size="small"
                      sx={{
                        flex: 1,
                        '& .MuiOutlinedInput-root': {
                          border: 'none',
                          borderRadius: 0,
                          height: 48,
                          '& fieldset': { border: 'none' },
                        },
                        '& .MuiInputBase-input': {
                          padding: '8px 0',
                          fontSize: '1rem',
                        },
                      }}
                      onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                    />
                    <Button
                      onClick={handleSearch}
                      sx={{
                        backgroundColor: '#00AAFF',
                        color: '#fff',
                        borderRadius: 0,
                        height: 48,
                        px: 3,
                        textTransform: 'none',
                        fontWeight: 400,
                        fontSize: '0.95rem',
                        '&:hover': { backgroundColor: '#0088cc' },
                      }}
                    >
                      Найти
                    </Button>
                  </Box>
                )}
                sx={{ width: '100%' }}
              />
            </Box>
          </Box>

          {/* Иконка профиля */}
          <Box sx={{ flexShrink: 0, mr: 2 }}>
            <IconButton color="inherit" size="large">
              <Person fontSize="large" />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Paper elevation={0} sx={{ p: 3, borderRadius: 4 }}>
          <Box
            sx={{
              display: 'flex',
              flexDirection: { xs: 'column', md: 'row' },
              gap: 4,
            }}
          >
            <Box sx={{ flex: '0 0 60%', position: 'relative' }}>
              <Box
                component="img"
                src={product.image}
                alt={product.title}
                sx={{
                  width: '100%',
                  height: 'auto',
                  borderRadius: 4,
                  maxHeight: 600,
                  objectFit: 'cover',
                }}
              />
              {product.is_limited && (
                <Chip
                  label="Лимитированный"
                  sx={{
                    position: 'absolute',
                    top: 16,
                    left: 16,
                    backgroundColor: '#FF6163',
                    color: '#fff',
                    fontWeight: 600,
                    fontSize: '1rem',
                    borderRadius: 3,
                    px: 2,
                    py: 1,
                  }}
                />
              )}
            </Box>

            <Box
              sx={{
                flex: 1,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Box
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'flex-start',
                  }}
                >
                  <Typography variant="h4" sx={{ fontWeight: 400, mb: 1 }}>
                    {product.title}
                  </Typography>
                  <IconButton
                    onClick={toggleFavorite}
                    sx={{ color: isFavorite ? '#FF6163' : '#999' }}
                  >
                    {isFavorite ? <Favorite /> : <FavoriteBorder />}
                  </IconButton>
                </Box>

                <Typography
                  variant="h3"
                  sx={{ fontWeight: 700, color: '#000000', mb: 1 }}
                >
                  {product.price.toLocaleString()} ₽
                </Typography>

                <Typography variant="body1" sx={{ mb: 2, fontSize: '1.1rem' }}>
                  <Box component="span" sx={{ fontWeight: 500, color: '#000000' }}>
                    {stockInfo.prefix}
                  </Box>
                  <Box
                    component="span"
                    sx={{ fontWeight: 600, color: '#00C853' }}
                  >
                    {stockInfo.suffix}
                  </Box>
                </Typography>

                <Typography
                  variant="body1"
                  color="text.secondary"
                  sx={{ mb: 0.5 }}
                >
                  Категория: {product.category}
                </Typography>
                <Typography
                  variant="body2"
                  color="text.secondary"
                  sx={{ mb: 3 }}
                >
                  Опубликовано: {product.date}
                </Typography>
              </Box>

              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
                <Button
                  variant="contained"
                  size="large"
                  sx={{
                    backgroundColor: '#00AAFF',
                    color: '#fff',
                    borderRadius: 4,
                    textTransform: 'none',
                    fontWeight: 600,
                    fontSize: '1.2rem',
                    py: 1.5,
                    '&:hover': { backgroundColor: '#0088cc' },
                  }}
                  onClick={handleBuy}
                >
                  Купить
                </Button>
                <Button
                  variant="outlined"
                  size="large"
                  sx={{
                    borderColor: '#00AAFF',
                    color: '#00AAFF',
                    borderRadius: 4,
                    textTransform: 'none',
                    fontWeight: 600,
                    fontSize: '1.2rem',
                    py: 1.5,
                    '&:hover': {
                      borderColor: '#0088cc',
                      backgroundColor: 'rgba(0,170,255,0.04)',
                    },
                  }}
                >
                  Добавить в корзину
                </Button>
              </Box>
            </Box>
          </Box>
        </Paper>
      </Container>
    </Box>
  );
};

export default ProductDetail;