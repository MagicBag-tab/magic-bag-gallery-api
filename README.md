# Magic Bag Gallery

Aplicación web completa para gestionar el inventario y las ventas de una galería de arte contemporáneo. Desarrollada con Go (backend), React (frontend) y PostgreSQL (base de datos). Todo el stack se levanta con un único comando Docker Compose.

---

## Tecnologías

| Capa        | Tecnología                                          |
|-------------|-----------------------------------------------------|
| Frontend    | React 18, React Router v6, Recharts, CSS Modules    |
| Backend     | Go 1.25, Gorilla Mux, GORM v2, gorilla/sessions     |
| Base datos  | PostgreSQL 16                                       |
| Despliegue  | Docker & Docker Compose                             |
| Calidad     | ESLint, Vitest, @testing-library/react              |

---

## Proyecto 3

Esta rama `proyecto-3` extiende el Proyecto 2 con:

- CRUD del backend migrado a GORM v2 (`gorm.io/gorm` + `gorm.io/driver/postgres`).
- Autenticacion con cookie de sesion (`gorilla/sessions`) y rutas protegidas por rol.
- Funciones PL/pgSQL para operaciones criticas: registro de clientes/empleados, reservas, ventas y eliminacion de pinturas.
- 5 roles de base de datos: `mbg_catalogo`, `mbg_cliente`, `mbg_guia`, `mbg_asesor`, `mbg_reclutador`.
- Script SQL de inicializacion en `db/zz_proyecto3_roles_procedures.sql`; Docker Compose ya monta `./db` en `/docker-entrypoint-initdb.d`.

> Si ya existe un volumen de PostgreSQL de una corrida anterior, reiniciar con `docker compose down -v && docker compose up` para ejecutar los scripts nuevos.

---

## Requisitos previos

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (incluye Docker Compose v2)

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

Al terminar de iniciar (puede tardar ~30 s la primera vez mientras descarga imágenes):

| Servicio    | URL                         |
|-------------|-----------------------------|
| Frontend    | http://localhost:3000       |
| Backend API | http://localhost:8888       |
| Base datos  | `localhost:5432` (proy2db)  |

La base de datos se inicializa automáticamente con el esquema DDL y los datos de prueba la primera vez que se levanta el contenedor.

> **Reiniciar desde cero** (borra todos los datos): `docker compose down -v && docker compose up`

---

## Variables de entorno

Copiar `.env.example` a `.env`. Las variables requeridas son:

```env
POSTGRES_USER=proy2
POSTGRES_PASSWORD=secret
POSTGRES_DB=proy2db
DB_HOST=database
DB_PORT=5432
JWT_SECRET=your_super_secret_jwt_key_here_change_in_production
SESSION_SECRET=your_super_secret_session_key_here_change_in_production
```

---

## Credenciales de usuarios de prueba

**Contraseña de todos los usuarios: `secret`**

### Empleados — acceso al panel `/admin` y `/reportes`

| Nombre           | Correo                        | Tipo de empleado |
|------------------|-------------------------------|------------------|
| Ana Solís        | ana.solis@magicbag.gt         | guia             |
| Roberto Lima     | roberto.lima@magicbag.gt      | asesor           |
| Patricia Aguilar | patricia.aguilar@magicbag.gt  | reclutador       |
| Miguel Ramos     | miguel.ramos@magicbag.gt      | guia             |
| Carmen Cifuentes | carmen.cifuentes@magicbag.gt  | asesor           |

### Clientes — acceso a `/mi-cuenta`

| Nombre           | Correo                     | Tipo    |
|------------------|----------------------------|---------|
| María Pérez      | maria.perez@gmail.com      | vip     |
| Carlos Méndez    | carlos.mendez@gmail.com    | regular |
| Lucía Hernández  | lucia.hernandez@gmail.com  | vip     |
| Sofía Ramírez    | sofia.ramirez@gmail.com    | vip     |
| Diego López      | diego.lopez@gmail.com      | regular |

---

## Páginas de la aplicación

