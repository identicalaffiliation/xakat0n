import React, { createContext, useContext, useState } from 'react';

export type QueueStatus = 
  | 'CHECKOUT'
  | 'QUEUED'
  | 'OFFERED'
  | 'EXPIRED'
  | 'CANCELLED'
  | 'PURCHASED';

interface QueueState {
  productId: number | null;
  status: QueueStatus | null;
  queuePosition: number | null;
  totalInQueue: number | null;
  timeLeft: number | null;
}

interface QueueContextType {
  state: QueueState;
  occupiedProductId: number | null;
  startCheckout: (productId: number) => void;
  joinQueue: (productId: number) => void;
  leaveQueue: () => void;
  expireOffer: () => void;
  confirmPurchase: () => void;
  reset: () => void;
  isProductOccupied: (productId: number) => boolean;
  forceStatus: (productId: number, newStatus: QueueStatus, timeLeft?: number) => void;
}

const QueueContext = createContext<QueueContextType | undefined>(undefined);

export const QueueProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, setState] = useState<QueueState>({
    productId: null,
    status: null,
    queuePosition: null,
    totalInQueue: null,
    timeLeft: null,
  });
  const [occupiedProductId, setOccupiedProductId] = useState<number | null>(null);
  const [queue, setQueue] = useState<number[]>([]);

  const startCheckout = (productId: number) => {
    setOccupiedProductId(productId);
    setState({
      productId,
      status: 'CHECKOUT',
      queuePosition: null,
      totalInQueue: null,
      timeLeft: 120,
    });
  };

  const joinQueue = (productId: number) => {
    const userId = 1;
    setQueue(prev => [...prev, userId]);
    const position = queue.length + 1;
    setOccupiedProductId(productId);
    setState({
      productId,
      status: 'QUEUED',
      queuePosition: position,
      totalInQueue: position + 2,
      timeLeft: 480,
    });
  };

  const leaveQueue = () => {
    const userId = 1;
    setQueue(prev => prev.filter(id => id !== userId));
    setState(prev => ({
      ...prev,
      status: 'CANCELLED',
      queuePosition: null,
      totalInQueue: null,
      timeLeft: null,
    }));
    setOccupiedProductId(null);
  };

  const expireOffer = () => {
    setState(prev => ({
      ...prev,
      status: 'EXPIRED',
      timeLeft: null,
    }));
    setOccupiedProductId(null);
  };

  const confirmPurchase = () => {
    setState(prev => ({
      ...prev,
      status: 'PURCHASED',
      timeLeft: null,
    }));
    setOccupiedProductId(null);
  };

  const reset = () => {
    setState({
      productId: null,
      status: null,
      queuePosition: null,
      totalInQueue: null,
      timeLeft: null,
    });
    setOccupiedProductId(null);
    setQueue([]);
  };

  const isProductOccupied = (productId: number) => occupiedProductId === productId;

  const forceStatus = (productId: number, newStatus: QueueStatus, timeLeft?: number) => {
    setState({
      productId,
      status: newStatus,
      queuePosition: newStatus === 'QUEUED' ? 2 : null,
      totalInQueue: newStatus === 'QUEUED' ? 4 : null,
      timeLeft: timeLeft || null,
    });
    if (newStatus === 'CHECKOUT' || newStatus === 'QUEUED') {
      setOccupiedProductId(productId);
    } else {
      setOccupiedProductId(null);
    }
  };

  return (
    <QueueContext.Provider value={{
      state,
      occupiedProductId,
      startCheckout,
      joinQueue,
      leaveQueue,
      expireOffer,
      confirmPurchase,
      reset,
      isProductOccupied,
      forceStatus,
    }}>
      {children}
    </QueueContext.Provider>
  );
};

export const useQueue = () => {
  const context = useContext(QueueContext);
  if (!context) throw new Error('useQueue must be used within QueueProvider');
  return context;
};