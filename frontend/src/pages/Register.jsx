import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { registerCliente } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import styles from './Register.module.css';

export default function Register() {
  const [form, setForm] = useState({
    nombre: '', apellido: '', correo_electronico: '',
    telefono: '', contrasena: '', tipo_cliente: 'regular'
  });
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
      const data = await registerCliente(form);
      loginUser(data.token, data.role);
      navigate('/catalogo');
    } catch (err) {
      setError(err.message || 'Error al registrarse.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.ambient} />
      <div className={styles.card}>
        <div className={styles.header}>
          <p className={styles.eyebrow}>Nueva cuenta</p>
          <h1 className={styles.title}>Únete a la galería</h1>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.row}>
            <div className={styles.field}>
              <label className={styles.label}>Nombre</label>
              <input className={styles.input} name="nombre" value={form.nombre} onChange={handleChange} required />
            </div>
            <div className={styles.field}>
              <label className={styles.label}>Apellido</label>
              <input className={styles.input} name="apellido" value={form.apellido} onChange={handleChange} required />
            </div>
          </div>
          <div className={styles.field}>
            <label className={styles.label}>Correo electrónico</label>
            <input className={styles.input} type="email" name="correo_electronico" value={form.correo_electronico} onChange={handleChange} required />
          </div>
          <div className={styles.field}>
            <label className={styles.label}>Teléfono</label>
            <input className={styles.input} name="telefono" value={form.telefono} onChange={handleChange} required />
          </div>
          <div className={styles.field}>
            <label className={styles.label}>Contraseña</label>
            <input className={styles.input} type="password" name="contrasena" value={form.contrasena} onChange={handleChange} required />
          </div>

          {error && <p className={styles.error}>{error}</p>}

          <button className={styles.btn} type="submit" disabled={loading}>
            {loading ? 'Registrando...' : 'Crear cuenta'}
          </button>
        </form>

        <p className={styles.login}>
          ¿Ya tienes cuenta? <Link to="/login" className={styles.link}>Inicia sesión</Link>
        </p>
      </div>
    </div>
  );
}