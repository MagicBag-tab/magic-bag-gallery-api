const BASE = '/api';

function getHeaders(auth = false) {
  const h = { 'Content-Type': 'application/json' };
  if (auth) {
    const token = localStorage.getItem('token');
    if (token) h['Authorization'] = `Bearer ${token}`;
  }
  return h;
}

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, options);
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Error ${res.status}`);
  }
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('json')) return res.json();
  return res;
}

export const login = (correo_electronico, contrasena) =>
  request('/login', { method: 'POST', headers: getHeaders(), body: JSON.stringify({ correo_electronico, contrasena }) });

export const registerCliente = (data) =>
  request('/register/cliente', { method: 'POST', headers: getHeaders(), body: JSON.stringify(data) });

export const getPinturas  = ()         => request('/pinturas',      { headers: getHeaders() });
export const getPintura   = (id)       => request(`/pinturas/${id}`,{ headers: getHeaders() });
export const createPintura= (data)     => request('/pinturas',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const updatePintura= (id, data) => request(`/pinturas/${id}`,{ method: 'PUT',    headers: getHeaders(true), body: JSON.stringify(data) });
export const deletePintura= (id)       => request(`/pinturas/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getArtistas  = ()         => request('/artistas',      { headers: getHeaders() });
export const getArtista   = (id)       => request(`/artistas/${id}`,{ headers: getHeaders() });
export const createArtista= (data)     => request('/artistas',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const updateArtista= (id, data) => request(`/artistas/${id}`,{ method: 'PUT',    headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteArtista= (id)       => request(`/artistas/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getColecciones  = ()         => request('/colecciones',      { headers: getHeaders() });
export const getColeccion    = (id)       => request(`/colecciones/${id}`,{ headers: getHeaders() });
export const createColeccion = (data)     => request('/colecciones',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const updateColeccion = (id, data) => request(`/colecciones/${id}`,{ method: 'PUT',    headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteColeccion = (id)       => request(`/colecciones/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getTecnicas  = ()         => request('/tecnicas',      { headers: getHeaders() });
export const createTecnica= (data)     => request('/tecnicas',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const updateTecnica= (id, data) => request(`/tecnicas/${id}`,{ method: 'PUT',    headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteTecnica= (id)       => request(`/tecnicas/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getVentas  = ()     => request('/ventas',      { headers: getHeaders(true) });
export const createVenta= (data) => request('/ventas',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteVenta= (id)   => request(`/ventas/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getTours  = ()         => request('/tours',      { headers: getHeaders() });
export const getTour   = (id)       => request(`/tours/${id}`,{ headers: getHeaders() });
export const createTour= (data)     => request('/tours',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const updateTour= (id, data) => request(`/tours/${id}`,{ method: 'PUT',    headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteTour= (id)       => request(`/tours/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getReservas  = ()     => request('/reservas',      { headers: getHeaders(true) });
export const createReserva= (data) => request('/reservas',      { method: 'POST',   headers: getHeaders(true), body: JSON.stringify(data) });
export const deleteReserva= (id)   => request(`/reservas/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getUsuarios  = ()   => request('/usuarios',      { headers: getHeaders(true) });
export const deleteUsuario= (id) => request(`/usuarios/${id}`,{ method: 'DELETE', headers: getHeaders(true) });

export const getReportePinturasCompleto  = () => request('/reportes/pinturas-completo',       { headers: getHeaders(true) });
export const getReporteVentasDetalle     = () => request('/reportes/ventas-detalle',          { headers: getHeaders(true) });
export const getReporteArtistasResumen   = () => request('/reportes/artistas-resumen',        { headers: getHeaders(true) });
export const getReporteTopArtistas       = () => request('/reportes/top-artistas-ventas',     { headers: getHeaders(true) });
export const getReporteColeccionesValor  = () => request('/reportes/colecciones-valor',       { headers: getHeaders(true) });
export const getReporteVentasPorMes      = () => request('/reportes/ventas-por-mes',          { headers: getHeaders(true) });
export const getReporteTecnicasPopulares = () => request('/reportes/tecnicas-populares',      { headers: getHeaders(true) });

async function downloadCSV(path, filename) {
  const token = localStorage.getItem('token');
  const res = await fetch(`/api${path}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  if (!res.ok) throw new Error(`Error ${res.status}`);
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}

export const exportVentasCSV   = () => downloadCSV('/exportar/ventas-csv',   'ventas_magic_bag_gallery.csv');
export const exportPinturasCSV = () => downloadCSV('/exportar/pinturas-csv', 'catalogo_pinturas.csv');
export const exportArtistasCSV = () => downloadCSV('/exportar/artistas-csv', 'artistas_resumen.csv');

export const getMisReservas = () => request('/me/reservas', { headers: getHeaders(true) });
export const getMisVentas   = () => request('/me/ventas',   { headers: getHeaders(true) });
