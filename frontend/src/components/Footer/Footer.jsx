import styles from './Footer.module.css';

export default function Footer() {
  return (
    <footer className={styles.footer}>
      <div className={styles.divider} />
      <div className={styles.inner}>
        <p className={styles.brand}>Magic Bag Gallery</p>
        <p className={styles.copy}>© {new Date().getFullYear()} — Arte que trasciende</p>
      </div>
    </footer>
  );
}