# Magic Bag Gallery API

API REST para la galería de arte **Magic Bag Gallery**, desarrollada en Go con PostgreSQL.

## Requisitos

- [Docker](https://www.docker.com/) y Docker Compose instalados.

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

La API estará disponible en `http://localhost:8888`.  
La base de datos se inicializa automáticamente con el DDL y datos de prueba.

---

## Credenciales de base de datos

| Variable          | Valor               |
|-------------------|---------------------|
| `POSTGRES_USER`   | `proy2`             |
| `POSTGRES_PASSWORD` | `secret`          |
| `POSTGRES_DB`     | `magic_bag_gallery` |

---

## Endpoints principales

### Autenticación
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/login` | Login, devuelve JWT |
| POST | `/api/register/cliente` | Registro de cliente |

### Catálogo (público)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/pinturas` | Listar pinturas |
| GET | `/api/pinturas/{id}` | Detalle de pintura |
| GET | `/api/artistas` | Listar artistas |
| GET | `/api/colecciones` | Listar colecciones |
| GET | `/api/tecnicas` | Listar técnicas |
| GET | `/api/tours` | Listar tours |

### Reportes (con JOINs, subqueries, GROUP BY, CTEs y VIEWs)
| Método | Ruta | Tipo de consulta |
|--------|------|-----------------|
| GET | `/api/reportes/pinturas-completo` | VIEW + JOIN múltiple |
| GET | `/api/reportes/ventas-detalle` | VIEW + JOIN múltiple |
| GET | `/api/reportes/artistas-resumen` | VIEW + JOIN múltiple |
| GET | `/api/reportes/artistas-con-ventas` | Subquery EXISTS |
| GET | `/api/reportes/clientes-vip-compradores` | Subquery IN |
| GET | `/api/reportes/ventas-por-mes` | GROUP BY + HAVING + agregación |
| GET | `/api/reportes/ventas-por-mes/{anio}` | GROUP BY + HAVING filtrado |
| GET | `/api/reportes/tecnicas-populares` | GROUP BY + HAVING + agregación |
| GET | `/api/reportes/top-artistas-ventas` | CTE (WITH) + RANK() |
| GET | `/api/reportes/colecciones-valor` | CTE (WITH) + RANK() |

### Exportación CSV
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/exportar/ventas-csv` | Reporte de ventas en CSV |
| GET | `/api/exportar/pinturas-csv` | Catálogo de pinturas en CSV |
| GET | `/api/exportar/artistas-csv` | Resumen de artistas en CSV |

### CRUD (requiere JWT de empleado)
- `/api/pinturas`, `/api/artistas`, `/api/colecciones`, `/api/tecnicas`
- `/api/ventas`, `/api/detalles-venta`, `/api/envios`
- `/api/tours`, `/api/reservas`
- `/api/usuarios`

---

## Estructura del proyecto

```
magic-bag-gallery-api/
├── backend/
│   ├── internal/
│   │   ├── handlers/       # Handlers HTTP
│   │   ├── middleware/     # JWT y control de roles
│   │   └── models/         # Modelos de datos
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── db/
│   ├── ddl_magic_bag_gallery.sql       # Esquema + índices + vistas
│   └── dml_datos_iniciales_magic_bag_gallery.sql  # Datos de prueba
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Características técnicas

- **SQL explícito** — sin ORM, todas las queries escritas a mano
- **Transacciones** — BEGIN / COMMIT / ROLLBACK explícito en operaciones críticas
- **Vistas SQL** — `vista_pinturas_completa`, `vista_ventas_detalle`, `vista_artistas_resumen`
- **CTEs** — consultas con `WITH` para rankings y reportes complejos
- **Autenticación JWT** — login con roles `cliente` / `empleado`
- **Exportación CSV** — descarga directa desde la UI