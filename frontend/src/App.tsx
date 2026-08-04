import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AuthPage from './components/AuthPage/AuthPage';
import ProductCatalog from './components/ProductCatalog/ProductCatalog.tsx';
import ProductDetail from './components/ProductDetail/ProductDetail.tsx';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<AuthPage />} />
        <Route path="/products" element={<ProductCatalog />} />
        <Route path="/product/:id" element={<ProductDetail />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;