| Ruta           | Descripción                                          | Acceso              |
|----------------|------------------------------------------------------|---------------------|
| `/catalogo`    | Catálogo completo de pinturas con filtros            | Público             |
| `/artistas`    | Lista de artistas con resumen de su obra             | Público             |
| `/colecciones` | Colecciones disponibles en la galería                | Público             |
| `/tours`       | Tours guiados con precio, horario y reserva          | Público / Auth      |
| `/login`       | Inicio de sesión                                     | Público             |
| `/register`    | Registro de nuevos clientes                          | Público             |
| `/reportes`    | Reportes con gráficas y exportación CSV              | Solo empleados      |
| `/admin`       | Panel CRUD adaptado al tipo de empleado              | Solo empleados      |
| `/mi-cuenta`   | Reservas y compras del cliente                       | Solo clientes       |

### Panel de administración por rol

| Tipo empleado | Tabs disponibles en `/admin`              |
|---------------|-------------------------------------------|
| `guia`        | Tours, Reservas                           |
| `asesor`      | Ventas, Usuarios                          |
| `reclutador`  | Artistas, Colecciones                     |
| Sin tipo      | Todos: Pinturas, Artistas, Colecciones, Técnicas, Ventas, Tours, Reservas, Usuarios |

---

## Comandos de desarrollo (frontend)

```bash
cd frontend
npm install

npm start        # Servidor de desarrollo en http://localhost:3000
npm run lint     # ESLint — debe terminar sin errores
npm test         # Vitest — ejecuta los tests unitarios
npm run build    # Build de producción
```

---

## API REST — Documentación de endpoints

> **Base URL**: `http://localhost:8888/api`
>
> Autenticacion: cookie de sesion HTTP-only creada por `/login`. El frontend envia la cookie automaticamente en las peticiones protegidas.

### Autenticación

| Método | Ruta                       | Auth | Descripción                        | Body (JSON)                                                                 |
|--------|----------------------------|------|------------------------------------|-----------------------------------------------------------------------------|
| POST   | `/login`                   | No   | Inicia sesion y crea cookie HTTP-only | `{ "correo_electronico": "", "contrasena": "" }`                         |
| POST   | `/logout`                  | No   | Cierra la sesion del navegador     | N/A                                                                         |
| GET    | `/session`                 | Sesion | Devuelve la sesion activa        | N/A                                                                         |
| POST   | `/register/cliente`        | No   | Registra un nuevo cliente          | `{ "nombre", "apellido", "correo_electronico", "telefono", "contrasena" }`  |
| POST   | `/auth/register/empleado`  | Empleado | Registra un nuevo empleado    | `{ "nombre", "apellido", "correo_electronico", "telefono", "contrasena", "tipo_empleado" }` |

**Respuesta de `/login`:**
```json
{ "role": "empleado", "id_usuario": 11, "nombre": "Ana", "tipo_empleado": "guia" }
```

---

### Pinturas

| Método | Ruta                                    | Auth     | Descripción                              |
|--------|-----------------------------------------|----------|------------------------------------------|
| GET    | `/pinturas`                             | No       | Lista todas las pinturas con artista y colección |
| GET    | `/pinturas/{id}`                        | No       | Detalle de una pintura                   |
| GET    | `/pinturas/artista/{id_artista}`        | No       | Pinturas filtradas por artista           |
| GET    | `/pinturas/coleccion/{id_coleccion}`    | No       | Pinturas filtradas por colección         |
| GET    | `/pinturas/tecnica/{id_tecnica}`        | No       | Pinturas filtradas por técnica           |
| POST   | `/pinturas`                             | Empleado | Crear nueva pintura                      |
| PUT    | `/pinturas/{id}`                        | Empleado | Actualizar pintura                       |
| DELETE | `/pinturas/{id}`                        | Empleado | Eliminar pintura y sus técnicas          |

**Body para POST/PUT `/pinturas`:**
```json
{
  "titulo": "Untitled",
  "descripcion": "...",
  "fecha_creacion": "1981-01-01",
  "precio": 95000.00,
  "exclusiva": true,
  "imagen_path": "/uploads/pinturas/img.jpg",
  "imagen_tipo": "image/jpeg",
  "imagen_nombre": "img.jpg",
  "id_artista": 1,
  "id_coleccion": 1,
  "tecnicas": [1, 4]
}
```

---

### Artistas

