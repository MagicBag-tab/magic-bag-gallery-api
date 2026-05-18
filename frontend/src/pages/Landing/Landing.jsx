import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getPinturas } from '../../api/api';
import styles from './Landing.module.css';

const coleccionesDestacadas = [
  { nombre: 'Neo-Expresionismo', desc: 'Obras clave del movimiento neo-expresionista', exclusiva: true },
  { nombre: 'Arte Callejero', desc: 'Street art de los artistas más influyentes del mundo', exclusiva: false },
  { nombre: 'Identidad y Cultura', desc: 'Obras que exploran la identidad racial y cultural', exclusiva: false },
];

export default function Landing() {
  const navigate = useNavigate();
  const [pinturas, setPinturas] = useState([]);

  useEffect(() => {
    getPinturas()
      .then((data) => setPinturas((data ?? []).slice(0, 8)))
      .catch(() => setPinturas([]));
  }, []);

  const carouselItems = useMemo(() => {
    if (pinturas.length === 0) return [];
    return [...pinturas, ...pinturas];
  }, [pinturas]);

  return (
    <div className={styles.landingPage}>
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <p className={styles.eyebrow}>Galería de arte contemporáneo</p>
          <h1 className={styles.heroTitle}>Magic Bag Gallery</h1>
          <p className={styles.heroSub}>Descubre los artistas más influyentes del arte contemporáneo global</p>
          <div className={styles.heroActions}>
            <button className={styles.btnPrimary} onClick={() => navigate('/catalogo')}>Explorar catálogo</button>
            <button className={styles.btnOutline} onClick={() => navigate('/login')}>Iniciar sesión</button>
          </div>
        </div>
      </section>

      <section className={styles.stats} aria-label="Estadísticas de Magic Bag Gallery">
        <div className={styles.statCard}>
          <span className={styles.statNumber}>25</span>
          <p className={styles.statLabel}>Artistas internacionales</p>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statNumber}>25</span>
          <p className={styles.statLabel}>Colecciones exclusivas</p>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statNumber}>25+</span>
          <p className={styles.statLabel}>Obras disponibles</p>
        </div>
      </section>

      <section className={styles.carouselSection}>
        <div className={styles.carouselHeader}>
          <p className={styles.eyebrow}>Obras en exhibición</p>
          <h2 className={styles.sectionTitle}>Piezas destacadas</h2>
        </div>

        {carouselItems.length > 0 ? (
          <div className={styles.carouselShell} aria-label="Carrusel de obras destacadas">
            <div className={styles.carouselTrack}>
              {carouselItems.map((pintura, index) => (
                <article
                  key={`${pintura.id_pintura}-${index}`}
                  className={styles.paintingSlide}
                  onClick={() => navigate('/catalogo')}
                >
                  <div className={styles.paintingImage}>
                    {pintura.imagen_path ? (
                      <img src={`http://localhost:8888${pintura.imagen_path}`} alt={pintura.titulo} />
                    ) : (
                      <span className={styles.paintingInitial}>{pintura.titulo?.[0]}</span>
                    )}
                  </div>
                  <div className={styles.paintingMeta}>
                    <p className={styles.paintingArtist}>{pintura.artista}</p>
                    <h3 className={styles.paintingTitle}>{pintura.titulo}</h3>
                    <p className={styles.paintingPrice}>Q {Number(pintura.precio).toLocaleString()}</p>
                  </div>
                </article>
              ))}
            </div>
          </div>
        ) : (
          <p className={styles.carouselEmpty}>Explora el catálogo completo para ver las obras disponibles.</p>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Colecciones destacadas</h2>
        <div className={styles.collectionsGrid}>
          {coleccionesDestacadas.map((coleccion) => (
            <article key={coleccion.nombre} className={styles.collectionCard}>
              {coleccion.exclusiva && <span className={styles.badge}>Exclusiva</span>}
              <h3 className={styles.collectionName}>{coleccion.nombre}</h3>
              <p className={styles.collectionDesc}>{coleccion.desc}</p>
              <button className={styles.btnLink} onClick={() => navigate('/colecciones')}>Ver colección</button>
            </article>
          ))}
        </div>
      </section>

      <section className={styles.cta}>
        <h2 className={styles.ctaTitle}>¿Listo para descubrir el arte?</h2>
        <p className={styles.ctaSub}>Únete a nuestra comunidad de coleccionistas y amantes del arte</p>
        <button className={styles.btnWhite} onClick={() => navigate('/register')}>Crear cuenta gratis</button>
      </section>
    </div>
  );
}
