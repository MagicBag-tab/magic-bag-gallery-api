import { useState, useEffect } from 'react';
import { getArtistas } from '../../services/api';
import Loader from '../../components/Loader/Loader';
import Modal from '../../components/Modal/Modal';
import styles from './Artists.module.css';

export default function Artists() {
  const [artistas, setArtistas] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState(null);

  useEffect(() => {
    getArtistas().then(setArtistas).finally(() => setLoading(false));
  }, []);

  const initials = (name) => name.split(' ').slice(0, 2).map(w => w[0]).join('');

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Nuestros creadores</p>
        <h1 className={styles.heroTitle}>Artistas</h1>
        <p className={styles.heroSub}>Los artistas más influyentes del arte contemporáneo global</p>
      </div>

      {loading ? <Loader fullPage /> : (
        <div className={styles.grid}>
          {artistas.map((a, i) => (
            <article
              key={a.id_artista}
              className={styles.card}
              style={{ animationDelay: `${i * 0.06}s` }}
              onClick={() => setSelected(a)}
            >
              <div className={styles.avatar}>
                <span>{initials(a.nombre_completo)}</span>
              </div>
              <div className={styles.info}>
                <h3 className={styles.name}>{a.nombre_completo}</h3>
                <p className={styles.nationality}>{a.nacionalidad}</p>
              </div>
              <span className={styles.arrow}>→</span>
            </article>
          ))}
        </div>
      )}

      {selected && (
        <Modal title="Artista" onClose={() => setSelected(null)}>
          <div className={styles.detail}>
            <div className={styles.detailAvatar}>
              <span>{initials(selected.nombre_completo)}</span>
            </div>
            <h2 className={styles.detailName}>{selected.nombre_completo}</h2>
            <p className={styles.detailNat}>{selected.nacionalidad}</p>

            {selected.pinturas?.length > 0 && (
              <div className={styles.detailSection}>
                <p className={styles.detailLabel}>Obras en galería</p>
                <ul className={styles.list}>
                  {selected.pinturas.map(p => <li key={p}>{p}</li>)}
                </ul>
              </div>
            )}
            {selected.colecciones?.length > 0 && (
              <div className={styles.detailSection}>
                <p className={styles.detailLabel}>Colecciones</p>
                <ul className={styles.list}>
                  {selected.colecciones.map(c => <li key={c}>{c}</li>)}
                </ul>
              </div>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
}