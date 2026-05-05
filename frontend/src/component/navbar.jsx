import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import styles from './Navbar.module.css';

export default function Navbar() {
  const { isAuthenticated, isEmpleado, logoutUser } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logoutUser();
    navigate('/login');
  };

  return (
    <header className={styles.navbar}>
      <div className={styles.logo}>
        Magic <span>Bag</span> Gallery
      </div>

      <nav className={styles.nav}>
        <NavLink to="/catalogo" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
          Catálogo
        </NavLink>
        <NavLink to="/artistas" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
          Artistas
        </NavLink>
        <NavLink to="/colecciones" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
          Colecciones
        </NavLink>
        <NavLink to="/tours" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
          Tours
        </NavLink>
        <NavLink to="/reportes" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
          Reportes
        </NavLink>
        {isEmpleado && (
          <NavLink to="/admin" className={({ isActive }) => `${styles.navLink}${isActive ? ' ' + styles.active : ''}`}>
            Admin
          </NavLink>
        )}
      </nav>

      <div className={styles.actions}>
        {isAuthenticated ? (
          <>
            {isEmpleado && <span className={styles.badge}>Empleado</span>}
            <button className={styles.btnOutline} onClick={handleLogout}>Salir</button>
          </>
        ) : (
          <>
            <button className={styles.btnOutline} onClick={() => navigate('/login')}>Iniciar sesión</button>
            <button className={styles.btnGold} onClick={() => navigate('/register')}>Registrarse</button>
          </>
        )}
      </div>
    </header>
  );
}