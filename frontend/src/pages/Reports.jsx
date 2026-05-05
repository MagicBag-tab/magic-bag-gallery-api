import { useState, useEffect } from 'react';
import {
  getReporteVentasPorMes, getReporteTopArtistas,
  getReporteColeccionesValor, getReporteTecnicasPopulares,
  exportVentasCSV, exportPinturasCSV, exportArtistasCSV
} from '../../services/api';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, LineChart, Line, CartesianGrid } from 'recharts';
import Loader from '../../components/Loader/Loader';
import styles from './Reports.module.css';

const CustomTooltip = ({ active, payload, label }) => {
  if (active && payload?.length) {
    return (
      <div className={styles.tooltip}>
        <p className={styles.tooltipLabel}>{label}</p>
        {payload.map(p => (
          <p key={p.dataKey} style={{ color: p.color }}>
            {p.name}: {typeof p.value === 'number' && p.value > 1000 ? `Q ${p.value.toLocaleString()}` : p.value}
          </p>
        ))}
      </div>
    );
  }
  return null;
};

export default function Reports() {
  const [ventasMes, setVentasMes] = useState([]);
  const [topArtistas, setTopArtistas] = useState([]);
  const [colecciones, setColecciones] = useState([]);
  const [tecnicas, setTecnicas] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      getReporteVentasPorMes(),
      getReporteTopArtistas(),
      getReporteColeccionesValor(),
      getReporteTecnicasPopulares()
    ]).then(([vm, ta, col, tec]) => {
      setVentasMes(vm);
      setTopArtistas(ta);
      setColecciones(col);
      setTecnicas(tec);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) return <Loader fullPage />;

  return (
    <div className={styles.page}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Análisis de datos</p>
        <h1 className={styles.heroTitle}>Reportes & Estadísticas</h1>
        <p className={styles.heroSub}>Métricas de ventas, artistas y colecciones de Magic Bag Gallery</p>
      </div>

      <div className={styles.exports}>
        <p className={styles.exportsLabel}>Exportar datos:</p>
        <a href={exportVentasCSV()} download className={styles.exportBtn}>Ventas CSV</a>
        <a href={exportPinturasCSV()} download className={styles.exportBtn}>Pinturas CSV</a>
        <a href={exportArtistasCSV()} download className={styles.exportBtn}>Artistas CSV</a>
      </div>

      <div className={styles.grid}>
        <div className={styles.chartCard}>
          <h2 className={styles.chartTitle}>Ingresos por mes</h2>
          <p className={styles.chartSub}>Total de ventas mensuales</p>
          <ResponsiveContainer width="100%" height={240}>
            <LineChart data={ventasMes}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.04)" />
              <XAxis dataKey="mes" tick={{ fill: '#a09880', fontSize: 11 }} />
              <YAxis tick={{ fill: '#a09880', fontSize: 11 }} />
              <Tooltip content={<CustomTooltip />} />
              <Line type="monotone" dataKey="ingresos_totales" name="Ingresos" stroke="#c9a84c" strokeWidth={2} dot={{ fill: '#c9a84c', r: 3 }} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className={styles.chartCard}>
          <h2 className={styles.chartTitle}>Top artistas por ventas</h2>
          <p className={styles.chartSub}>Ingresos generados por artista</p>
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={topArtistas} layout="vertical">
              <XAxis type="number" tick={{ fill: '#a09880', fontSize: 11 }} />
              <YAxis dataKey="nombre_completo" type="category" width={120} tick={{ fill: '#a09880', fontSize: 10 }} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="ingresos_totales" name="Ingresos" fill="#c9a84c" radius={[0, 2, 2, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className={styles.chartCard}>
          <h2 className={styles.chartTitle}>Valor de colecciones</h2>
          <p className={styles.chartSub}>Valor total por colección</p>
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={colecciones}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.04)" />
              <XAxis dataKey="nombre" tick={{ fill: '#a09880', fontSize: 10 }} />
              <YAxis tick={{ fill: '#a09880', fontSize: 11 }} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="valor_total" name="Valor total" fill="#9e7c2c" radius={[2, 2, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className={styles.chartCard}>
          <h2 className={styles.chartTitle}>Técnicas populares</h2>
          <p className={styles.chartSub}>Pinturas por técnica</p>
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={tecnicas}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.04)" />
              <XAxis dataKey="tecnica" tick={{ fill: '#a09880', fontSize: 10 }} />
              <YAxis tick={{ fill: '#a09880', fontSize: 11 }} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="total_pinturas" name="Pinturas" fill="#e2c97e" radius={[2, 2, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className={styles.tableSection}>
        <h2 className={styles.tableTitle}>Ranking de artistas</h2>
        <div className={styles.tableWrapper}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>#</th>
                <th>Artista</th>
                <th>Nacionalidad</th>
                <th>Ventas</th>
                <th>Ingresos</th>
              </tr>
            </thead>
            <tbody>
              {topArtistas.map(a => (
                <tr key={a.id_artista}>
                  <td className={styles.pos}>{a.posicion}</td>
                  <td>{a.nombre_completo}</td>
                  <td className={styles.muted}>{a.nacionalidad}</td>
                  <td className={styles.center}>{a.ventas_realizadas}</td>
                  <td className={styles.gold}>Q {Number(a.ingresos_totales).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}