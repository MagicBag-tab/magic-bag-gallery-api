import { useState, useEffect } from 'react';
import { getPinturas } from '../../api/api';
import PaintingCard from '../../components/PaintingCard/PaintingCard';
import Loader from '../../components/Loader/Loader';
import Modal from '../../components/Modal/Modal';
import styles from './Catalog.module.css';

export default function Catalog() {
  const [pinturas, setPinturas] = useState([]);
  const [filtered, setFiltered] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [filterExclusiva, setFilterExclusiva] = useState('all');
  const [selected, setSelected] = useState(null);

  useEffect(() => {
    getPinturas()
      .then(data => { setPinturas(data); setFiltered(data); })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    let result = pinturas;
    if (search) {
      const q = search.toLowerCase();
      result = result.filter(p =>
        p.titulo.toLowerCase().includes(q) ||
        p.artista.toLowerCase().includes(q) ||
        p.coleccion?.toLowerCase().includes(q)
      );
    }
    if (filterExclusiva === 'si') result = result.filter(p => p.exclusiva);
    if (filterExclusiva === 'no') result = result.filter(p => !p.exclusiva);
    setFiltered(result);
  }, [search, filterExclusiva, pinturas]);

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Colección permanente</p>
        <h1 className={styles.heroTitle}>Catálogo de Obras</h1>
        <p className={styles.heroSub}>Arte contemporáneo de los artistas más influyentes del mundo</p>
      </div>

      <div className={styles.controls}>
        <input
          className={styles.search}
          type="text"
          placeholder="Buscar por título, artista o colección..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <div className={styles.filters}>
          <span className={styles.filterLabel}>Exclusividad:</span>
          {['all', 'si', 'no'].map(v => (
            <button key={v} className={`${styles.filterBtn} ${filterExclusiva === v ? styles.active : ''}`} onClick={() => setFilterExclusiva(v)}>
              {v === 'all' ? 'Todas' : v === 'si' ? 'Exclusivas' : 'Estándar'}
            </button>
          ))}
        </div>
        <span className={styles.count}>{filtered.length} obras</span>
      </div>

      {loading ? <Loader fullPage /> : (
        <div className={styles.grid}>
          {filtered.map((p, i) => (
            <div key={p.id_pintura} style={{ animationDelay: `${i * 0.04}s` }}>
              <PaintingCard pintura={p} onClick={setSelected} />
            </div>
          ))}
          {filtered.length === 0 && (
            <div className={styles.empty}><p>No se encontraron obras con esos criterios.</p></div>
          )}
        </div>
      )}

      {selected && (
        <Modal title="Detalle de obra" onClose={() => setSelected(null)}>
          <div className={styles.detail}>
            <div className={styles.detailImageWrapper}>
              <div className={styles.detailImage}>
                <span className={styles.detailInitial}>{selected.titulo?.[0]}</span>
              </div>
              {selected.exclusiva && <span className={styles.detailBadge}>Obra Exclusiva</span>}
            </div>
            <p className={styles.detailArtist}>{selected.artista}</p>
            <h2 className={styles.detailTitle}>{selected.titulo}</h2>
            <p className={styles.detailPrice}>Q {Number(selected.precio).toLocaleString()}</p>
            {selected.coleccion && <p className={styles.detailCollection}>Colección: {selected.coleccion}</p>}
            <p className={styles.detailDesc}>{selected.descripcion}</p>
            {selected.tecnicas?.length > 0 && (
              <div className={styles.detailTecnicas}>
                <p className={styles.detailLabel}>Técnicas:</p>
                <div className={styles.tags}>
                  {selected.tecnicas.map(t => <span key={t} className={styles.tag}>{t}</span>)}
                </div>
              </div>
            )}
            {selected.fecha_creacion && (
              <p className={styles.detailMeta}>Fecha de creación: {new Date(selected.fecha_creacion).getFullYear()}</p>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
}