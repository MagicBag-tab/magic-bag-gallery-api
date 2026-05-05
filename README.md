# Magic Bag Gallery

Aplicación web para gestionar el inventario y las ventas de una galería de arte. Desarrollada con Go (backend), React (frontend) y PostgreSQL (base de datos). Todo el stack se levanta con Docker Compose.

---

## Tecnologías

| Capa       | Tecnología                          |
|------------|-------------------------------------|
| Frontend   | React 18, React Router v6, Recharts |
| Backend    | Go 1.25, Gorilla Mux, JWT           |
| Base datos | PostgreSQL 16                       |
| Despliegue | Docker & Docker Compose             |

---

## Requisitos previos

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (incluye Docker Compose)

No se necesita instalar Go, Node.js ni PostgreSQL de forma local.

---

## Levantar el proyecto

```bash
# 1. Clonar el repositorio
git clone <url-del-repositorio>
cd magic-bag-gallery-api

# 2. Crear el archivo de variables de entorno
cp .env.example .env

# 3. Levantar todos los servicios
docker compose up
```

Al terminar de iniciar:

| Servicio  | URL                          |
|-----------|------------------------------|
| Frontend  | http://localhost:3000        |
| Backend   | http://localhost:8888        |
| Base datos| `localhost:5432` (proy2db)   |

La base de datos se inicializa automáticamente con el esquema DDL y los datos de prueba la primera vez que se levanta el contenedor.

---

## Variables de entorno

El archivo `.env.example` incluye todas las variables requeridas:

```env
POSTGRES_USER=proy2
POSTGRES_PASSWORD=secret
POSTGRES_DB=proy2db
DB_HOST=database
DB_PORT=5432
JWT_SECRET=your_super_secret_jwt_key_here_change_in_production
```

> Las credenciales `proy2` / `secret` son requeridas para la calificación del proyecto y no deben modificarse.

---

## Estructura del proyecto

