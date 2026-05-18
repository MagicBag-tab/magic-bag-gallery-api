import React, { createContext, useContext, useState, useCallback } from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    const nombre = localStorage.getItem('nombre');
    const tipoEmpleado = localStorage.getItem('tipoEmpleado');
    return token ? { token, role, nombre, tipoEmpleado } : null;
  });

  const loginUser = useCallback((token, role, nombre, tipoEmpleado = null) => {
    localStorage.setItem('token', token);
    localStorage.setItem('role', role);
    localStorage.setItem('nombre', nombre);
    if (tipoEmpleado) localStorage.setItem('tipoEmpleado', tipoEmpleado);
    setUser({ token, role, nombre, tipoEmpleado });
  }, []);

  const logoutUser = useCallback(() => {
    localStorage.clear();
    setUser(null);
  }, []);

  const isEmpleado = user?.role === 'empleado';
  const isAuthenticated = !!user;
  const nombre = user?.nombre;
  const tipoEmpleado = user?.tipoEmpleado;

  return (
    <AuthContext.Provider value={{ user, loginUser, logoutUser, isEmpleado, isAuthenticated, nombre, tipoEmpleado }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);