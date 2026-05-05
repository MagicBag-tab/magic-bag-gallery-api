import { useState, useEffect } from 'react';
import { getTours, createReserva } from '../../api/api';
import { useAuth } from '../../context/AuthContext';
import Loader from '../../components/Loader/Loader';
import Modal from '../../components/Modal/Modal';
import styles from './Tours.module.css';

export default function Tours() {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState(null);
  const [reservaForm, setReservaForm] = useState({ id_cliente: '', fecha_reserva: '' });
  const [reservaMsg, setReservaMsg] = useState('');
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    getTours().then(setTours).finally(() => setLoading(false));
  }, []);

  const handleReserva = async (e) => {
    e.preventDefault();
    setReservaMsg('');
    try {
      await createReserva({ ...reservaForm, id_tour: selected.id_tour, id_cliente: Number(reservaForm.id_cliente) });
      setReservaMsg('¡Reserva creada exitosamente!');
    } catch (err) {
      setReservaMsg('Error: ' + err.message);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Experiencias guiadas</p>
        <h1 className={styles.heroTitle}>Tours & Reservas</h1>
        <p className={styles.heroSub}>Recorre la galería con nuestros expertos en arte contemporáneo</p>
      </div>

      {loading ? <Loader fullPage /> : (
        <div className={styles.grid}>
          {tours.map((t, i) => (
            <article key={t.id_tour} className={styles.card} style={{ animationDelay: `${i * 0.05}s` }}>
              <div className={styles.cardHeader}>
                <span className={styles.price}>Q {Number(t.precio).toLocaleString()}</span>
                <span className={styles.horario}>{t.horario}</span>
              </div>
              <h3 className={styles.cardName}>{t.nombre}</h3>
              <p className={styles.cardDesc}>{t.descripcion}</p>
              <div className={styles.cardMeta}>
                <span>🗓 {new Date(t.fecha_inicio).toLocaleDateString('es-GT')}</span>
                <span>👤 {t.nombre_guia}</span>
              </div>
              {isAuthenticated && (
                <button className={styles.btnReserva} onClick={() => { setSelected(t); setReservaMsg(''); setReservaForm({ id_cliente: '', fecha_reserva: '' }); }}>
                  Reservar
                </button>
              )}
            </article>
          ))}
        </div>
      )}

      {selected && (
        <Modal title={`Reservar: ${selected.nombre}`} onClose={() => setSelected(null)}>
          <form onSubmit={handleReserva} className={styles.reservaForm}>
            <div className={styles.field}>
              <label className={styles.label}>ID de Cliente</label>
              <input className={styles.input} type="number" value={reservaForm.id_cliente} onChange={e => setReservaForm(f => ({ ...f, id_cliente: e.target.value }))} placeholder="Ej: 1" required />
            </div>
            <div className={styles.field}>
              <label className={styles.label}>Fecha de reserva</label>
              <input className={styles.input} type="date" value={reservaForm.fecha_reserva} onChange={e => setReservaForm(f => ({ ...f, fecha_reserva: e.target.value }))} required />
            </div>
            {reservaMsg && (
              <p className={reservaMsg.startsWith('Error') ? styles.error : styles.success}>{reservaMsg}</p>
            )}
            <button className={styles.btnSubmit} type="submit">Confirmar reserva</button>
          </form>
        </Modal>
      )}

      {!isAuthenticated && (
        <p className={styles.loginNote}>
          <a href="/login" className={styles.loginLink}>Inicia sesión</a> para hacer una reserva.
        </p>
      )}
    </div>
  );
}