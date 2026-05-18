import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { catalogReducer, initialState } from '../hooks/useCatalogFilters';
import { AuthProvider, useAuth } from '../context/AuthContext';

describe('catalogReducer', () => {
  it('SET_PINTURAS carga las pinturas y deshabilita el loading', () => {
    const pinturas = [
      { id_pintura: 1, titulo: 'Skull', artista: 'Basquiat', exclusiva: true },
      { id_pintura: 2, titulo: 'Balloon Girl', artista: 'Banksy', exclusiva: false },
    ];
    const nextState = catalogReducer(initialState, {
      type: 'SET_PINTURAS',
      payload: pinturas,
    });
    expect(nextState.pinturas).toHaveLength(2);
    expect(nextState.loading).toBe(false);
    expect(nextState.error).toBeNull();
  });

  it('SET_SEARCH actualiza el campo search', () => {
    const nextState = catalogReducer(initialState, {
      type: 'SET_SEARCH',
      payload: 'Basquiat',
    });
    expect(nextState.search).toBe('Basquiat');
  });

  it('SET_FILTER_EXCLUSIVA actualiza el filtro', () => {
    const nextState = catalogReducer(initialState, {
      type: 'SET_FILTER_EXCLUSIVA',
      payload: 'si',
    });
    expect(nextState.filterExclusiva).toBe('si');
  });

  it('RESET_FILTERS vuelve search y filtro a valores iniciales', () => {
    const modifiedState = {
      ...initialState,
      search: 'algo',
      filterExclusiva: 'si',
    };
    const nextState = catalogReducer(modifiedState, { type: 'RESET_FILTERS' });
    expect(nextState.search).toBe('');
    expect(nextState.filterExclusiva).toBe('all');
  });

  it('SET_ERROR guarda el mensaje y desactiva loading', () => {
    const nextState = catalogReducer(initialState, {
      type: 'SET_ERROR',
      payload: 'Error de conexión',
    });
    expect(nextState.error).toBe('Error de conexión');
    expect(nextState.loading).toBe(false);
  });
});

describe('AuthContext', () => {
  it('empieza sin usuario autenticado cuando no hay token en localStorage', () => {
    localStorage.clear();
    const wrapper = ({ children }) => <AuthProvider>{children}</AuthProvider>;
    const { result } = renderHook(() => useAuth(), { wrapper });
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.isEmpleado).toBe(false);
  });

  it('loginUser guarda el token y actualiza el estado', () => {
    localStorage.clear();
    const wrapper = ({ children }) => <AuthProvider>{children}</AuthProvider>;
    const { result } = renderHook(() => useAuth(), { wrapper });

    act(() => {
      result.current.loginUser('token-abc', 'empleado');
    });

    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.isEmpleado).toBe(true);
    expect(localStorage.getItem('token')).toBe('token-abc');
  });

  it('logoutUser limpia el estado y el localStorage', () => {
    localStorage.setItem('token', 'token-xyz');
    localStorage.setItem('role', 'cliente');
    const wrapper = ({ children }) => <AuthProvider>{children}</AuthProvider>;
    const { result } = renderHook(() => useAuth(), { wrapper });

    act(() => {
      result.current.logoutUser();
    });

    expect(result.current.isAuthenticated).toBe(false);
    expect(localStorage.getItem('token')).toBeNull();
  });
});

function validateLoginForm({ correo, contrasena }) {
  const errors = {};
  if (!correo || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(correo)) {
    errors.correo = 'Correo electrónico inválido';
  }
  if (!contrasena || contrasena.length < 6) {
    errors.contrasena = 'La contraseña debe tener al menos 6 caracteres';
  }
  return errors;
}

describe('validateLoginForm', () => {
  it('retorna errores con campos vacíos', () => {
    const errors = validateLoginForm({ correo: '', contrasena: '' });
    expect(errors.correo).toBeDefined();
    expect(errors.contrasena).toBeDefined();
  });

  it('retorna error de correo con formato inválido', () => {
    const errors = validateLoginForm({ correo: 'no-es-email', contrasena: '123456' });
    expect(errors.correo).toBeDefined();
    expect(errors.contrasena).toBeUndefined();
  });

  it('retorna error de contraseña si tiene menos de 6 caracteres', () => {
    const errors = validateLoginForm({ correo: 'user@test.com', contrasena: '123' });
    expect(errors.contrasena).toBeDefined();
    expect(errors.correo).toBeUndefined();
  });

  it('no retorna errores con datos válidos', () => {
    const errors = validateLoginForm({
      correo: 'maria.perez@gmail.com',
      contrasena: 'secreto123',
    });
    expect(Object.keys(errors)).toHaveLength(0);
  });
});