import React, { useState, useMemo } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  InputBase,
  Box,
  Chip,
  Grid,
  Card,
  CardMedia,
  CardContent,
  IconButton,
} from '@mui/material';
import { Search, Person } from '@mui/icons-material';
import { styled, alpha } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import type { Product } from '../../data/products';
import { products } from '../../data/products';

const SearchWrapper = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.black, 0.06),
  '&:hover': {
    backgroundColor: alpha(theme.palette.common.black, 0.1),
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: '100%',
  maxWidth: 500,
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: '100%',
  position: 'absolute',
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  width: '100%',
  '& .MuiInputBase-input': {
    padding: theme.spacing(1, 1, 1, 0),
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
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

  return (
    <Box sx={{ backgroundColor: '#f5f5f5', minHeight: '100vh' }}>
      <AppBar position="static" color="default" elevation={1} sx={{ backgroundColor: '#fff' }}>
        <Toolbar>
          <Typography
            variant="h6"
            sx={{ fontWeight: 700, color: '#00A2E3', mr: 3, cursor: 'pointer' }}
            onClick={() => navigate('/products')}
          >
            Авито
          </Typography>
          <SearchWrapper>
            <SearchIconWrapper>
              <Search />
            </SearchIconWrapper>
            <StyledInputBase
              placeholder="Поиск по объявлениям"
              inputProps={{ 'aria-label': 'search' }}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </SearchWrapper>
          <Box sx={{ flexGrow: 1 }} />
          <IconButton color="inherit" sx={{ mr: 1 }}>
            <Person />
          </IconButton>
        </Toolbar>
      </AppBar>

      <Box sx={{ maxWidth: 1200, mx: 'auto', px: 3, py: 2 }}>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 3 }}>
          {categories.map((cat) => (
            <Chip
              key={cat}
              label={cat}
              onClick={() => setSelectedCategory(cat)}
              color={selectedCategory === cat ? 'primary' : 'default'}
              variant={selectedCategory === cat ? 'filled' : 'outlined'}
              sx={{ fontSize: '0.875rem', fontWeight: 500 }}
            />
          ))}
        </Box>

        <Grid container spacing={3}>
          {filteredProducts.map((product) => (
            <Grid
              size={{ xs: 12, sm: 6, md: 4, lg: 3 }}
              key={product.id}
            >
              <Card
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  cursor: 'pointer',
                  transition: '0.2s',
                  '&:hover': { boxShadow: 6 },
                }}
                onClick={() => handleProductClick(product.id)}
              >
                <CardMedia
                  component="img"
                  height="200"
                  image={product.image}
                  alt={product.title}
                />
                <CardContent sx={{ flexGrow: 1 }}>
                  <Typography
                    gutterBottom
                    variant="h6"
                    component="div"
                    sx={{ fontWeight: 600, fontSize: '1rem' }}
                  >
                    {product.title}
                  </Typography>
                  <Typography
                    variant="body1"
                    sx={{ fontWeight: 700, color: '#00A2E3' }}
                  >
                    {product.price.toLocaleString()} ₽
                  </Typography>
                  <Typography
                    variant="caption"
                    color="text.secondary"
                    sx={{ display: 'block', mt: 1 }}
                  >
                    {product.date}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          ))}
          {filteredProducts.length === 0 && (
            <Box sx={{ width: '100%', textAlign: 'center', py: 4 }}>
              <Typography variant="body1" color="text.secondary">
                Ничего не найдено
              </Typography>
            </Box>
          )}
        </Grid>
      </Box>
    </Box>
  );
};

export default ProductCatalog;