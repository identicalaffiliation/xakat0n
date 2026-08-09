import { products } from '../data/products';
export const mockItems = products.map(p => ({
  item_id: String(p.id),
  title: p.title,
  description: p.description || '',
  price: p.price,
  category: p.category,
  is_limited: p.is_limited || false,
  stock: p.stock,
  sold_out: p.stock === 0,
}));

export const mockLoginResponse = (username: string) => ({
  userId: `mock-user-${Date.now()}`,
  username,
  token: 'mock-token-for-testing',
});