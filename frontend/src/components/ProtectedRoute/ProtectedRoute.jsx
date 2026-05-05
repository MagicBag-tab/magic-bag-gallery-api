import { Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

export default function ProtectedRoute({ children, requireEmpleado = false }) {
  const { isAuthenticated, isEmpleado } = useAuth();

  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (requireEmpleado && !isEmpleado) return <Navigate to="/catalogo" replace />;

  return children;
}