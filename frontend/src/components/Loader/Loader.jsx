import styles from './Loader.module.css';

export default function Loader({ fullPage = false }) {
  return (
    <div className={`${styles.wrapper} ${fullPage ? styles.fullPage : ''}`}>
      <div className={styles.ring}>
        <div /><div /><div /><div />
      </div>
    </div>
  );
}