import React, { useState, useMemo } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  InputBase,
  Box,
  Chip,
  Card,
  CardMedia,
  CardContent,
  IconButton,
  Button,
  InputAdornment,
} from '@mui/material';
import { Search, Person, Close } from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { products } from '../../data/products';

const SearchWrapper = styled('div')({
  display: 'flex',
  alignItems: 'center',
  width: '100%',
  maxWidth: 900,
  border: '2px solid #00AAFF',
  borderRadius: '20px',
  backgroundColor: '#fff',
  overflow: 'hidden',
  height: 48,
});

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 1.5),
  height: '100%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: '#00AAFF',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  flex: 1,
  '& .MuiInputBase-input': {
    padding: theme.spacing(0.8, 1, 0.8, 0),
    paddingLeft: '8px',
    fontSize: '1rem',
    transition: theme.transitions.create('width'),
    width: '100%',
    height: '100%',
  },
}));

const SearchButton = styled(Button)(({ theme }) => ({
  backgroundColor: '#00AAFF',
  color: '#fff',
  borderRadius: 0,
  height: '100%',
  padding: theme.spacing(0, 3),
  textTransform: 'none',
  fontWeight: 400,
  fontSize: '0.95rem',
  minHeight: 48,
  '&:hover': {
    backgroundColor: '#0088cc',
  },
}));

const categories = ['Все', ...new Set(products.map((p) => p.category))];

const ProductCatalog: React.FC = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('Все');

  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const matchesCategory = selectedCategory === 'Все' || product.category === selectedCategory;
      const matchesSearch = product.title.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesCategory && matchesSearch;
    });
  }, [searchQuery, selectedCategory]);

  const handleProductClick = (id: number) => {
    navigate(`/product/${id}`);
  };

  const handleSearch = () => {
    console.log('Поиск:', searchQuery);
  };

  const handleClearSearch = () => {
    setSearchQuery('');
  };

  const getStockText = (stock: number) => {
    if (stock === 0) return 'Нет в наличии';
    if (stock === 1) return 'В наличии: один';
    return 'В наличии: несколько';
  };

  return (
    <Box sx={{ backgroundColor: '#ffffff', minHeight: '100vh' }}>
      <AppBar
        position="sticky"
        color="default"
        elevation={0}
        sx={{ backgroundColor: '#fff', top: 0, zIndex: 1100 }}
      >
        <Toolbar sx={{ justifyContent: 'space-between', py: 2, mt: 1 }}>
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              ml: 8,
              cursor: 'pointer',
              flexShrink: 0,
            }}
            onClick={() => navigate('/products')}
          >
            <img
              src="/avito-logo.svg"
              alt="Avito"
              style={{ width: 36, height: 36, marginRight: 10 }}
            />
            <Typography variant="h5" sx={{ fontWeight: 700, color: '#000000', fontSize: '1.8rem' }}>
              Avito
            </Typography>
          </Box>

          <SearchWrapper>
            <SearchIconWrapper>
              <Search />
            </SearchIconWrapper>
            <StyledInputBase
              placeholder="Поиск по объявлениям"
              inputProps={{ 'aria-label': 'search' }}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              endAdornment={
                searchQuery ? (
                  <InputAdornment position="end">
                    <IconButton
                      size="small"
                      onClick={handleClearSearch}
                      onMouseDown={(e) => e.preventDefault()}
                      sx={{ mr: 0.5, color: '#999', p: 0.5 }}
                    >
                      <Close fontSize="small" />
                    </IconButton>
                  </InputAdornment>
                ) : null
              }
            />
            <SearchButton onClick={handleSearch}>Найти</SearchButton>
          </SearchWrapper>

          <Box sx={{ flexShrink: 0, mr: 2 }}>
            <IconButton color="inherit" size="large">
              <Person fontSize="large" />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Box sx={{ maxWidth: 1400, mx: 'auto', px: 2, pt: 3, pb: 2 }}>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.8, mb: 1.5 }}>
          {categories.map((cat) => (
            <Chip
              key={cat}
              label={cat}
              onClick={() => setSelectedCategory(cat)}
              sx={{
                fontSize: '0.95rem',
                fontWeight: 500,
                borderRadius: 4,
                height: 32,
                ...(selectedCategory === cat
                  ? {
                      backgroundColor: '#ff5722',
                      color: '#fff',
                      border: 'none',
                    }
                  : {
                      backgroundColor: '#fff',
                      color: '#000',
                      border: '2px solid #00AAFF',
                    }),
              }}
            />
          ))}
        </Box>

        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: 'repeat(5, 1fr)',
            gap: 2.5,
          }}
        >
          {filteredProducts.map((product) => (
            <Box
              key={product.id}
              sx={{ cursor: 'pointer' }}
              onClick={() => handleProductClick(product.id)}
            >
              <Card
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  borderRadius: 4, 
                  boxShadow: 'none',
                  transition: '0.2s',
                  overflow: 'hidden',
                }}
              >
                <Box sx={{ position: 'relative' }}>
                  <CardMedia
                    component="img"
                    height="280"
                    image={product.image}
                    alt={product.title}
                    sx={{
                      borderRadius: 4,
                    }}
                  />
                  {product.stock === 1 && (
                    <Chip
                      label="Лимитированный"
                      sx={{
                        position: 'absolute',
                        top: 12,
                        left: 12,
                        backgroundColor: '#ff5722',
                        color: '#fff',
                        fontWeight: 600,
                        fontSize: '0.8rem',
                        borderRadius: 3,
                        height: 28,
                      }}
                    />
                  )}
                </Box>
                <CardContent sx={{ flexGrow: 1, p: 1.5, textAlign: 'left' }}>
                  <Typography
                    gutterBottom
                    variant="h6"
                    component="div"
                    sx={{ fontWeight: 400, fontSize: '1.1rem', lineHeight: 1.3 }}
                  >
                    {product.title}
                  </Typography>
                  <Typography
                    variant="body1"
                    sx={{ fontWeight: 700, color: '#000000', fontSize: '1.2rem' }}
                  >
                    {product.price.toLocaleString()} ₽
                  </Typography>
                  <Typography
                    variant="caption"
                    color="text.secondary"
                    sx={{ display: 'block', mt: 0.5, fontSize: '0.9rem' }}
                  >
                    {product.date}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: 500,
                      color: '#00C853',
                      mt: 0.5,
                      fontSize: '0.9rem',
                    }}
                  >
                    {getStockText(product.stock)}
                  </Typography>
                </CardContent>
              </Card>
            </Box>
          ))}
          {filteredProducts.length === 0 && (
            <Box sx={{ gridColumn: '1 / -1', textAlign: 'center', py: 4 }}>
              <Typography variant="body1" color="text.secondary">
                Ничего не найдено
              </Typography>
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default ProductCatalog;