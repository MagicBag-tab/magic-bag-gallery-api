import { useState, useEffect, useMemo } from 'react';
import {
  getPinturas, deletePintura,
  getArtistas, deleteArtista, createArtista, updateArtista,
  getColecciones, deleteColeccion, createColeccion, updateColeccion,
  getTecnicas, deleteTecnica, createTecnica, updateTecnica,
  getVentas, deleteVenta,
  getUsuarios, deleteUsuario,
  getTours, deleteTour,
  getReservas,
} from '../../api/api';
import { useAuth } from '../../context/AuthContext';
import Modal from '../../components/Modal/Modal';
import Loader from '../../components/Loader/Loader';
import styles from './Admin.module.css';

const TABS_BY_TIPO = {
  guia:       ['Tours', 'Reservas'],
  asesor:     ['Ventas', 'Usuarios'],
  reclutador: ['Artistas', 'Colecciones'],
};
const ALL_TABS = ['Pinturas', 'Artistas', 'Colecciones', 'Técnicas', 'Ventas', 'Tours', 'Reservas', 'Usuarios'];

export default function Admin() {
  const { tipoEmpleado } = useAuth();

  const allowedTabs = useMemo(
    () => TABS_BY_TIPO[tipoEmpleado] ?? ALL_TABS,
    [tipoEmpleado]
  );

  const [tab,     setTab]     = useState(allowedTabs[0]);
  const [data,    setData]    = useState([]);
  const [loading, setLoading] = useState(true);
  const [modal,   setModal]   = useState(null);
  const [form,    setForm]    = useState({});
  const [msg,     setMsg]     = useState('');

  useEffect(() => {
    if (!allowedTabs.includes(tab)) setTab(allowedTabs[0]);
  }, [allowedTabs, tab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const fetchers = {
        Pinturas:    getPinturas,
        Artistas:    getArtistas,
        Colecciones: getColecciones,
        'Técnicas':  getTecnicas,
        Ventas:      getVentas,
        Tours:       getTours,
        Reservas:    getReservas,
        Usuarios:    getUsuarios,
      };
      const result = await fetchers[tab]();
      setData(result || []);
    } catch {
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
        Pinturas:    deletePintura,
        Artistas:    deleteArtista,
        Colecciones: deleteColeccion,
        'Técnicas':  deleteTecnica,
        Ventas:      deleteVenta,
        Tours:       deleteTour,
        Usuarios:    deleteUsuario,
      };
      if (deleters[tab]) await deleters[tab](id);
      fetchData();
    } catch (e) {
      setMsg('Error al eliminar: ' + e.message);
    }
  };

  const openCreate = () => { setForm({}); setModal({ type: 'create' }); setMsg(''); };
  const openEdit   = (item) => { setForm(item); setModal({ type: 'edit', item }); setMsg(''); };
  const handleFormChange = (e) => setForm(f => ({ ...f, [e.target.name]: e.target.value }));

  const handleSave = async (e) => {
    e.preventDefault();
    setMsg('');
    try {
      const payload = { ...form };

      const idKeys = {
        Pinturas:    'id_pintura',
        Artistas:    'id_artista',
        Colecciones: 'id_coleccion',
        'Técnicas':  'id_tecnica',
      };
      
      // Limpieza de campos para creación
      if (modal.type === 'create' && idKeys[tab]) {
        delete payload[idKeys[tab]];
      }

      if (tab === 'Artistas') {
        payload.id_reclutador = parseInt(payload.id_reclutador, 10);
        
        if (modal.type === 'create') await createArtista(payload);
        else await updateArtista(modal.item.id_artista, payload);
      } else if (tab === 'Colecciones') {
        const collPayload = { ...payload, exclusiva: payload.exclusiva === 'true' };
        if (modal.type === 'create') await createColeccion(collPayload);
        else await updateColeccion(modal.item.id_coleccion, collPayload);
      } else if (tab === 'Técnicas') {
        if (modal.type === 'create') await createTecnica(payload);
        else await updateTecnica(modal.item.id_tecnica, payload);
      }
      setModal(null);
      fetchData();
    } catch (err) {
      setMsg('Error: ' + err.message);
    }
  };

  const getIdField = (item) => {
    const map = {
      Pinturas:    'id_pintura',
      Artistas:    'id_artista',
      Colecciones: 'id_coleccion',
      'Técnicas':  'id_tecnica',
      Ventas:      'id_venta',
      Tours:       'id_tour',
      Reservas:    'id_cliente_tour',
      Usuarios:    'id_usuario',
    };
    return item[map[tab]];
  };

  const canCreate = ['Artistas', 'Colecciones', 'Técnicas'].includes(tab);
  const canEdit   = ['Artistas', 'Colecciones', 'Técnicas'].includes(tab);
  const canDelete = !['Reservas'].includes(tab);

  const renderRow = (item) => {
    switch (tab) {
      case 'Pinturas':    return <><td>{item.titulo}</td><td>{item.artista}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.exclusiva ? 'Sí' : 'No'}</td></>;
      case 'Artistas':    return <><td>{item.nombre_completo}</td><td>{item.nacionalidad}</td></>;
      case 'Colecciones': return <><td>{item.nombre}</td><td>{item.exclusiva ? 'Sí' : 'No'}</td><td>{item.fecha_lanzamiento}</td></>;
      case 'Técnicas':    return <><td>{item.nombre}</td><td className={styles.desc}>{item.descripcion}</td></>;
      case 'Ventas':      return <><td>{item.fecha_venta}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.id_cliente}</td></>;
      case 'Tours':       return <><td>{item.nombre}</td><td>{item.nombre_guia}</td><td>Q {Number(item.precio).toLocaleString()}</td><td>{item.fecha_inicio}</td></>;
      case 'Reservas':    return <><td>{item.nombre_tour}</td><td>{item.id_cliente}</td><td>{item.fecha_reserva}</td></>;
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
      case 'Ventas':      return ['Fecha', 'Total', 'Cliente'];
      case 'Tours':       return ['Nombre', 'Guía', 'Precio', 'Inicio'];
      case 'Reservas':    return ['Tour', 'Cliente', 'Fecha reserva'];
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

  const roleDesc = {
    guia:       'Guía de tours',
    asesor:     'Asesor de ventas',
    reclutador: 'Reclutador de artistas',
  }[tipoEmpleado] ?? 'Administración general';

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <p className={styles.eyebrow}>Panel de administración — {roleDesc}</p>
        <h1 className={styles.title}>Gestión de contenido</h1>
      </div>

      <div className={styles.tabs}>
        {allowedTabs.map(t => (
          <button key={t} className={`${styles.tab} ${tab === t ? styles.active : ''}`} onClick={() => setTab(t)}>{t}</button>
        ))}
      </div>

      <div className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>{tab}</h2>
          {canCreate && <button className={styles.btnCreate} onClick={openCreate}>+ Nuevo</button>}
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
                      {canEdit   && <button className={styles.btnEdit}   onClick={() => openEdit(item)}>Editar</button>}
                      {canDelete && <button className={styles.btnDelete} onClick={() => handleDelete(getIdField(item))}>Eliminar</button>}
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
          title={`${modal.type === 'create' ? 'Crear' : 'Editar'} ${tab.replace(/s$/, '')}`}
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