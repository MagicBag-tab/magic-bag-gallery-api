import React, { useState, useCallback } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import styles from './Navbar.module.css';

export default function Navbar() {
  const { isAuthenticated, isEmpleado, logoutUser } = useAuth();
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = useCallback(() => {
    logoutUser();
    navigate('/login');
    setMenuOpen(false);
  }, [logoutUser, navigate]);

  const closeMenu = useCallback(() => setMenuOpen(false), []);

  const linkClass = ({ isActive }) =>
    `${styles.navLink}${isActive ? ' ' + styles.active : ''}`;

  const links = (
    <>
      <NavLink to="/catalogo"    className={linkClass} onClick={closeMenu}>Catálogo</NavLink>
      <NavLink to="/artistas"    className={linkClass} onClick={closeMenu}>Artistas</NavLink>
      <NavLink to="/colecciones" className={linkClass} onClick={closeMenu}>Colecciones</NavLink>
      <NavLink to="/tours"       className={linkClass} onClick={closeMenu}>Tours</NavLink>
      <NavLink to="/reportes"    className={linkClass} onClick={closeMenu}>Reportes</NavLink>
      {isEmpleado && (
        <NavLink to="/admin" className={linkClass} onClick={closeMenu}>Admin</NavLink>
      )}
    </>
  );

  return (
    <header className={styles.navbar}>
      <div className={styles.logo}>
        Magic <span>Bag</span> Gallery
      </div>

      <nav className={styles.nav} aria-label="Navegación principal">
        {links}
      </nav>

      <div className={styles.actions}>
        {isAuthenticated ? (
          <>
            {isEmpleado && <span className={styles.badge}>Empleado</span>}
            <button className={styles.btnOutline} onClick={handleLogout}>Salir</button>
          </>
        ) : (
          <>
            <button className={styles.btnOutline} onClick={() => { navigate('/login'); closeMenu(); }}>
              Iniciar sesión
            </button>
            <button className={styles.btnGold} onClick={() => { navigate('/register'); closeMenu(); }}>
              Registrarse
            </button>
          </>
        )}
      </div>

      <button
        className={styles.hamburger}
        onClick={() => setMenuOpen((o) => !o)}
        aria-label={menuOpen ? 'Cerrar menú' : 'Abrir menú'}
        aria-expanded={menuOpen}
      >
        <span className={`${styles.bar} ${menuOpen ? styles.barTop : ''}`} />
        <span className={`${styles.bar} ${menuOpen ? styles.barMid : ''}`} />
        <span className={`${styles.bar} ${menuOpen ? styles.barBot : ''}`} />
      </button>

      {menuOpen && (
        <div className={styles.mobileMenu} role="dialog" aria-label="Menú móvil">
          <nav className={styles.mobileNav}>
            {links}
          </nav>
          <div className={styles.mobileActions}>
            {isAuthenticated ? (
              <button className={styles.btnGold} onClick={handleLogout}>Cerrar sesión</button>
            ) : (
              <>
                <button className={styles.btnOutline} onClick={() => { navigate('/login'); closeMenu(); }}>
                  Iniciar sesión
                </button>
                <button className={styles.btnGold} onClick={() => { navigate('/register'); closeMenu(); }}>
                  Registrarse
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </header>
  );
}