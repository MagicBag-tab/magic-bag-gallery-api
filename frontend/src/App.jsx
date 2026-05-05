import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar/Navbar';
import Footer from './components/Footer/Footer';
import ProtectedRoute from './components/ProtectedRoute/ProtectedRoute';

import Login from './pages/Login/Login';
import Register from './pages/Register/Register';
import Catalog from './pages/Catalog/Catalog';
import Artists from './pages/Artists/Artists';
import Collections from './pages/Collection/Collection';
import Tours from './pages/Tours/Tours';
import Reports from './pages/Reports/Reports';
import Admin from './pages/Admin/Admin';
import NotFound from './pages/NotFound/NotFound';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route path="/" element={<Navigate to="/catalogo" replace />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/catalogo" element={<Catalog />} />
          <Route path="/artistas" element={<Artists />} />
          <Route path="/colecciones" element={<Collections />} />
          <Route path="/tours" element={<Tours />} />
          <Route path="/reportes" element={<Reports />} />
          <Route path="/admin" element={
            <ProtectedRoute requireEmpleado>
              <Admin />
            </ProtectedRoute>
          } />
          <Route path="*" element={<NotFound />} />
        </Routes>
        <Footer />
      </BrowserRouter>
    </AuthProvider>
  );
}