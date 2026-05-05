package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ReportePinturasCompletoHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_pintura, titulo, descripcion, precio, fecha_creacion,
		       exclusiva, imagen_path, imagen_tipo, imagen_nombre,
		       artista, nacionalidad_artista, coleccion,
		       coleccion_exclusiva, tecnicas
		FROM vista_pinturas_completa
		ORDER BY artista, titulo
	`)
	if err != nil {
		http.Error(w, "Error al obtener reporte de pinturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PinturaCompleta struct {
		ID                  int     `json:"id_pintura"`
		Titulo              string  `json:"titulo"`
		Descripcion         string  `json:"descripcion"`
		Precio              float64 `json:"precio"`
		FechaCreacion       string  `json:"fecha_creacion"`
		Exclusiva           bool    `json:"exclusiva"`
		ImagenPath          string  `json:"imagen_path"`
		ImagenTipo          string  `json:"imagen_tipo"`
		ImagenNombre        string  `json:"imagen_nombre"`
		Artista             string  `json:"artista"`
		NacionalidadArtista string  `json:"nacionalidad_artista"`
		Coleccion           string  `json:"coleccion"`
		ColeccionExclusiva  *bool   `json:"coleccion_exclusiva"`
		Tecnicas            string  `json:"tecnicas"`
	}

	resultado := []PinturaCompleta{}
	for rows.Next() {
		var p PinturaCompleta
		if err := rows.Scan(
			&p.ID, &p.Titulo, &p.Descripcion, &p.Precio, &p.FechaCreacion,
			&p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
			&p.Artista, &p.NacionalidadArtista, &p.Coleccion,
			&p.ColeccionExclusiva, &p.Tecnicas,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteVentasDetalleHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_venta, fecha_venta, total_venta, cliente,
		       correo_cliente, tipo_cliente, empleado, correo_empleado, cantidad_items
		FROM vista_ventas_detalle
		ORDER BY fecha_venta DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener reporte de ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VentaDetalle struct {
		ID             int     `json:"id_venta"`
		FechaVenta     string  `json:"fecha_venta"`
		TotalVenta     float64 `json:"total_venta"`
		Cliente        string  `json:"cliente"`
		CorreoCliente  string  `json:"correo_cliente"`
		TipoCliente    string  `json:"tipo_cliente"`
		Empleado       string  `json:"empleado"`
		CorreoEmpleado string  `json:"correo_empleado"`
		CantidadItems  int     `json:"cantidad_items"`
	}

	resultado := []VentaDetalle{}
	for rows.Next() {
		var v VentaDetalle
		if err := rows.Scan(
			&v.ID, &v.FechaVenta, &v.TotalVenta, &v.Cliente,
			&v.CorreoCliente, &v.TipoCliente, &v.Empleado,
			&v.CorreoEmpleado, &v.CantidadItems,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteArtistasResumenHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_artista, nombre_completo, nacionalidad, reclutador,
		       total_pinturas, total_colecciones, valor_total_obra, precio_promedio
		FROM vista_artistas_resumen
		ORDER BY valor_total_obra DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener reporte de artistas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ArtistaResumen struct {
		ID               int     `json:"id_artista"`
		NombreCompleto   string  `json:"nombre_completo"`
		Nacionalidad     string  `json:"nacionalidad"`
		Reclutador       string  `json:"reclutador"`
		TotalPinturas    int     `json:"total_pinturas"`
		TotalColecciones int     `json:"total_colecciones"`
		ValorTotalObra   float64 `json:"valor_total_obra"`
		PrecioPromedio   float64 `json:"precio_promedio"`
	}

	resultado := []ArtistaResumen{}
	for rows.Next() {
		var a ArtistaResumen
		if err := rows.Scan(
			&a.ID, &a.NombreCompleto, &a.Nacionalidad, &a.Reclutador,
			&a.TotalPinturas, &a.TotalColecciones, &a.ValorTotalObra, &a.PrecioPromedio,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteArtistasConVentasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT a.id_artista, a.nombre_completo, a.nacionalidad
		FROM artista a
		WHERE EXISTS (
			SELECT 1
			FROM pintura p
			JOIN detalle_venta dv ON p.id_pintura = dv.id_pintura
			WHERE p.id_artista = a.id_artista
		)
		ORDER BY a.nombre_completo
	`)
	if err != nil {
		http.Error(w, "Error al obtener artistas con ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Artista struct {
		ID             int    `json:"id_artista"`
		NombreCompleto string `json:"nombre_completo"`
		Nacionalidad   string `json:"nacionalidad"`
	}

	resultado := []Artista{}
	for rows.Next() {
		var a Artista
		if err := rows.Scan(&a.ID, &a.NombreCompleto, &a.Nacionalidad); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteClientesVIPCompradoresHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT u.id_usuario, u.nombre, u.apellido, u.correo_electronico, cl.tipo_cliente
		FROM usuario u
		JOIN cliente cl ON u.id_usuario = cl.id_usuario
		WHERE cl.tipo_cliente = 'vip'
		  AND cl.id_cliente IN (
			  SELECT DISTINCT id_cliente FROM venta
		  )
		ORDER BY u.apellido, u.nombre
	`)
	if err != nil {
		http.Error(w, "Error al obtener clientes VIP compradores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ClienteVIP struct {
		IDUsuario   int    `json:"id_usuario"`
		Nombre      string `json:"nombre"`
		Apellido    string `json:"apellido"`
		Correo      string `json:"correo_electronico"`
		TipoCliente string `json:"tipo_cliente"`
	}

	resultado := []ClienteVIP{}
	for rows.Next() {
		var c ClienteVIP
		if err := rows.Scan(&c.IDUsuario, &c.Nombre, &c.Apellido, &c.Correo, &c.TipoCliente); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteVentasPorMesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT
			TO_CHAR(fecha_venta, 'YYYY-MM')   AS mes,
			COUNT(id_venta)                    AS total_ventas,
			SUM(precio)                        AS ingresos_totales,
			AVG(precio)                        AS ingreso_promedio,
			MAX(precio)                        AS venta_maxima,
			MIN(precio)                        AS venta_minima
		FROM venta
		GROUP BY TO_CHAR(fecha_venta, 'YYYY-MM')
		HAVING COUNT(id_venta) >= 1
		ORDER BY mes DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener ventas por mes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VentaMes struct {
		Mes             string  `json:"mes"`
		TotalVentas     int     `json:"total_ventas"`
		IngresosTotales float64 `json:"ingresos_totales"`
		IngresoPromedio float64 `json:"ingreso_promedio"`
		VentaMaxima     float64 `json:"venta_maxima"`
		VentaMinima     float64 `json:"venta_minima"`
	}

	resultado := []VentaMes{}
	for rows.Next() {
		var v VentaMes
		if err := rows.Scan(
			&v.Mes, &v.TotalVentas, &v.IngresosTotales,
			&v.IngresoPromedio, &v.VentaMaxima, &v.VentaMinima,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteTecnicasPopularesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT
			t.id_tecnica,
			t.nombre              AS tecnica,
			COUNT(pt.id_pintura)  AS total_pinturas,
			SUM(p.precio)         AS valor_acumulado,
			AVG(p.precio)         AS precio_promedio
		FROM tecnica t
		JOIN pintura_tecnica pt ON t.id_tecnica = pt.id_tecnica
		JOIN pintura p ON pt.id_pintura = p.id_pintura
		GROUP BY t.id_tecnica, t.nombre
		HAVING COUNT(pt.id_pintura) > 1
		ORDER BY total_pinturas DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener técnicas populares", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type TecnicaPopular struct {
		IDTecnica      int     `json:"id_tecnica"`
		Tecnica        string  `json:"tecnica"`
		TotalPinturas  int     `json:"total_pinturas"`
		ValorAcumulado float64 `json:"valor_acumulado"`
		PrecioPromedio float64 `json:"precio_promedio"`
	}

	resultado := []TecnicaPopular{}
	for rows.Next() {
		var t TecnicaPopular
		if err := rows.Scan(
			&t.IDTecnica, &t.Tecnica, &t.TotalPinturas,
			&t.ValorAcumulado, &t.PrecioPromedio,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteTopArtistasPorVentasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		WITH ingresos_por_artista AS (
			SELECT
				a.id_artista,
				a.nombre_completo,
				a.nacionalidad,
				COUNT(DISTINCT dv.id_venta)  AS ventas_realizadas,
				SUM(dv.precio_unitario * dv.cantidad) AS ingresos_totales
			FROM artista a
			JOIN pintura p ON a.id_artista = p.id_artista
			JOIN detalle_venta dv ON p.id_pintura = dv.id_pintura
			GROUP BY a.id_artista, a.nombre_completo, a.nacionalidad
		),
		ranking AS (
			SELECT *,
				RANK() OVER (ORDER BY ingresos_totales DESC) AS posicion
			FROM ingresos_por_artista
		)
		SELECT id_artista, nombre_completo, nacionalidad,
		       ventas_realizadas, ingresos_totales, posicion
		FROM ranking
		ORDER BY posicion
	`)
	if err != nil {
		http.Error(w, "Error al obtener top artistas por ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type TopArtista struct {
		IDArticle        int     `json:"id_artista"`
		NombreCompleto   string  `json:"nombre_completo"`
		Nacionalidad     string  `json:"nacionalidad"`
		VentasRealizadas int     `json:"ventas_realizadas"`
		IngresosTotales  float64 `json:"ingresos_totales"`
		Posicion         int     `json:"posicion"`
	}

	resultado := []TopArtista{}
	for rows.Next() {
		var t TopArtista
		if err := rows.Scan(
			&t.IDArticle, &t.NombreCompleto, &t.Nacionalidad,
			&t.VentasRealizadas, &t.IngresosTotales, &t.Posicion,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ReporteColeccionesValorHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		WITH valor_colecciones AS (
			SELECT
				c.id_coleccion,
				c.nombre,
				c.exclusiva,
				c.fecha_lanzamiento,
				COUNT(p.id_pintura)  AS total_pinturas,
				SUM(p.precio)        AS valor_total,
				AVG(p.precio)        AS precio_promedio
			FROM coleccion c
			LEFT JOIN pintura p ON c.id_coleccion = p.id_coleccion
			GROUP BY c.id_coleccion, c.nombre, c.exclusiva, c.fecha_lanzamiento
		)
		SELECT
			id_coleccion, nombre, exclusiva, fecha_lanzamiento,
			total_pinturas, valor_total, precio_promedio,
			RANK() OVER (ORDER BY valor_total DESC NULLS LAST) AS ranking_valor
		FROM valor_colecciones
		ORDER BY ranking_valor
	`)
	if err != nil {
		http.Error(w, "Error al obtener valor de colecciones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ColeccionValor struct {
		IDColeccion      int     `json:"id_coleccion"`
		Nombre           string  `json:"nombre"`
		Exclusiva        bool    `json:"exclusiva"`
		FechaLanzamiento string  `json:"fecha_lanzamiento"`
		TotalPinturas    int     `json:"total_pinturas"`
		ValorTotal       float64 `json:"valor_total"`
		PrecioPromedio   float64 `json:"precio_promedio"`
		RankingValor     int     `json:"ranking_valor"`
	}

	resultado := []ColeccionValor{}
	for rows.Next() {
		var c ColeccionValor
		if err := rows.Scan(
			&c.IDColeccion, &c.Nombre, &c.Exclusiva, &c.FechaLanzamiento,
			&c.TotalPinturas, &c.ValorTotal, &c.PrecioPromedio, &c.RankingValor,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ExportarVentasCSVHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_venta, fecha_venta, total_venta, cliente,
		       correo_cliente, tipo_cliente, empleado, correo_empleado, cantidad_items
		FROM vista_ventas_detalle
		ORDER BY fecha_venta DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener datos para exportar", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=ventas_magic_bag_gallery.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"ID Venta", "Fecha", "Total (Q)", "Cliente", "Correo Cliente",
		"Tipo Cliente", "Empleado", "Correo Empleado", "Cantidad Items",
	})

	for rows.Next() {
		var (
			id             int
			fechaVenta     string
			totalVenta     float64
			cliente        string
			correoCliente  string
			tipoCliente    string
			empleado       string
			correoEmpleado string
			cantidadItems  int
		)
		if err := rows.Scan(
			&id, &fechaVenta, &totalVenta, &cliente,
			&correoCliente, &tipoCliente, &empleado,
			&correoEmpleado, &cantidadItems,
		); err != nil {
			continue
		}
		writer.Write([]string{
			strconv.Itoa(id),
			fechaVenta,
			strconv.FormatFloat(totalVenta, 'f', 2, 64),
			cliente,
			correoCliente,
			tipoCliente,
			empleado,
			correoEmpleado,
			strconv.Itoa(cantidadItems),
		})
	}
}

func ExportarPinturasCSVHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_pintura, titulo, precio, fecha_creacion,
		       exclusiva, artista, coleccion, tecnicas
		FROM vista_pinturas_completa
		ORDER BY artista, titulo
	`)
	if err != nil {
		http.Error(w, "Error al obtener datos para exportar", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=catalogo_pinturas.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"ID", "Título", "Precio (Q)", "Fecha Creación",
		"Exclusiva", "Artista", "Colección", "Técnicas",
	})

	for rows.Next() {
		var (
			id            int
			titulo        string
			precio        float64
			fechaCreacion string
			exclusiva     bool
			artista       string
			coleccion     string
			tecnicas      string
		)
		if err := rows.Scan(
			&id, &titulo, &precio, &fechaCreacion,
			&exclusiva, &artista, &coleccion, &tecnicas,
		); err != nil {
			continue
		}
		exclusivaStr := "No"
		if exclusiva {
			exclusivaStr = "Sí"
		}
		writer.Write([]string{
			strconv.Itoa(id),
			titulo,
			strconv.FormatFloat(precio, 'f', 2, 64),
			fechaCreacion,
			exclusivaStr,
			artista,
			coleccion,
			tecnicas,
		})
	}
}

func ExportarArtistasCSVHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_artista, nombre_completo, nacionalidad, reclutador,
		       total_pinturas, total_colecciones, valor_total_obra, precio_promedio
		FROM vista_artistas_resumen
		ORDER BY valor_total_obra DESC
	`)
	if err != nil {
		http.Error(w, "Error al obtener datos para exportar", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=artistas_resumen.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"ID", "Artista", "Nacionalidad", "Reclutador",
		"Total Pinturas", "Total Colecciones", "Valor Total Obra (Q)", "Precio Promedio (Q)",
	})

	for rows.Next() {
		var (
			id               int
			nombreCompleto   string
			nacionalidad     string
			reclutador       string
			totalPinturas    int
			totalColecciones int
			valorTotalObra   float64
			precioPromedio   float64
		)
		if err := rows.Scan(
			&id, &nombreCompleto, &nacionalidad, &reclutador,
			&totalPinturas, &totalColecciones, &valorTotalObra, &precioPromedio,
		); err != nil {
			continue
		}
		writer.Write([]string{
			strconv.Itoa(id),
			nombreCompleto,
			nacionalidad,
			reclutador,
			strconv.Itoa(totalPinturas),
			strconv.Itoa(totalColecciones),
			strconv.FormatFloat(valorTotalObra, 'f', 2, 64),
			strconv.FormatFloat(precioPromedio, 'f', 2, 64),
		})
	}
}

func ReporteVentasPorAnioHandler(w http.ResponseWriter, r *http.Request) {
	anio, err := strconv.Atoi(mux.Vars(r)["anio"])
	if err != nil || anio < 2000 || anio > 2100 {
		http.Error(w, "Año inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT
			TO_CHAR(fecha_venta, 'YYYY-MM')   AS mes,
			COUNT(id_venta)                    AS total_ventas,
			SUM(precio)                        AS ingresos_totales,
			AVG(precio)                        AS ingreso_promedio,
			MAX(precio)                        AS venta_maxima,
			MIN(precio)                        AS venta_minima
		FROM venta
		WHERE EXTRACT(YEAR FROM fecha_venta) = $1
		GROUP BY TO_CHAR(fecha_venta, 'YYYY-MM')
		HAVING COUNT(id_venta) >= 1
		ORDER BY mes
	`, anio)
	if err != nil {
		http.Error(w, "Error al obtener ventas por año", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VentaMes struct {
		Mes             string  `json:"mes"`
		TotalVentas     int     `json:"total_ventas"`
		IngresosTotales float64 `json:"ingresos_totales"`
		IngresoPromedio float64 `json:"ingreso_promedio"`
		VentaMaxima     float64 `json:"venta_maxima"`
		VentaMinima     float64 `json:"venta_minima"`
	}

	resultado := []VentaMes{}
	for rows.Next() {
		var v VentaMes
		if err := rows.Scan(
			&v.Mes, &v.TotalVentas, &v.IngresosTotales,
			&v.IngresoPromedio, &v.VentaMaxima, &v.VentaMinima,
		); err != nil {
			http.Error(w, "Error al leer datos", http.StatusInternalServerError)
			return
		}
		resultado = append(resultado, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}