```
magic-bag-gallery-api/
├── backend/
│   ├── internal/
│   │   ├── handlers/       # Handlers HTTP por entidad
│   │   ├── middleware/      # JWT y control de roles
│   │   └── models/          # Structs de datos
│   ├── main.go              # Punto de entrada y rutas
│   ├── go.mod / go.sum
│   └── Dockerfile
├── db/
│   ├── ddl_magic_bag_gallery.sql            # Esquema, índices y vistas
│   └── dml_datos_iniciales_magic_bag_gallery.sql  # Datos de prueba
├── frontend/
│   ├── src/
│   │   ├── api/             # Funciones fetch al backend
│   │   ├── components/      # Componentes reutilizables
│   │   ├── context/         # AuthContext (JWT)
│   │   └── pages/           # Vistas de la aplicación
│   ├── package.json
│   └── Dockerfile
├── uploads/
│   └── pinturas/            # Imágenes de obras
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Páginas de la aplicación

| Ruta          | Descripción                                              | Acceso        |
|---------------|----------------------------------------------------------|---------------|
| `/catalogo`   | Catálogo completo de pinturas con filtros                | Público       |
| `/artistas`   | Lista de artistas con resumen de su obra                 | Público       |
| `/colecciones`| Colecciones disponibles en la galería                    | Público       |
| `/tours`      | Tours guiados disponibles con precio y horario           | Público       |
| `/reportes`   | Reportes con gráficas: ventas, artistas, técnicas        | Público       |
| `/login`      | Inicio de sesión (cliente o empleado)                    | Público       |
| `/register`   | Registro de nuevos clientes                              | Público       |
| `/admin`      | Panel de administración con CRUD completo                | Solo empleado |

---

## Credenciales de prueba

Para ingresar como **empleado** y acceder al panel de administración:

```
Correo:     ana.solis@magicbag.gt
Contraseña: (la del hash en el DML — reemplazar con contraseña real al hacer setup)
```

> Los hashes del DML son placeholders. Para crear un usuario de prueba funcional, usar el endpoint `POST /api/auth/register/empleado` con un JWT de empleado existente, o insertar un hash bcrypt válido directamente en la base de datos.

---

## API — Endpoints principales

### Autenticación

| Método | Ruta                          | Descripción                     |
|--------|-------------------------------|---------------------------------|
| POST   | `/api/login`                  | Login, devuelve JWT             |
| POST   | `/api/register/cliente`       | Registro de cliente             |
| POST   | `/api/auth/register/empleado` | Registro de empleado (JWT req.) |

### Catálogo (público)

| Método | Ruta                                  | Descripción                      |
|--------|---------------------------------------|----------------------------------|
| GET    | `/api/pinturas`                       | Listar todas las pinturas        |
| GET    | `/api/pinturas/{id}`                  | Detalle de una pintura           |
| GET    | `/api/pinturas/artista/{id_artista}`  | Pinturas por artista             |
| GET    | `/api/pinturas/coleccion/{id_coleccion}` | Pinturas por colección        |
| GET    | `/api/pinturas/tecnica/{id_tecnica}`  | Pinturas por técnica             |
| GET    | `/api/artistas`                       | Listar artistas                  |
| GET    | `/api/colecciones`                    | Listar colecciones               |
| GET    | `/api/tecnicas`                       | Listar técnicas                  |
| GET    | `/api/tours`                          | Listar tours                     |

### Reportes (público)

| Método | Ruta                                    | Tipo de consulta SQL             |
|--------|-----------------------------------------|----------------------------------|
| GET    | `/api/reportes/pinturas-completo`       | VIEW + JOIN múltiple             |
| GET    | `/api/reportes/ventas-detalle`          | VIEW + JOIN múltiple             |
| GET    | `/api/reportes/artistas-resumen`        | VIEW + JOIN múltiple             |
| GET    | `/api/reportes/artistas-con-ventas`     | Subquery `EXISTS`                |
| GET    | `/api/reportes/clientes-vip-compradores`| Subquery `IN`                    |
| GET    | `/api/reportes/ventas-por-mes`          | `GROUP BY` + `HAVING` + agregación |
| GET    | `/api/reportes/ventas-por-mes/{anio}`   | `GROUP BY` filtrado por año      |
| GET    | `/api/reportes/tecnicas-populares`      | `GROUP BY` + `HAVING` + agregación |
| GET    | `/api/reportes/top-artistas-ventas`     | CTE (`WITH`) + `RANK()`          |
| GET    | `/api/reportes/colecciones-valor`       | CTE (`WITH`) + `RANK()`          |

### Exportación CSV (público)

| Método | Ruta                          | Descripción                    |
|--------|-------------------------------|--------------------------------|
| GET    | `/api/exportar/ventas-csv`    | Reporte de ventas en CSV       |
| GET    | `/api/exportar/pinturas-csv`  | Catálogo de pinturas en CSV    |
| GET    | `/api/exportar/artistas-csv`  | Resumen de artistas en CSV     |

### CRUD administrativo (requiere JWT de empleado)

| Método          | Ruta                          | Entidad         |
|-----------------|-------------------------------|-----------------|
| GET/POST        | `/api/pinturas`               | Pinturas        |
| GET/PUT/DELETE  | `/api/pinturas/{id}`          | Pinturas        |
| GET/POST        | `/api/artistas`               | Artistas        |
| GET/PUT/DELETE  | `/api/artistas/{id}`          | Artistas        |
| GET/POST        | `/api/colecciones`            | Colecciones     |
| GET/PUT/DELETE  | `/api/colecciones/{id}`       | Colecciones     |
| GET/POST        | `/api/tecnicas`               | Técnicas        |
| GET/PUT/DELETE  | `/api/tecnicas/{id}`          | Técnicas        |
| GET/POST        | `/api/ventas`                 | Ventas          |
| GET/PUT/DELETE  | `/api/ventas/{id}`            | Ventas          |
| GET/POST        | `/api/detalles-venta`         | Detalles venta  |
| GET/PUT/DELETE  | `/api/detalles-venta/{id}`    | Detalles venta  |
| GET/POST        | `/api/envios`                 | Envíos          |
| GET/PUT/DELETE  | `/api/envios/{id}`            | Envíos          |
| GET/POST        | `/api/tours`                  | Tours           |
| GET/PUT/DELETE  | `/api/tours/{id}`             | Tours           |
| GET/POST        | `/api/usuarios`               | Usuarios        |
| GET/PUT/DELETE  | `/api/usuarios/{id}`          | Usuarios        |

### Reservas (requiere JWT)

| Método          | Ruta                  | Descripción          |
|-----------------|-----------------------|----------------------|
| GET/POST        | `/api/reservas`       | Listar / crear reserva |
| GET/PUT/DELETE  | `/api/reservas/{id}`  | Detalle / editar / eliminar |

---

## Diseño de base de datos

### Entidades principales

- **usuario** — base común para clientes y empleados (nombre, correo, teléfono, contraseña bcrypt)
- **cliente** — extiende usuario con tipo de cliente (`vip` / `regular`)
- **empleado** — extiende usuario con tipo de empleado (`guia` / `asesor` / `reclutador`)
- **artista** — artista con nacionalidad y empleado reclutador
- **pintura** — obra con precio, imagen, colección y técnicas asociadas
- **coleccion** — agrupación temática de pinturas, puede ser exclusiva
- **tecnica** — técnica artística; relación N:M con pintura via `pintura_tecnica`
- **tour** — tour guiado con fechas, horario y precio
- **venta** — transacción entre cliente y empleado
- **detalle_venta** — líneas de cada venta (pintura, cantidad, precio unitario)
- **envio** — despacho asociado a una venta
- **direccion** — direcciones de envío por cliente
- **cliente_tour** — reservas de clientes a tours

### Vistas SQL

| Vista                    | Descripción                                                 |
|--------------------------|-------------------------------------------------------------|
| `vista_pinturas_completa`| Pintura + artista + colección + técnicas (STRING_AGG)       |
| `vista_ventas_detalle`   | Venta + cliente + empleado + conteo de ítems                |
| `vista_artistas_resumen` | Artista + reclutador + totales de obras y valor             |

### Índices definidos

| Índice                  | Tabla    | Columna               | Justificación                         |
|-------------------------|----------|-----------------------|---------------------------------------|
| `idx_usuario_correo`    | usuario  | correo_electronico    | Búsqueda en login por correo          |
| `idx_pintura_artista`   | pintura  | id_artista            | Filtrar obras por artista             |
| `idx_venta_cliente`     | venta    | id_cliente            | Historial de ventas por cliente       |
| `idx_envio_venta`       | envio    | id_venta              | Consulta de envíos asociados a venta  |
| `idx_pintura_coleccion` | pintura  | id_coleccion          | Filtrar pinturas por colección        |

---

## Características técnicas

- **SQL explícito** — sin ORM; todas las queries escritas a mano con `database/sql`
- **Transacciones** — `BEGIN / COMMIT / ROLLBACK` explícito en operaciones críticas (ej. eliminación en cascada de ventas)
- **Autenticación JWT** — login con roles `cliente` / `empleado`; rutas protegidas por middleware
- **Exportación CSV** — descarga directa desde la UI para ventas, pinturas y artistas
- **Manejo de errores** — mensajes descriptivos en frontend y backend para validaciones y fallos de BD

---

## Notas

- Al reiniciar los contenedores con `docker compose down -v` se elimina el volumen de la base de datos y los datos se reinicializan desde cero con el DML.
- Para conservar los datos entre reinicios, usar `docker compose down` (sin `-v`).
- El backend recarga el código en cada inicio porque usa `go run main.go`; en producción se recomienda compilar el binario.