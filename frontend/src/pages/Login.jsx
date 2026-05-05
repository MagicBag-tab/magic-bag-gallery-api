import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { login } from '../../api/api';
import { useAuth } from '../../context/AuthContext';
import styles from './Login.module.css';

export default function Login() {
  const [form, setForm] = useState({ correo_electronico: '', contrasena: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { loginUser } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e) => setForm(f => ({ ...f, [e.target.name]: e.target.value }));

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const data = await login(form.correo_electronico, form.contrasena);
      loginUser(data.token, data.role);
      navigate('/catalogo');
    } catch {
      setError('Credenciales inválidas. Intenta de nuevo.');
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

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label className={styles.label}>Correo electrónico</label>
            <input className={styles.input} type="email" name="correo_electronico" value={form.correo_electronico} onChange={handleChange} placeholder="tu@correo.com" required />
          </div>
          <div className={styles.field}>
            <label className={styles.label}>Contraseña</label>
            <input className={styles.input} type="password" name="contrasena" value={form.contrasena} onChange={handleChange} placeholder="••••••••" required />
          </div>
          {error && <p className={styles.error}>{error}</p>}
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