| Método | Ruta              | Auth     | Descripción                                      |
|--------|-------------------|----------|--------------------------------------------------|
| GET    | `/artistas`       | No       | Lista todos los artistas                         |
| GET    | `/artistas/{id}`  | No       | Detalle con pinturas y colecciones del artista   |
| POST   | `/artistas`       | Empleado | Crear artista                                    |
| PUT    | `/artistas/{id}`  | Empleado | Actualizar artista                               |
| DELETE | `/artistas/{id}`  | Empleado | Eliminar artista y sus pinturas                  |

**Body para POST/PUT `/artistas`:**
```json
{
  "nombre_completo": "Jean-Michel Basquiat",
  "nacionalidad": "Estadounidense",
  "id_reclutador": 3,
  "id_pinturas": [1, 2, 3]
}
```

---

### Colecciones

| Método | Ruta                  | Auth     | Descripción                    |
|--------|-----------------------|----------|--------------------------------|
| GET    | `/colecciones`        | No       | Lista todas las colecciones    |
| GET    | `/colecciones/{id}`   | No       | Detalle con pinturas incluidas |
| POST   | `/colecciones`        | Empleado | Crear colección                |
| PUT    | `/colecciones/{id}`   | Empleado | Actualizar colección           |
| DELETE | `/colecciones/{id}`   | Empleado | Eliminar colección             |

**Body para POST/PUT `/colecciones`:**
```json
{
  "nombre": "Neo-Expresionismo",
  "descripcion": "...",
  "exclusiva": true,
  "fecha_lanzamiento": "2024-01-15",
  "id_pinturas": [1, 2, 3]
}
```

---

### Técnicas

| Método | Ruta               | Auth     | Descripción          |
|--------|--------------------|----------|----------------------|
| GET    | `/tecnicas`        | No       | Lista todas          |
| GET    | `/tecnicas/{id}`   | No       | Detalle de una       |
| POST   | `/tecnicas`        | Empleado | Crear técnica        |
| PUT    | `/tecnicas/{id}`   | Empleado | Actualizar técnica   |
| DELETE | `/tecnicas/{id}`   | Empleado | Eliminar técnica     |

---

### Tours y Reservas

| Método | Ruta               | Auth     | Descripción                              |
|--------|--------------------|----------|------------------------------------------|
| GET    | `/tours`           | No       | Lista todos los tours con nombre de guía |
| GET    | `/tours/{id}`      | No       | Detalle de un tour                       |
| POST   | `/tours`           | Empleado | Crear tour                               |
| PUT    | `/tours/{id}`      | Empleado | Actualizar tour                          |
| DELETE | `/tours/{id}`      | Empleado | Eliminar tour y sus reservas             |
| GET    | `/reservas`        | Sesion   | Lista todas las reservas                 |
| GET    | `/reservas/{id}`   | Sesion   | Detalle de una reserva                   |
| POST   | `/reservas`        | Sesion   | Crear reserva                            |
| PUT    | `/reservas/{id}`   | Sesion   | Actualizar reserva                       |
| DELETE | `/reservas/{id}`   | Sesion   | Eliminar reserva                         |

**Body para POST `/tours`:**
```json
{
  "id_guia": 1,
  "nombre": "Neo-Expresionismo y Basquiat",
  "descripcion": "...",
  "fecha_inicio": "2025-01-10",
  "fecha_fin": "2025-01-10",
  "horario": "10:00 - 12:00",
  "precio": "150.00"
}
```

**Body para POST `/reservas`:**
```json
{
  "id_cliente": 1,
  "id_tour": 1,
  "fecha_reserva": "2025-01-05"
}
```

---

### Ventas y Detalles

| Método | Ruta                          | Auth     | Descripción                        |
|--------|-------------------------------|----------|------------------------------------|
| GET    | `/ventas`                     | Empleado | Lista todas las ventas             |
| GET    | `/ventas/{id}`                | Empleado | Detalle de una venta               |
| POST   | `/ventas`                     | Empleado | Crear venta                        |
| PUT    | `/ventas/{id}`                | Empleado | Actualizar venta                   |
| DELETE | `/ventas/{id}`                | Empleado | Eliminar venta y detalles/envíos   |
| GET    | `/detalles-venta`             | Empleado | Lista todos los ítems de venta     |
| GET    | `/detalles-venta/{id}`        | Empleado | Detalle de un ítem                 |
| GET    | `/ventas/{id_venta}/detalles` | Empleado | Ítems de una venta específica      |
| POST   | `/detalles-venta`             | Empleado | Agregar ítem a una venta           |
| PUT    | `/detalles-venta/{id}`        | Empleado | Actualizar ítem                    |
| DELETE | `/detalles-venta/{id}`        | Empleado | Eliminar ítem                      |

