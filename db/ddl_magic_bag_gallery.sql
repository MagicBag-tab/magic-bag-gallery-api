-- DDL MAGIC BAG GALLERY --
/*
 * Autor: Sarah Rachel Estrada Bonilla
 * Fecha: 2024-06-17
 * Descripción: Script SQL para crear la estrucutura de la galería de arte
*/


-- Usuario --
CREATE TABLE IF NOT EXISTS usuario (
    id_usuario SERIAL UNIQUE NOT NULL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    correo_electronico VARCHAR(100) UNIQUE NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    contrasena VARCHAR(255) NOT NULL
);

-- Cliente --
CREATE TABLE IF NOT EXISTS cliente (
    id_cliente SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_usuario INT NOT NULL,
    tipo_cliente VARCHAR(50) NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario)
);

-- Empleado --
CREATE TABLE IF NOT EXISTS empleado (
    id_empleado SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_usuario INT NOT NULL,
    tipo_empleado VARCHAR(50) NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario)
);

-- Dirección --
CREATE TABLE IF NOT EXISTS direccion (
    id_direccion SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_cliente INT NOT NULL,
    detalle VARCHAR(100) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    ciudad VARCHAR(50) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    codigo_postal VARCHAR(20) NOT NULL,
    pais VARCHAR(50) NOT NULL,
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente)
);

-- Artista --
CREATE TABLE IF NOT EXISTS artista (
    id_artista SERIAL UNIQUE NOT NULL PRIMARY KEY,
    nombre_completo VARCHAR(100) NOT NULL,
    nacionalidad VARCHAR(50) NOT NULL,
    id_reclutador INT NOT NULL,
    FOREIGN KEY (id_reclutador) REFERENCES empleado(id_empleado)
);

-- Colección --
CREATE TABLE IF NOT EXISTS coleccion (
    id_coleccion SERIAL UNIQUE NOT NULL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    exclusiva BOOLEAN NOT NULL,
    fecha_lanzamiento DATE NOT NULL
);

-- Pintura --
CREATE TABLE IF NOT EXISTS pintura (
    id_pintura SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_artista INT NOT NULL,
    titulo VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    precio DECIMAL(10, 2) NOT NULL,
    fecha_creacion DATE NOT NULL,
    imagen_path VARCHAR(255) NOT NULL,
    imagen_tipo VARCHAR(50) NOT NULL,
    imagen_nombre VARCHAR(255) NOT NULL,
    exclusiva BOOLEAN NOT NULL,
    id_coleccion INT,
    FOREIGN KEY (id_artista) REFERENCES artista(id_artista),
    FOREIGN KEY (id_coleccion) REFERENCES coleccion(id_coleccion)
);

-- Técnica --
CREATE TABLE IF NOT EXISTS tecnica (
    id_tecnica SERIAL UNIQUE NOT NULL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL
);

-- Pintura_Tecnica --
CREATE TABLE IF NOT EXISTS pintura_tecnica (
    id_pintura_tecnica SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_pintura INT NOT NULL,
    id_tecnica INT NOT NULL,
    FOREIGN KEY (id_pintura) REFERENCES pintura(id_pintura),
    FOREIGN KEY (id_tecnica) REFERENCES tecnica(id_tecnica)
);

-- Tour -- 
CREATE TABLE IF NOT EXISTS tour (
    id_tour SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_guia INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    horario VARCHAR(50) NOT NULL,
    precio DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (id_guia) REFERENCES empleado(id_empleado)
);

-- Cliente_Tour --
CREATE TABLE IF NOT EXISTS cliente_tour (
    id_cliente_tour SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_cliente INT NOT NULL,
    id_tour INT NOT NULL,
    fecha_reserva DATE NOT NULL,
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente),
    FOREIGN KEY (id_tour) REFERENCES tour(id_tour)
);

-- Venta --
CREATE TABLE IF NOT EXISTS venta (
    id_venta SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_cliente INT NOT NULL,
    id_empleado INT NOT NULL,
    fecha_venta DATE NOT NULL,
    precio DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente),
    FOREIGN KEY (id_empleado) REFERENCES empleado(id_empleado)
);

-- Detalle_Venta --
CREATE TABLE IF NOT EXISTS detalle_venta (
    id_detalle_venta SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_venta INT NOT NULL,
    id_pintura INT NOT NULL,
    cantidad INT NOT NULL,
    precio_unitario DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (id_venta) REFERENCES venta(id_venta),
    FOREIGN KEY (id_pintura) REFERENCES pintura(id_pintura)
);

-- Envío --
CREATE TABLE IF NOT EXISTS envio (
    id_envio SERIAL UNIQUE NOT NULL PRIMARY KEY,
    id_venta INT NOT NULL,
    direccion_envio VARCHAR(255) NOT NULL,
    fecha_envio DATE NOT NULL,
    estado_envio VARCHAR(50) NOT NULL,
    FOREIGN KEY (id_venta) REFERENCES venta(id_venta)
);

-- Indices --

-- Búsqueda por correo
CREATE INDEX idx_usuario_correo
    ON usuario(correo_electronico);

-- Filtrar pinturas por artista
CREATE INDEX idx_pintura_artista
    ON pintura(id_artista);

-- Historial de ventas por cliente
CREATE INDEX idx_venta_cliente
    ON venta(id_cliente);

-- Consultas de envíos por venta
CREATE INDEX idx_envio_venta
    ON envio(id_venta);

-- Filtrar pinturas por colección
CREATE INDEX idx_pintura_coleccion
    ON pintura(id_coleccion);