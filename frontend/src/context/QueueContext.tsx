import React, { createContext, useContext, useState } from 'react';

export type QueueStatus = 
  | 'CHECKOUT'
  | 'QUEUED'
  | 'OFFERED'
  | 'EXPIRED'
  | 'CANCELLED'
  | 'PURCHASED';

interface QueueState {
  itemId: number | null;
  status: QueueStatus | null;
  queuePosition: number | null;
  totalInQueue: number | null;
  timeLeft: number | null;
  isLimited: boolean | null;
}

interface QueueContextType {
  state: QueueState;
  occupiedItemId: number | null;
  startCheckout: (itemId: number, isLimited?: boolean) => void;
  joinQueue: (itemId: number) => void;
  leaveQueue: () => void;
  expireOffer: () => void;
  confirmPurchase: () => void;
  reset: () => void;
  isProductOccupied: (itemId: number) => boolean;
  forceStatus: (itemId: number, newStatus: QueueStatus, timeLeft?: number, isLimited?: boolean) => void;
}

const QueueContext = createContext<QueueContextType | undefined>(undefined);

export const QueueProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, setState] = useState<QueueState>({
    itemId: null,
    status: null,
    queuePosition: null,
    totalInQueue: null,
    timeLeft: null,
    isLimited: null,
  });
  const [occupiedItemId, setOccupiedItemId] = useState<number | null>(null);
  const [queue, setQueue] = useState<number[]>([]);

  const startCheckout = (itemId: number, isLimited: boolean = true) => {
    setOccupiedItemId(itemId);
    setState({
      itemId,
      status: 'CHECKOUT',
      queuePosition: null,
      totalInQueue: null,
      timeLeft: isLimited ? 120 : null,
      isLimited,
    });
  };

  const joinQueue = (itemId: number) => {
    const userId = 1;
    setQueue(prev => [...prev, userId]);
    const position = queue.length + 1;
    setOccupiedItemId(itemId);
    setState({
      itemId,
      status: 'QUEUED',
      queuePosition: position,
      totalInQueue: position + 2,
      timeLeft: 480,
      isLimited: true,
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
      isLimited: null,
    }));
    setOccupiedItemId(null);
  };

  const expireOffer = () => {
    setState(prev => ({
      ...prev,
      status: 'EXPIRED',
      timeLeft: null,
      isLimited: null,
    }));
    setOccupiedItemId(null);
  };

  const confirmPurchase = () => {
    setState(prev => ({
      ...prev,
      status: 'PURCHASED',
      timeLeft: null,
      isLimited: null,
    }));
    setOccupiedItemId(null);
  };

  const reset = () => {
    setState({
      itemId: null,
      status: null,
      queuePosition: null,
      totalInQueue: null,
      timeLeft: null,
      isLimited: null,
    });
    setOccupiedItemId(null);
    setQueue([]);
  };

  const isProductOccupied = (itemId: number) => occupiedItemId === itemId;

  const forceStatus = (
    itemId: number,
    newStatus: QueueStatus,
    timeLeft?: number,
    isLimited: boolean = true,
  ) => {
    setState({
      itemId,
      status: newStatus,
      queuePosition: newStatus === 'QUEUED' ? 2 : null,
      totalInQueue: newStatus === 'QUEUED' ? 4 : null,
      timeLeft: timeLeft || null,
      isLimited,
    });
    if (newStatus === 'CHECKOUT' || newStatus === 'QUEUED') {
      setOccupiedItemId(itemId);
    } else {
      setOccupiedItemId(null);
    }
  };

  return (
    <QueueContext.Provider value={{
      state,
      occupiedItemId,
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
