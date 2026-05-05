import { useState, useEffect } from 'react';
import {
  getPinturas, deletePintura,
  getArtistas, deleteArtista, createArtista, updateArtista,
  getColecciones, deleteColeccion, createColeccion, updateColeccion,
  getTecnicas, deleteTecnica, createTecnica, updateTecnica,
  getVentas, deleteVenta,
  getUsuarios, deleteUsuario,
  getTours, deleteTour
} from '../../services/api';
import Modal from '../../components/Modal/Modal';
import Loader from '../../components/Loader/Loader';
import styles from './Admin.module.css';

const TABS = ['Pinturas', 'Artistas', 'Colecciones', 'Técnicas', 'Ventas', 'Tours', 'Usuarios'];

export default function Admin() {
  const [tab, setTab] = useState('Pinturas');
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modal, setModal] = useState(null); // { type: 'create'|'edit', item? }
  const [form, setForm] = useState({});
  const [msg, setMsg] = useState('');

  const fetchData = async () => {
    setLoading(true);
    try {
      const fetchers = {
        Pinturas: getPinturas, Artistas: getArtistas, Colecciones: getColecciones,
        Técnicas: getTecnicas, Ventas: getVentas, Tours: getTours, Usuarios: getUsuarios
      };
      const result = await fetchers[tab]();
      setData(result || []);
    } catch (e) {
      setData([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); setMsg(''); }, [tab]);

  const handleDelete = async (id) => {
    if (!window.confirm('¿Eliminar este registro?')) return;
    try {
      const deleters = {
        Pinturas: deletePintura, Artistas: deleteArtista, Colecciones: deleteColeccion,
        Técnicas: deleteTecnica, Ventas: deleteVenta, Tours: deleteTour, Usuarios: deleteUsuario
      };
      await deleters[tab](id);
      fetchData();
    } catch (e) {
      setMsg('Error al eliminar: ' + e.message);
    }
  };

  const openCreate = () => { setForm({}); setModal({ type: 'create' }); setMsg(''); };
  const openEdit = (item) => { setForm(item); setModal({ type: 'edit', item }); setMsg(''); };

  const handleFormChange = (e) => setForm(f => ({ ...f, [e.target.name]: e.target.value }));

  const handleSave = async (e) => {
    e.preventDefault();
    setMsg('');
    try {
      if (tab === 'Artistas') {
        if (modal.type === 'create') await createArtista(form);
        else await updateArtista(modal.item.id_artista, form);
      } else if (tab === 'Colecciones') {
        if (modal.type === 'create') await createColeccion({ ...form, exclusiva: form.exclusiva === 'true' });
        else await updateColeccion(modal.item.id_coleccion, { ...form, exclusiva: form.exclusiva === 'true' });
      } else if (tab === 'Técnicas') {
        if (modal.type === 'create') await createTecnica(form);
        else await updateTecnica(modal.item.id_tecnica, form);
      }
      setModal(null);
      fetchData();
    } catch (err) {
      setMsg('Error: ' + err.message);
    }
  };

  const getIdField = (item) => {
    const map = { Pinturas: 'id_pintura', Artistas: 'id_artista', Colecciones: 'id_coleccion', Técnicas: 'id_tecnica', Ventas: 'id_venta', Tours: 'id_tour', Usuarios: 'id_usuario' };
    return item[map[tab]];
  };

  const canCreate = ['Artistas', 'Colecciones', 'Técnicas'].includes(tab);
  const canEdit   = ['Artistas', 'Colecciones', 'Técnicas'].includes(tab);

  const renderRow = (item) => {
    switch (tab) {
      case 'Pinturas':    return <><td>{item.titulo}</td><td>{item.artista}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.exclusiva ? 'Sí' : 'No'}</td></>;
      case 'Artistas':    return <><td>{item.nombre_completo}</td><td>{item.nacionalidad}</td></>;
      case 'Colecciones': return <><td>{item.nombre}</td><td>{item.exclusiva ? 'Sí' : 'No'}</td><td>{item.fecha_lanzamiento}</td></>;
      case 'Técnicas':    return <><td>{item.nombre}</td><td className={styles.desc}>{item.descripcion}</td></>;
      case 'Ventas':      return <><td>{item.id_venta}</td><td>{item.fecha_venta}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.id_cliente}</td></>;
      case 'Tours':       return <><td>{item.nombre}</td><td>{item.nombre_guia}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.fecha_inicio}</td></>;
      case 'Usuarios':    return <><td>{item.nombre} {item.apellido}</td><td>{item.correo_electronico}</td><td>{item.telefono}</td></>;
      default: return null;
    }
  };

  const renderHeaders = () => {
    switch (tab) {
      case 'Pinturas':    return ['Título', 'Artista', 'Precio', 'Exclusiva'];
      case 'Artistas':    return ['Nombre', 'Nacionalidad'];
      case 'Colecciones': return ['Nombre', 'Exclusiva', 'Lanzamiento'];
      case 'Técnicas':    return ['Nombre', 'Descripción'];
      case 'Ventas':      return ['ID', 'Fecha', 'Total', 'Cliente'];
      case 'Tours':       return ['Nombre', 'Guía', 'Precio', 'Inicio'];
      case 'Usuarios':    return ['Nombre', 'Correo', 'Teléfono'];
      default: return [];
    }
  };

  const renderForm = () => {
    switch (tab) {
      case 'Artistas':
        return (
          <>
            <div className={styles.field}><label className={styles.label}>Nombre completo</label><input className={styles.input} name="nombre_completo" value={form.nombre_completo || ''} onChange={handleFormChange} required /></div>
            <div className={styles.field}><label className={styles.label}>Nacionalidad</label><input className={styles.input} name="nacionalidad" value={form.nacionalidad || ''} onChange={handleFormChange} required /></div>
            <div className={styles.field}><label className={styles.label}>ID Reclutador</label><input className={styles.input} type="number" name="id_reclutador" value={form.id_reclutador || ''} onChange={handleFormChange} required /></div>
          </>
        );
      case 'Colecciones':
        return (
          <>
            <div className={styles.field}><label className={styles.label}>Nombre</label><input className={styles.input} name="nombre" value={form.nombre || ''} onChange={handleFormChange} required /></div>
            <div className={styles.field}><label className={styles.label}>Descripción</label><textarea className={styles.input} name="descripcion" value={form.descripcion || ''} onChange={handleFormChange} rows={3} required /></div>
            <div className={styles.field}>
              <label className={styles.label}>Exclusiva</label>
              <select className={styles.input} name="exclusiva" value={String(form.exclusiva) || 'false'} onChange={handleFormChange}>
                <option value="false">No</option>
                <option value="true">Sí</option>
              </select>
            </div>
            <div className={styles.field}><label className={styles.label}>Fecha lanzamiento</label><input className={styles.input} type="date" name="fecha_lanzamiento" value={form.fecha_lanzamiento || ''} onChange={handleFormChange} required /></div>
          </>
        );
      case 'Técnicas':
        return (
          <>
            <div className={styles.field}><label className={styles.label}>Nombre</label><input className={styles.input} name="nombre" value={form.nombre || ''} onChange={handleFormChange} required /></div>
            <div className={styles.field}><label className={styles.label}>Descripción</label><textarea className={styles.input} name="descripcion" value={form.descripcion || ''} onChange={handleFormChange} rows={3} required /></div>
          </>
        );
      default: return <p className={styles.noForm}>Este módulo no soporta edición desde aquí.</p>;
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Panel de administración</p>
          <h1 className={styles.title}>Gestión de contenido</h1>
        </div>
      </div>

      <div className={styles.tabs}>
        {TABS.map(t => (
          <button
            key={t}
            className={`${styles.tab} ${tab === t ? styles.active : ''}`}
            onClick={() => setTab(t)}
          >
            {t}
          </button>
        ))}
      </div>

      <div className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>{tab}</h2>
          {canCreate && (
            <button className={styles.btnCreate} onClick={openCreate}>+ Nuevo</button>
          )}
        </div>

        {msg && <p className={styles.error}>{msg}</p>}

        {loading ? <Loader /> : (
          <div className={styles.tableWrapper}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  {renderHeaders().map(h => <th key={h}>{h}</th>)}
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {data.map(item => (
                  <tr key={getIdField(item)}>
                    <td className={styles.idCell}>{getIdField(item)}</td>
                    {renderRow(item)}
                    <td className={styles.actions}>
                      {canEdit && (
                        <button className={styles.btnEdit} onClick={() => openEdit(item)}>Editar</button>
                      )}
                      <button className={styles.btnDelete} onClick={() => handleDelete(getIdField(item))}>Eliminar</button>
                    </td>
                  </tr>
                ))}
                {data.length === 0 && (
                  <tr><td colSpan={20} className={styles.empty}>Sin registros</td></tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {modal && (
        <Modal
          title={`${modal.type === 'create' ? 'Crear' : 'Editar'} ${tab.slice(0, -1)}`}
          onClose={() => setModal(null)}
        >
          <form onSubmit={handleSave} className={styles.form}>
            {renderForm()}
            {msg && <p className={styles.error}>{msg}</p>}
            <div className={styles.formActions}>
              <button type="button" className={styles.btnCancel} onClick={() => setModal(null)}>Cancelar</button>
              <button type="submit" className={styles.btnSave}>Guardar</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  );
}