import React, { useState, useMemo, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AppBar, Toolbar, Typography, InputBase, Box, Chip, Card, CardMedia, CardContent,
  IconButton, Button, InputAdornment
} from '@mui/material';
import { Search, Person, Close } from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { getItems, type Item } from '../../api/items';
const SearchWrapper = styled('div')({
  display: 'flex', alignItems: 'center', width: '100%', maxWidth: 900,
  border: '2px solid #00AAFF', borderRadius: '20px', bgcolor: '#fff',
  overflow: 'hidden', height: 48,
});
import { getProductImage } from '../../utils/imageUtils';
const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 1.5), height: '100%', display: 'flex',
  alignItems: 'center', justifyContent: 'center', color: '#00AAFF',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit', flex: 1,
  '& .MuiInputBase-input': {
    padding: theme.spacing(0.8, 1, 0.8, 0), paddingLeft: '8px',
    fontSize: '1rem', transition: theme.transitions.create('width'),
    width: '100%', height: '100%',
  },
}));

const SearchButton = styled(Button)(({ theme }) => ({
  backgroundColor: '#00AAFF', color: '#fff', borderRadius: 0,
  height: '100%', padding: theme.spacing(0, 3), textTransform: 'none',
  fontWeight: 400, fontSize: '0.95rem', minHeight: 48,
  '&:hover': { backgroundColor: '#0088cc' },
}));

const ProductCatalog: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('Все');

  // useEffect(() => {
  //   getItems()
  //     .then(data => setItems(data || []))
  //     .catch(console.error)
  // }, []);

  const categories = useMemo(() => {
    const cats = items.map(p => p.category).filter(Boolean) as string[];
    return ['Все', ...new Set(cats)];
  }, [items]);


  // useEffect(() => {
  //   setItems(mockItems);
  // }, []);
  useEffect(() => {
    getItems()
      .then(data => {
        console.log('Товары из API:', data);
        setItems(data);
      })
      .catch(err => console.error('Ошибка загрузки товаров:', err));
  }, []);

  // const filteredItems = useMemo(() => {
  //   return items.filter(item => {
  //     if (!item || !item.title) return false;
  //     const matchCat = selectedCategory === 'Все' || item.category === selectedCategory;
  //     const matchSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase());
  //     return matchCat && matchSearch;
  //   });
  // }, [items, searchQuery, selectedCategory]);
  const filteredItems = useMemo(() => {
    console.log('Товары для фильтрации:', items);
    return items.filter(item => {
      if (!item || !item.title) return false;
      const matchCat = selectedCategory === 'Все' || item.category === selectedCategory;
      const matchSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase());
      return matchCat && matchSearch;
    });
  }, [items, searchQuery, selectedCategory]);
  const getStockText = (stock?: number | null) => {
    if (stock === 0) return 'Нет в наличии';
    if (stock === 1) return 'В наличии: один';
    return 'В наличии: несколько';
  };

  const handleProductClick = (itemId: string) => navigate(`/product/${itemId}`);
  const handleClearSearch = () => setSearchQuery('');

  return (
    <Box sx={{ bgcolor: '#fff', minHeight: '100vh' }}>
      <AppBar position="sticky" color="default" elevation={0} sx={{ bgcolor: '#fff', top: 0, zIndex: 1100 }}>
        <Toolbar sx={{ justifyContent: 'space-between', py: 2, mt: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', ml: 8, cursor: 'pointer', flexShrink: 0 }} onClick={() => navigate('/products')}>
            <img src="/avito-logo.svg" alt="Avito" style={{ width: 36, height: 36, marginRight: 10 }} />
            <Typography variant="h5" sx={{ fontWeight: 700, color: '#000', fontSize: '1.8rem' }}>Avito</Typography>
          </Box>
          <SearchWrapper>
            <SearchIconWrapper><Search /></SearchIconWrapper>
            <StyledInputBase
              placeholder="Поиск по объявлениям"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              endAdornment={searchQuery ? (
                <InputAdornment position="end">
                  <IconButton size="small" onClick={handleClearSearch} onMouseDown={(e) => e.preventDefault()} sx={{ mr: 0.5, color: '#999', p: 0.5 }}>
                    <Close fontSize="small" />
                  </IconButton>
                </InputAdornment>
              ) : null}
            />
            <SearchButton onClick={() => {}}>Найти</SearchButton>
          </SearchWrapper>
          <Box sx={{ flexShrink: 0, mr: 2 }}>
            <IconButton color="inherit" size="large"><Person fontSize="large" /></IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Box sx={{ maxWidth: 1400, mx: 'auto', px: 2, pt: 3, pb: 2 }}>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.8, mb: 1.5 }}>
          {categories.map(cat => (
            <Chip
              key={cat}
              label={cat}
              onClick={() => setSelectedCategory(cat)}
              sx={{
                fontSize: '0.95rem', fontWeight: 500, borderRadius: 4, height: 32,
                ...(selectedCategory === cat
                  ? { bgcolor: '#FF6163', color: '#fff', border: 'none' }
                  : { bgcolor: '#fff', color: '#000', border: '2px solid #00AAFF' }),
              }}
            />
          ))}
        </Box>

        <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(5, 1fr)', gap: 2.5 }}>
          {filteredItems.map(item => (
            <Box key={item.item_id} sx={{ cursor: 'pointer' }} onClick={() => handleProductClick(item.item_id)}>
              <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column', borderRadius: 4, boxShadow: 'none', overflow: 'hidden' }}>
                <Box sx={{ position: 'relative' }}>
                  <CardMedia
                    component="img"
                    height="280"
                    image={getProductImage(item.item_id)}
                    // image={`https://picsum.photos/seed/${item.item_id}/200/200`}
                    alt={item.title}
                    onError={(e) => {
                      (e.target as HTMLImageElement).src = '/images/products/placeholder.jpg';
                    }}
                    sx={{ borderRadius: 4 }}
                  />
                  {item.is_limited && (
                    <Chip
                      label="Лимитированный"
                      sx={{ position: 'absolute', top: 12, left: 12, bgcolor: '#FF6163', color: '#fff', fontWeight: 600, fontSize: '0.8rem', borderRadius: 3, height: 28 }}
                    />
                  )}
                </Box>
                <CardContent sx={{ flexGrow: 1, p: 1.5, textAlign: 'left' }}>
                  <Typography gutterBottom variant="h6" component="div" sx={{ fontWeight: 400, fontSize: '1.1rem', lineHeight: 1.3 }}>
                    {item.title}
                  </Typography>
                  <Typography variant="body1" sx={{ fontWeight: 700, color: '#000', fontSize: '1.2rem' }}>
                    {item.price.toLocaleString()} ₽
                  </Typography>
                  <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 0.5, fontSize: '0.9rem' }}>
                    {item.category || 'Без категории'}
                  </Typography>
                  <Typography variant="body2" sx={{ fontWeight: 500, color: '#00C853', mt: 0.5, fontSize: '0.9rem' }}>
                    {getStockText(item.stock)}
                  </Typography>
                </CardContent>
              </Card>
            </Box>
          ))}
          {filteredItems.length === 0 && (
            <Box sx={{ gridColumn: '1 / -1', textAlign: 'center', py: 4 }}>
              <Typography variant="body1" color="text.secondary">Ничего не найдено</Typography>
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default ProductCatalog;