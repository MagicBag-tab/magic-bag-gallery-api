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

| Servicio   | URL                        |
|------------|----------------------------|
| Frontend   | http://localhost:3000      |
| Backend    | http://localhost:8888      |
| Base datos | `localhost:5432` (proy2db) |

La base de datos se inicializa automáticamente con el esquema DDL y los datos de prueba la primera vez que se levanta el contenedor.

> Para reiniciar desde cero (borrar todos los datos): `docker compose down -v && docker compose up`

---

## Variables de entorno

```env
POSTGRES_USER=proy2
POSTGRES_PASSWORD=secret
POSTGRES_DB=proy2db
DB_HOST=database
DB_PORT=5432
JWT_SECRET=your_super_secret_jwt_key_here_change_in_production
```

> Las credenciales `proy2` / `secret` son requeridas para la calificación del proyecto.

---

## Credenciales de usuarios de prueba

**Contraseña de todos los usuarios: `secret`**

### Empleados — acceso completo al panel `/admin`

| Nombre           | Correo electrónico           | Tipo       |
|------------------|------------------------------|------------|
| Ana Solís        | ana.solis@magicbag.gt        | guia       |
| Roberto Lima     | roberto.lima@magicbag.gt     | asesor     |
| Patricia Aguilar | patricia.aguilar@magicbag.gt | reclutador |
| Miguel Ramos     | miguel.ramos@magicbag.gt     | guia       |
| Carmen Cifuentes | carmen.cifuentes@magicbag.gt | asesor     |

### Clientes

| Nombre          | Correo electrónico        | Tipo    |
|-----------------|---------------------------|---------|
| María Pérez     | maria.perez@gmail.com     | vip     |
| Carlos Méndez   | carlos.mendez@gmail.com   | regular |
| Lucía Hernández | lucia.hernandez@gmail.com | vip     |
| Sofía Ramírez   | sofia.ramirez@gmail.com   | vip     |
| Diego López     | diego.lopez@gmail.com     | regular |

---

## Páginas de la aplicación

| Ruta           | Descripción                                       | Acceso        |
|----------------|---------------------------------------------------|---------------|
| `/catalogo`    | Catálogo completo de pinturas con filtros          | Público       |
| `/artistas`    | Lista de artistas con resumen de su obra           | Público       |
| `/colecciones` | Colecciones disponibles en la galería              | Público       |
| `/tours`       | Tours guiados disponibles con precio y horario     | Público       |
| `/reportes`    | Reportes con gráficas y exportación CSV            | Público       |
| `/login`       | Inicio de sesión                                  | Público       |
| `/register`    | Registro de nuevos clientes                        | Público       |
| `/admin`       | Panel de administración con CRUD completo          | Solo empleado |

---

## Estructura del proyecto

```
magic-bag-gallery-api/
├── backend/
│   ├── internal/
│   │   ├── handlers/        # Handlers HTTP por entidad
│   │   ├── middleware/       # JWT y control de roles
│   │   └── models/           # Structs de datos
│   ├── main.go
│   ├── go.mod / go.sum
│   └── Dockerfile
├── db/
│   ├── ddl_magic_bag_gallery.sql                  # Esquema + índices + vistas
│   └── dml_datos_iniciales_magic_bag_gallery.sql  # Datos de prueba (45 usuarios)
├── frontend/
│   ├── src/
│   │   ├── api/             # Funciones fetch al backend
│   │   ├── components/      # Componentes reutilizables
│   │   ├── context/         # AuthContext (JWT)
│   │   └── pages/           # Vistas de la aplicación
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Diseño de base de datos

### Vistas SQL utilizadas por el backend

| Vista                     | Descripción                                            |
|---------------------------|--------------------------------------------------------|
| `vista_pinturas_completa` | Pintura + artista + colección + técnicas (STRING_AGG)  |
| `vista_ventas_detalle`    | Venta + cliente + empleado + conteo de ítems           |
| `vista_artistas_resumen`  | Artista + reclutador + totales de obras y valor        |

### Índices definidos

| Índice                  | Tabla   | Columna             | Justificación                   |
|-------------------------|---------|---------------------|---------------------------------|
| `idx_usuario_correo`    | usuario | correo_electronico  | Búsqueda en login por correo    |
| `idx_pintura_artista`   | pintura | id_artista          | Filtrar obras por artista       |
| `idx_venta_cliente`     | venta   | id_cliente          | Historial de ventas por cliente |
| `idx_envio_venta`       | envio   | id_venta            | Consulta de envíos por venta    |
| `idx_pintura_coleccion` | pintura | id_coleccion        | Filtrar pinturas por colección  |

---

## Características técnicas

- **SQL explícito** — sin ORM; todas las queries escritas a mano con `database/sql`
- **Transacciones** — `BEGIN / COMMIT / ROLLBACK` explícito en operaciones críticas
- **Vistas SQL** — 3 vistas usadas por el backend para alimentar la UI
- **CTEs** — queries con `WITH` y `RANK()` para rankings de artistas y colecciones
- **Subqueries** — `EXISTS` e `IN` en reportes de artistas y clientes VIP
- **GROUP BY + HAVING** — reportes de ventas mensuales y técnicas populares
- **Autenticación JWT** — roles `cliente` / `empleado` con middleware
- **Exportación CSV** — descarga directa desde la UI (ventas, pinturas, artistas)
- **Manejo de errores** — mensajes descriptivos en frontend y backend

