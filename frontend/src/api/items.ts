// import client from './client';

// export type Item = {
//   item_id: string;
//   title: string;
//   description?: string;
//   price: number;
//   category?: string;
//   is_limited: boolean;
//   stock?: number | null;
//   sold_out: boolean;
// };

// export const getItems = async (): Promise<Item[]> => {
//   const res = await client.get<{ items: Item[] | null }>('/items');
//   return res.data.items ?? [];
// };

// export const getItem = async (itemId: string): Promise<Item> => {
//   const res = await client.get<Item>(`/items/${itemId}`);
//   return res.data;
// };

// export const getSimilar = async (itemId: string, limit = 6): Promise<Item[]> => {
//   const res = await client.get<Item[]>(`/items/${itemId}/similar`, { params: { limit } });
//   return res.data;
// };
import client from './client';
import { USE_MOCK } from './mock';
import { mockItems } from './mockData';

export type Item = {
  item_id: string;
  title: string;
  description?: string;
  price: number;
  category?: string;
  is_limited: boolean;
  stock?: number | null;
  sold_out: boolean;
};

const mapItem = (raw: any): Item => {
  const data = raw.item || raw; 
  return {
    item_id: data.id,
    title: data.title,
    description: data.description || '',
    price: data.price,
    category: data.category || undefined,
    is_limited: data.isLimited ?? false,
    stock: data.stock ?? null,
    sold_out: (data.stock ?? 0) === 0,
  };
};

export const getItems = async (): Promise<Item[]> => {
  if (USE_MOCK) return mockItems;
  const res = await client.get<{ items: any[] }>('/items');
  return (res.data.items || []).map(mapItem);
};

export const getItem = async (itemId: string): Promise<Item> => {
  if (USE_MOCK) {
    const item = mockItems.find(i => i.item_id === itemId);
    if (!item) throw new Error('Item not found');
    return item;
  }
  const res = await client.get<any>(`/items/${itemId}`);
  return mapItem(res.data);
};

export const getSimilar = async (itemId: string, limit = 6): Promise<Item[]> => {
  if (USE_MOCK) {
    return mockItems.filter(i => i.item_id !== itemId).slice(0, limit);
  }
  const res = await client.get<any[]>(`/items/${itemId}/similar`, { params: { limit } });
  return (res.data || []).map(mapItem);
};