import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import styles from './Navbar.module.css';

export default function Navbar() {
  const { isAuthenticated, isEmpleado, nombre, logoutUser } = useAuth();
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    logoutUser();
    setMenuOpen(false);
    navigate('/login');
  };

  const navClass = ({ isActive }) =>
    `${styles.navLink}${isActive ? ' ' + styles.active : ''}`;

  return (
    <header className={styles.navbar}>
      <div className={styles.logo}>Magic <span>Bag</span> Gallery</div>

      <nav className={styles.nav}>
        {!isEmpleado && (
          <>
            <NavLink to="/catalogo"    className={navClass}>Catálogo</NavLink>
            <NavLink to="/artistas"    className={navClass}>Artistas</NavLink>
            <NavLink to="/colecciones" className={navClass}>Colecciones</NavLink>
            <NavLink to="/tours"       className={navClass}>Tours</NavLink>
          </>
        )}

        {isEmpleado && (
          <NavLink to="/reportes" className={navClass}>Reportes</NavLink>
        )}

        {isEmpleado && (
          <NavLink to="/admin" className={navClass}>Admin</NavLink>
        )}

        {(isAuthenticated && !isEmpleado) && (
          <NavLink to="/mi-cuenta" className={navClass}>Mi cuenta</NavLink>
        )}
      </nav>

      <div className={styles.actions}>
        {isAuthenticated ? (
          <>
            {nombre && (
              <span className={styles.greeting}>
                <span className={styles.hola}>Hola,</span>{" "}
                <span className={styles.nombreUsuario}>{nombre}</span>
              </span>
            )}
            {isEmpleado && <span className={styles.badge}>Empleado</span>}
            <button className={styles.btnOutline} onClick={handleLogout}>Salir</button>
          </>
        ) : (
          <>
            <button className={styles.btnOutline} onClick={() => navigate('/login')}>Iniciar sesión</button>
            <button className={styles.btnGold}    onClick={() => navigate('/register')}>Registrarse</button>
          </>
        )}
      </div>

      <button
        className={`${styles.hamburger} ${menuOpen ? styles.open : ''}`}
        onClick={() => setMenuOpen(o => !o)}
        aria-label="Menú"
      >
        <span /><span /><span />
      </button>

      {menuOpen && (
        <div className={styles.mobileMenu}>
          {!isEmpleado && (
            <>
              <NavLink to="/catalogo"    className={navClass} onClick={() => setMenuOpen(false)}>Catálogo</NavLink>
              <NavLink to="/artistas"    className={navClass} onClick={() => setMenuOpen(false)}>Artistas</NavLink>
              <NavLink to="/colecciones" className={navClass} onClick={() => setMenuOpen(false)}>Colecciones</NavLink>
              <NavLink to="/tours"       className={navClass} onClick={() => setMenuOpen(false)}>Tours</NavLink>
            </>
          )}

          {isEmpleado && (
            <NavLink to="/reportes" className={navClass} onClick={() => setMenuOpen(false)}>Reportes</NavLink>
          )}
          {isEmpleado && (
            <NavLink to="/admin" className={navClass} onClick={() => setMenuOpen(false)}>Admin</NavLink>
          )}
          {isAuthenticated && !isEmpleado && (
            <NavLink to="/mi-cuenta" className={navClass} onClick={() => setMenuOpen(false)}>Mi cuenta</NavLink>
          )}

          <div className={styles.mobileDivider} />

          {isAuthenticated ? (
            <button className={styles.btnOutline} onClick={handleLogout}>Cerrar sesión</button>
          ) : (
            <>
              <button className={styles.btnOutline} onClick={() => { navigate('/login');    setMenuOpen(false); }}>Iniciar sesión</button>
              <button className={styles.btnGold}    onClick={() => { navigate('/register'); setMenuOpen(false); }}>Registrarse</button>
            </>
          )}
        </div>
      )}
    </header>
  );
}