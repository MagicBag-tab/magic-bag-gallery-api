import { createContext, useContext, useState, useCallback } from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    return token ? { token, role } : null;
  });

  const loginUser = useCallback((token, role) => {
    localStorage.setItem('token', token);
    localStorage.setItem('role', role);
    setUser({ token, role });
  }, []);

  const logoutUser = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    setUser(null);
  }, []);

  const isEmpleado = user?.role === 'empleado';
  const isAuthenticated = !!user;

  return (
    <AuthContext.Provider value={{ user, loginUser, logoutUser, isEmpleado, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);