import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../../context/AuthContext';
import Loader from '../../components/Loader/Loader';
import styles from './MiCuenta.module.css';

export default function MiCuenta() {
  const { nombre } = useAuth();

  const [reservas,  setReservas]  = useState([]);
  const [ventas,    setVentas]    = useState([]);
  const [loadingR,  setLoadingR]  = useState(true);
  const [loadingV,  setLoadingV]  = useState(true);
  const [errorR,    setErrorR]    = useState('');
  const [errorV,    setErrorV]    = useState('');

  const fetchReservas = useCallback(async () => {
    setLoadingR(true);
    setErrorR('');
    try {
      const res = await fetch('/api/me/reservas', {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      });
      if (!res.ok) throw new Error(`Error ${res.status}`);
      const data = await res.json();
      setReservas(data || []);
    } catch (e) {
      setErrorR('No se pudieron cargar tus reservas: ' + e.message);
    } finally {
      setLoadingR(false);
    }
  }, []);

  const fetchVentas = useCallback(async () => {
    setLoadingV(true);
    setErrorV('');
    try {
      const res = await fetch('/api/me/ventas', {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      });
      if (!res.ok) throw new Error(`Error ${res.status}`);
      const data = await res.json();
      setVentas(data || []);
    } catch (e) {
      setErrorV('No se pudo cargar tu historial: ' + e.message);
    } finally {
      setLoadingV(false);
    }
  }, []);

  useEffect(() => {
    fetchReservas();
    fetchVentas();
  }, [fetchReservas, fetchVentas]);

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Tu espacio personal</p>
        <h1 className={styles.heroTitle}>
          {nombre ? `Hola, ${nombre}` : 'Mi cuenta'}
        </h1>
        <p className={styles.heroSub}>Tus reservas de tours y tu historial de compras</p>
      </div>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Mis reservas de tours</h2>

        {loadingR && <Loader />}

        {!loadingR && errorR && (
          <div className={styles.errorBox}>
            <span>{errorR}</span>
            <button className={styles.retryBtn} onClick={fetchReservas}>Reintentar</button>
          </div>
        )}

        {!loadingR && !errorR && reservas.length === 0 && (
          <p className={styles.empty}>Aún no tienes reservas de tours.</p>
        )}

        {!loadingR && !errorR && reservas.length > 0 && (
          <div className={styles.tableWrapper}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Tour</th>
                  <th>Fecha de reserva</th>
                  <th>Horario</th>
                  <th>Precio</th>
                </tr>
              </thead>
              <tbody>
                {reservas.map(r => (
                  <tr key={r.id_cliente_tour}>
                    <td>{r.nombre_tour}</td>
                    <td>{r.fecha_reserva}</td>
                    <td>{r.horario ?? '—'}</td>
                    <td>{r.precio ? `Q ${Number(r.precio).toLocaleString()}` : '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Mi historial de compras</h2>

        {loadingV && <Loader />}

        {!loadingV && errorV && (
          <div className={styles.errorBox}>
            <span>{errorV}</span>
            <button className={styles.retryBtn} onClick={fetchVentas}>Reintentar</button>
          </div>
        )}

        {!loadingV && !errorV && ventas.length === 0 && (
          <p className={styles.empty}>Todavía no has realizado compras en la galería.</p>
        )}

        {!loadingV && !errorV && ventas.length > 0 && (
          <div className={styles.tableWrapper}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Fecha</th>
                  <th>Pintura(s)</th>
                  <th>Total</th>
                  <th>Estado envío</th>
                </tr>
              </thead>
              <tbody>
                {ventas.map(v => (
                  <tr key={v.id_venta}>
                    <td>{v.fecha_venta}</td>
                    <td>{v.pinturas ?? '—'}</td>
                    <td className={styles.gold}>Q {Number(v.total).toLocaleString()}</td>
                    <td>
                      <span className={`${styles.badge} ${styles[v.estado_envio?.replace(' ', '_')] ?? ''}`}>
                        {v.estado_envio ?? 'Sin envío'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}