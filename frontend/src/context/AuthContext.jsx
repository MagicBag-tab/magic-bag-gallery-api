import React, { createContext, useContext, useState, useCallback } from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const role = localStorage.getItem('role');
    const nombre = localStorage.getItem('nombre');
    const tipoEmpleado = localStorage.getItem('tipoEmpleado');
    return role ? { role, nombre, tipoEmpleado } : null;
  });

  const loginUser = useCallback((role, nombre = null, tipoEmpleado = null) => {
    localStorage.setItem('role', role);
    if (nombre) localStorage.setItem('nombre', nombre);
    else localStorage.removeItem('nombre');
    if (tipoEmpleado) localStorage.setItem('tipoEmpleado', tipoEmpleado);
    else localStorage.removeItem('tipoEmpleado');
    setUser({ role, nombre, tipoEmpleado });
  }, []);

  const logoutUser = useCallback(async () => {
    try {
      await fetch('/api/logout', { method: 'POST', credentials: 'include' });
    } catch {
      // La limpieza local debe ocurrir aunque el backend no responda.
    } finally {
      localStorage.clear();
      setUser(null);
    }
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
