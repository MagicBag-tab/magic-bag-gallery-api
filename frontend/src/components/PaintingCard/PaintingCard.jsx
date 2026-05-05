import styles from './PaintingCard.module.css';

export default function PaintingCard({ pintura, onClick }) {
  return (
    <article className={styles.card} onClick={() => onClick?.(pintura)}>
      <div className={styles.imageWrapper}>
        <div className={styles.imagePlaceholder}>
          <span className={styles.initial}>{pintura.titulo?.[0]}</span>
        </div>
        {pintura.exclusiva && <span className={styles.badge}>Exclusiva</span>}
      </div>
      <div className={styles.info}>
        <p className={styles.artist}>{pintura.artista}</p>
        <h3 className={styles.title}>{pintura.titulo}</h3>
        <div className={styles.footer}>
          <span className={styles.price}>Q {Number(pintura.precio).toLocaleString()}</span>
          {pintura.coleccion && <span className={styles.collection}>{pintura.coleccion}</span>}
        </div>
        {pintura.tecnicas?.length > 0 && (
          <div className={styles.tecnicas}>
            {pintura.tecnicas.slice(0, 2).map(t => (
              <span key={t} className={styles.tag}>{t}</span>
            ))}
          </div>
        )}
      </div>
    </article>
  );
}