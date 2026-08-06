import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueueProvider } from './context/QueueContext';
import AuthPage from './components/AuthPage/AuthPage';
import ProductCatalog from './components/ProductCatalog/ProductCatalog';
import ProductDetail from './components/ProductDetail/ProductDetail';
import QueuePage from './components/QueuePage/QueuePage';

function App() {
  return (
    <BrowserRouter>
      <QueueProvider>
        <Routes>
          <Route path="/" element={<AuthPage />} />
          <Route path="/products" element={<ProductCatalog />} />
          <Route path="/product/:id" element={<ProductDetail />} />
          <Route path="/product/:id/queue" element={<QueuePage />} />
        </Routes>
      </QueueProvider>
    </BrowserRouter>
  );
}

export default App;