---

### Envíos

| Método | Ruta            | Auth     | Descripción        |
|--------|-----------------|----------|--------------------|
| GET    | `/envios`       | Empleado | Lista envíos       |
| GET    | `/envios/{id}`  | Empleado | Detalle de envío   |
| POST   | `/envios`       | Empleado | Crear envío        |
| PUT    | `/envios/{id}`  | Empleado | Actualizar envío   |
| DELETE | `/envios/{id}`  | Empleado | Eliminar envío     |

---

### Usuarios

| Método | Ruta               | Auth     | Descripción          |
|--------|--------------------|----------|----------------------|
| GET    | `/usuarios`        | Empleado | Lista todos          |
| GET    | `/usuarios/{id}`   | Empleado | Detalle de uno       |
| POST   | `/usuarios`        | Empleado | Crear usuario        |
| PUT    | `/usuarios/{id}`   | Empleado | Actualizar usuario   |
| DELETE | `/usuarios/{id}`   | Empleado | Eliminar usuario     |

---

### Endpoints personales (requieren sesion del usuario)

| Método | Ruta                   | Auth | Descripción                                    |
|--------|------------------------|------|------------------------------------------------|
| GET    | `/me/reservas`         | Sesion  | Reservas de tours del cliente autenticado      |
| GET    | `/me/ventas`           | Sesion  | Historial de compras del cliente autenticado   |
| GET    | `/me/tipo-empleado`    | Sesion  | Tipo de empleado del usuario autenticado       |

---

### Reportes (requieren sesion de empleado)

| Método | Ruta                                  | Auth     | Descripción                                    |
|--------|---------------------------------------|----------|------------------------------------------------|
| GET    | `/reportes/pinturas-completo`         | Empleado | Pinturas con artista, colección y técnicas     |
| GET    | `/reportes/ventas-detalle`            | Empleado | Ventas con cliente, empleado y cantidad items  |
| GET    | `/reportes/artistas-resumen`          | Empleado | Artistas con totales de obra y valor           |
| GET    | `/reportes/artistas-con-ventas`       | Empleado | Artistas que tienen al menos una venta         |
| GET    | `/reportes/clientes-vip-compradores`  | Empleado | Clientes VIP con historial de compras          |
| GET    | `/reportes/ventas-por-mes`            | Empleado | Ingresos mensuales (todos los años)            |
| GET    | `/reportes/ventas-por-mes/{anio}`     | Empleado | Ingresos mensuales de un año específico        |
| GET    | `/reportes/tecnicas-populares`        | Empleado | Técnicas con más de 1 pintura                  |
| GET    | `/reportes/top-artistas-ventas`       | Empleado | Ranking de artistas por ingresos (CTE + RANK)  |
| GET    | `/reportes/colecciones-valor`         | Empleado | Ranking de colecciones por valor total         |

---

### Exportación CSV (descarga directa)

| Método | Ruta                      | Auth | Descripción                        |
|--------|---------------------------|------|------------------------------------|
| GET    | `/exportar/ventas-csv`    | No   | Descarga CSV de ventas detalladas  |
| GET    | `/exportar/pinturas-csv`  | No   | Descarga CSV del catálogo completo |
| GET    | `/exportar/artistas-csv`  | No   | Descarga CSV de artistas           |

---

## Estructura del proyecto

