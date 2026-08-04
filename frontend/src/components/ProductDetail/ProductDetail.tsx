import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Typography, Button, Paper, Container } from '@mui/material';
import { products } from '../../data/products';

const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const product = products.find((p) => p.id === Number(id));

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

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Button
        variant="outlined"
        onClick={() => navigate('/products')}
        sx={{ mb: 3, borderRadius: 2 }}
      >
        ← Назад
      </Button>
      <Paper
        elevation={0}
        sx={{
          p: 4,
          borderRadius: 2,
          boxShadow: 'none',
          backgroundColor: '#fff',
        }}
      >
        <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 4 }}>
          <Box
            component="img"
            src={product.image}
            alt={product.title}
            sx={{
              width: '100%',
              maxWidth: 400,
              height: 'auto',
              borderRadius: 2,
            }}
          />
          <Box>
            <Typography variant="h4" gutterBottom>
              {product.title}
            </Typography>
            <Typography variant="h5" color="#00AAFF" sx={{ fontWeight: 700, mb: 2 }}>
              {product.price.toLocaleString()} ₽
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Категория: {product.category}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Опубликовано: {product.date}
            </Typography>
            <Button
              variant="contained"
              size="large"
              sx={{
                mt: 3,
                backgroundColor: '#00AAFF',
                borderRadius: 2,
                '&:hover': { backgroundColor: '#00AAFF' },
              }}
            >
              Купить
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default ProductDetail;