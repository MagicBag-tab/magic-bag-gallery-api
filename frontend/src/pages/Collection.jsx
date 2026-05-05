import { useState, useEffect } from 'react';
import { getColecciones } from '../../services/api';
import Loader from '../../components/Loader/Loader';
import Modal from '../../components/Modal/Modal';
import styles from './Collections.module.css';

export default function Collections() {
  const [colecciones, setColecciones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState(null);

  useEffect(() => {
    getColecciones().then(setColecciones).finally(() => setLoading(false));
  }, []);

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Agrupaciones curatoriales</p>
        <h1 className={styles.heroTitle}>Colecciones</h1>
        <p className={styles.heroSub}>Obras organizadas en torno a movimientos y narrativas artísticas</p>
      </div>

      {loading ? <Loader fullPage /> : (
        <div className={styles.grid}>
          {colecciones.map((c, i) => (
            <article
              key={c.id_coleccion}
              className={styles.card}
              style={{ animationDelay: `${i * 0.07}s` }}
              onClick={() => setSelected(c)}
            >
              <div className={styles.cardTop}>
                {c.exclusiva && <span className={styles.badge}>Exclusiva</span>}
              </div>
              <h3 className={styles.cardName}>{c.nombre}</h3>
              <p className={styles.cardDesc}>{c.descripcion}</p>
              <p className={styles.cardDate}>
                Lanzamiento: {new Date(c.fecha_lanzamiento).toLocaleDateString('es-GT', { year: 'numeric', month: 'long' })}
              </p>
            </article>
          ))}
        </div>
      )}

      {selected && (
        <Modal title={selected.nombre} onClose={() => setSelected(null)}>
          <div className={styles.detail}>
            {selected.exclusiva && <span className={styles.detailBadge}>Colección Exclusiva</span>}
            <p className={styles.detailDesc}>{selected.descripcion}</p>
            <p className={styles.detailMeta}>
              Fecha de lanzamiento: {new Date(selected.fecha_lanzamiento).toLocaleDateString('es-GT', { year: 'numeric', month: 'long', day: 'numeric' })}
            </p>
            {selected.pinturas?.length > 0 && (
              <div className={styles.detailSection}>
                <p className={styles.detailLabel}>Obras en esta colección</p>
                <ul className={styles.list}>
                  {selected.pinturas.map(p => <li key={p}>{p}</li>)}
                </ul>
              </div>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
}