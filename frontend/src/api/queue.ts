import client from './client';
import { USE_MOCK } from './mock';

export type QueueStatus =
  | 'QUEUED'
  | 'OFFERED'
  | 'CHECKOUT'
  | 'PURCHASED'
  | 'EXPIRED'
  | 'SOLD_OUT'
  | 'CANCELLED';

export type Ticket = {
  ticketId: string;
  itemId: string;
  status: QueueStatus;
  position: number | null;
  nextSlotFreeInSeconds: number | null;
  expiresInSeconds: number | null;
  expiresAt: string | null;
  createdAt: string;
  serverTime: string;
};

export type CheckoutStarted = {
  queueApplied: boolean;
  ticket: Ticket | null;
};

const mockTicket = (itemId: string, status: QueueStatus = 'OFFERED'): Ticket => ({
  ticketId: `mock-ticket-${Date.now()}`,
  itemId: itemId,
  status,
  position: status === 'QUEUED' ? 2 : null,
  nextSlotFreeInSeconds: status === 'QUEUED' ? 120 : null,
  expiresInSeconds: status === 'OFFERED' || status === 'CHECKOUT' ? 90 : null,
  expiresAt: status === 'OFFERED' || status === 'CHECKOUT' ? new Date(Date.now() + 90000).toISOString() : null,
  createdAt: new Date().toISOString(),
  serverTime: new Date().toISOString(),
});

export const enterQueue = async (itemId: string): Promise<Ticket> => {
  if (USE_MOCK) {
    const status: QueueStatus = Math.random() > 0.5 ? 'OFFERED' : 'QUEUED';
    return mockTicket(itemId, status);
  }
  const res = await client.post<Ticket>(`/items/${itemId}/queue`);
  return res.data;
};

export const getQueueStatus = async (itemId: string): Promise<Ticket> => {
  if (USE_MOCK) {
    return mockTicket(itemId, 'QUEUED');
  }
  const res = await client.get<Ticket>(`/items/${itemId}/queue/me`);
  return res.data;
};

export const cancelQueue = async (itemId: string): Promise<Ticket> => {
  if (USE_MOCK) {
    return mockTicket(itemId, 'CANCELLED');
  }
  const res = await client.delete<Ticket>(`/items/${itemId}/queue/me`);
  return res.data;
};

export const startCheckout = async (itemId: string): Promise<CheckoutStarted> => {
  if (USE_MOCK) {
    return {
      queueApplied: true,
      ticket: mockTicket(itemId, 'CHECKOUT'),
    };
  }
  const res = await client.post<CheckoutStarted>(`/items/${itemId}/checkout`);
  return res.data;
};

export const paymentCallback = async (itemId: string, result: 'paid' | 'failed'): Promise<Ticket> => {
  if (USE_MOCK) {
    return mockTicket(itemId, result === 'paid' ? 'PURCHASED' : 'CHECKOUT');
  }
  const res = await client.post<Ticket>(`/items/${itemId}/payment/callback`, { ticketId: 'mock', result });
  return res.data;
};