import { Link } from 'react-router-dom';
import styles from './NotFound.module.css';

export default function NotFound() {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>404</h1>
      <p className={styles.text}>Vaya, parece que te has perdido en nuestra galería.</p>
      <p className={styles.sub}>La página que buscas no existe o ha sido movida.</p>
      <Link to="/" className={styles.link}>Regresar al inicio</Link>
    </div>
  );
}