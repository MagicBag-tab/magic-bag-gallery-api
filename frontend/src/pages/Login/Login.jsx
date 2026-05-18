import React, { useState, useCallback } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { login } from '../../api/api';
import { useAuth } from '../../context/AuthContext';
import styles from './Login.module.css';

export function validateLoginForm({ correo_electronico, contrasena }) {
  const errors = {};
  if (!correo_electronico || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(correo_electronico)) {
    errors.correo_electronico = 'Ingresa un correo electrónico válido';
  }
  if (!contrasena || contrasena.length < 6) {
    errors.contrasena = 'La contraseña debe tener al menos 6 caracteres';
  }
  return errors;
}

export default function Login() {
  const [form, setForm] = useState({ correo_electronico: '', contrasena: '' });
  const [fieldErrors, setFieldErrors] = useState({});
  const [serverError, setServerError] = useState('');
  const [loading, setLoading] = useState(false);
  const { loginUser } = useAuth();
  const navigate = useNavigate();

  const handleChange = useCallback((e) => {
    const { name, value } = e.target;
    setForm((f) => ({ ...f, [name]: value }));
    if (fieldErrors[name]) {
      setFieldErrors((fe) => ({ ...fe, [name]: '' }));
    }
  }, [fieldErrors]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setServerError('');

    const errors = validateLoginForm(form);
    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }

    setLoading(true);
    try {
      const data = await login(form.correo_electronico, form.contrasena);
      loginUser(data.token, data.role);
      navigate('/catalogo');
    } catch {
      setServerError('Credenciales inválidas. Verifica tu correo y contraseña.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.ambient} />
      <div className={styles.card}>
        <div className={styles.header}>
          <p className={styles.eyebrow}>Bienvenido</p>
          <h1 className={styles.title}>Magic Bag Gallery</h1>
          <p className={styles.sub}>Accede a tu cuenta</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form} noValidate>
          <div className={styles.field}>
            <label className={styles.label} htmlFor="correo_electronico">
              Correo electrónico
            </label>
            <input
              id="correo_electronico"
              className={`${styles.input} ${fieldErrors.correo_electronico ? styles.inputError : ''}`}
              type="email"
              name="correo_electronico"
              value={form.correo_electronico}
              onChange={handleChange}
              placeholder="tu@correo.com"
              autoComplete="email"
              aria-describedby={fieldErrors.correo_electronico ? 'error-correo' : undefined}
            />
            {fieldErrors.correo_electronico && (
              <span id="error-correo" className={styles.fieldError} role="alert">
                {fieldErrors.correo_electronico}
              </span>
            )}
          </div>

          <div className={styles.field}>
            <label className={styles.label} htmlFor="contrasena">
              Contraseña
            </label>
            <input
              id="contrasena"
              className={`${styles.input} ${fieldErrors.contrasena ? styles.inputError : ''}`}
              type="password"
              name="contrasena"
              value={form.contrasena}
              onChange={handleChange}
              placeholder="••••••••"
              autoComplete="current-password"
              aria-describedby={fieldErrors.contrasena ? 'error-contra' : undefined}
            />
            {fieldErrors.contrasena && (
              <span id="error-contra" className={styles.fieldError} role="alert">
                {fieldErrors.contrasena}
              </span>
            )}
          </div>

          {serverError && (
            <p className={styles.error} role="alert">{serverError}</p>
          )}

          <button className={styles.btn} type="submit" disabled={loading}>
            {loading ? 'Iniciando...' : 'Iniciar sesión'}
          </button>
        </form>

        <p className={styles.register}>
          ¿No tienes cuenta? <Link to="/register" className={styles.link}>Regístrate</Link>
        </p>
      </div>
    </div>
  );
}