```
magic-bag-gallery-api/
├── backend/
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── artistasHandler.go
│   │   │   ├── authHandler.go
│   │   │   ├── coleccionHandler.go
│   │   │   ├── db.go
│   │   │   ├── envioHandler.go
│   │   │   ├── meHandler.go          # /me/reservas, /me/ventas, /me/tipo-empleado
│   │   │   ├── pinturasHandler.go
│   │   │   ├── reportesHandler.go
│   │   │   ├── tecnicaHandler.go
│   │   │   ├── toursHandler.go
│   │   │   ├── usuarioHandler.go
│   │   │   └── ventaHandler.go
│   │   ├── middleware/
│   │   │   └── auth.go               # SessionMiddleware, RequireRole
│   │   └── models/                   # Structs de datos
│   ├── main.go
│   ├── go.mod / go.sum
│   └── Dockerfile
├── db/
│   ├── ddl_magic_bag_gallery.sql     # Esquema, índices, vistas
│   └── dml_datos_iniciales_magic_bag_gallery.sql  # 45 usuarios de prueba
├── frontend/
│   ├── src/
│   │   ├── api/api.js                # Todas las funciones fetch
│   │   ├── components/
│   │   │   ├── Loader/
│   │   │   ├── Modal/
│   │   │   ├── Navbar/               # Responsive con hamburguesa
│   │   │   ├── PaintingCard/
│   │   │   └── ProtectedRoute/       # requireEmpleado + requireCliente
│   │   ├── context/
│   │   │   └── AuthContext.jsx       # sesion + tipoEmpleado + nombre
│   │   ├── hooks/
│   │   │   └── useCatalogFilters.js  # useReducer + useMemo
│   │   ├── pages/
│   │   │   ├── Admin/                # CRUD con tabs por tipo de empleado
│   │   │   ├── Artists/
│   │   │   ├── Catalog/              # useReducer + useMemo
│   │   │   ├── Collection/
│   │   │   ├── Login/                # Validación real del cliente
│   │   │   ├── MiCuenta/             # Solo clientes: reservas + historial
│   │   │   ├── NotFound/
│   │   │   ├── Register/
│   │   │   ├── Reports/              # 4 gráficas + exportación CSV
│   │   │   └── Tours/
│   │   ├── styles/
│   │   │   └── global.css
│   │   └── test/
│   │       ├── catalog.test.jsx      # Tests: reducer, AuthContext, validación
│   │       └── setup.js
│   ├── .eslintrc.json
│   ├── vitest.config.js
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Diseño de base de datos

### Vistas SQL

| Vista                     | Descripción                                           |
|---------------------------|-------------------------------------------------------|
| `vista_pinturas_completa` | Pintura + artista + colección + técnicas (STRING_AGG) |
| `vista_ventas_detalle`    | Venta + cliente + empleado + conteo de ítems          |
| `vista_artistas_resumen`  | Artista + reclutador + totales de obras y valor       |

### Índices

| Índice                  | Tabla   | Columna            | Justificación                   |
|-------------------------|---------|--------------------|---------------------------------|
| `idx_usuario_correo`    | usuario | correo_electronico | Búsqueda en login por correo    |
| `idx_pintura_artista`   | pintura | id_artista         | Filtrar obras por artista       |
| `idx_venta_cliente`     | venta   | id_cliente         | Historial de ventas por cliente |
| `idx_envio_venta`       | envio   | id_venta           | Consulta de envíos por venta    |
| `idx_pintura_coleccion` | pintura | id_coleccion       | Filtrar pinturas por colección  |

---

## Características técnicas destacadas

### Backend
- **SQL explícito** — sin ORM; queries escritas a mano con `database/sql`
- **Transacciones** — `BEGIN / COMMIT / ROLLBACK` en operaciones críticas
- **CTEs con RANK()** — rankings de artistas y colecciones
- **Subqueries** — `EXISTS` e `IN` en reportes
- **GROUP BY + HAVING** — reportes mensuales y técnicas populares
- **Autenticacion con sesion** — roles `cliente` / `empleado` con middleware
- **CORS configurado** — permite peticiones desde el frontend

### Frontend
- **useReducer** — estado del catálogo (filtros + búsqueda + datos)
- **useMemo** — lista filtrada memoizada para evitar recálculos
- **useCallback** — handlers estabilizados en Catalog, Navbar y MiCuenta
- **React Context** — AuthContext con sesion, rol, nombre y tipoEmpleado
- **React Router v6** — 9 rutas con ProtectedRoute por rol
- **Formularios controlados** — validación real en Login (email + min 6 chars)
- **ESLint** — configurado y sin errores (`npm run lint`)
- **Vitest** — tests unitarios del reducer, AuthContext y validación
- **Diseño responsivo** — Navbar con menú hamburguesa en móvil
- **Exportación CSV** — descarga directa de ventas, pinturas y artistas
