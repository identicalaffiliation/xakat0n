// import client from './client';

// export type QueueStatus =
//   | 'QUEUED'
//   | 'OFFERED'
//   | 'CHECKOUT'
//   | 'PURCHASED'
//   | 'EXPIRED'
//   | 'SOLD_OUT'
//   | 'CANCELLED';

// export type Ticket = {
//   ticket_id: string;
//   item_id: string;
//   status: QueueStatus;
//   position: number | null;
//   next_slot_free_in_seconds: number | null;
//   expires_in_seconds: number | null;
//   expires_at: string | null;
//   checkout_started_at: string | null;
//   payment_attempts: number;
//   created_at: string;
//   server_time: string;
// };

// export type CheckoutStarted = {
//   queue_applied: boolean;
//   ticket: Ticket | null;
// };

// export const enterQueue = async (itemId: string): Promise<Ticket> => {
//   const res = await client.post<Ticket>(`/items/${itemId}/queue`);
//   return res.data;
// };

// export const getQueueStatus = async (itemId: string): Promise<Ticket> => {
//   const res = await client.get<Ticket>(`/items/${itemId}/queue/me`);
//   return res.data;
// };

// export const cancelQueue = async (itemId: string): Promise<Ticket> => {
//   const res = await client.delete<Ticket>(`/items/${itemId}/queue/me`);
//   return res.data;
// };

// export const startCheckout = async (itemId: string): Promise<CheckoutStarted> => {
//   const res = await client.post<CheckoutStarted>(`/items/${itemId}/checkout`);
//   return res.data;
// };

// export const paymentCallback = async (itemId: string, result: 'paid' | 'failed'): Promise<Ticket> => {
//   const res = await client.post<Ticket>(`/items/${itemId}/payment/callback`, { result });
//   return res.data;
// };
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
  ticket_id: string;
  item_id: string;
  status: QueueStatus;
  position: number | null;
  next_slot_free_in_seconds: number | null;
  expires_in_seconds: number | null;
  expires_at: string | null;
  checkout_started_at: string | null;
  payment_attempts: number;
  created_at: string;
  server_time: string;
};

export type CheckoutStarted = {
  queue_applied: boolean;
  ticket: Ticket | null;
};

const mockTicket = (itemId: string, status: QueueStatus = 'OFFERED'): Ticket => ({
  ticket_id: `mock-ticket-${Date.now()}`,
  item_id: itemId,
  status,
  position: status === 'QUEUED' ? 2 : null,
  next_slot_free_in_seconds: status === 'QUEUED' ? 120 : null,
  expires_in_seconds: status === 'OFFERED' || status === 'CHECKOUT' ? 90 : null,
  expires_at: status === 'OFFERED' || status === 'CHECKOUT' ? new Date(Date.now() + 90000).toISOString() : null,
  checkout_started_at: status === 'CHECKOUT' ? new Date().toISOString() : null,
  payment_attempts: 0,
  created_at: new Date().toISOString(),
  server_time: new Date().toISOString(),
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
      queue_applied: true,
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
  const res = await client.post<Ticket>(`/items/${itemId}/payment/callback`, { result });
  return res.data;
};