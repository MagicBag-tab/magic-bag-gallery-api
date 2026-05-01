-- ============================================
-- DML - MAGIC BAG GALLERY
-- Datos de prueba
-- ============================================

-- ============================================
-- USUARIOS (base para clientes y empleados)
-- ============================================
INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena) VALUES
('María',     'Pérez López',      'maria.perez@gmail.com',      '50255512345', '$2b$12$abc123hashcliente1'),
('Carlos',    'Méndez García',    'carlos.mendez@gmail.com',    '50255523456', '$2b$12$abc123hashcliente2'),
('Lucía',     'Hernández Ruiz',   'lucia.hernandez@gmail.com',  '50255534567', '$2b$12$abc123hashcliente3'),
('Andrés',    'González Torres',  'andres.gonzalez@gmail.com',  '50255545678', '$2b$12$abc123hashcliente4'),
('Sofía',     'Ramírez Castro',   'sofia.ramirez@gmail.com',    '50255556789', '$2b$12$abc123hashcliente5'),
('Diego',     'López Fuentes',    'diego.lopez@gmail.com',      '50255567890', '$2b$12$abc123hashcliente6'),
('Valentina', 'Morales Cifuentes','valentina.morales@gmail.com','50255578901', '$2b$12$abc123hashcliente7'),
('Fernando',  'Castillo Reyes',   'fernando.castillo@gmail.com','50255589012', '$2b$12$abc123hashcliente8'),
('Isabella',  'Flores Alvarado',  'isabella.flores@gmail.com',  '50255590123', '$2b$12$abc123hashcliente9'),
('Sebastián', 'Vásquez Molina',   'sebastian.vasquez@gmail.com','50255501234', '$2b$12$abc123hashcliente10'),
-- Empleados
('Ana',       'Solís Gramajo',    'ana.solis@magicbag.gt',      '50244411111', '$2b$12$abc123hashempleado1'),
('Roberto',   'Lima Barrios',     'roberto.lima@magicbag.gt',   '50244422222', '$2b$12$abc123hashempleado2'),
('Patricia',  'Aguilar Choc',     'patricia.aguilar@magicbag.gt','50244433333','$2b$12$abc123hashempleado3'),
('Miguel',    'Ramos Tzul',       'miguel.ramos@magicbag.gt',   '50244444444', '$2b$12$abc123hashempleado4'),
('Carmen',    'Cifuentes Pop',    'carmen.cifuentes@magicbag.gt','50244455555','$2b$12$abc123hashempleado5'),
-- Clientes extra
('Jorge',     'Bolaños Méndez',   'jorge.bolanos@gmail.com',    '50255511111', '$2b$12$abc123hashcliente11'),
('Daniela',   'Monterroso Pac',   'daniela.monterroso@gmail.com','50255522222','$2b$12$abc123hashcliente12'),
('Ricardo',   'Orellana Sajché',  'ricardo.orellana@gmail.com', '50255533333', '$2b$12$abc123hashcliente13'),
('Gabriela',  'Ajú Toj',          'gabriela.aju@gmail.com',     '50255544444', '$2b$12$abc123hashcliente14'),
('Pablo',     'Xicay Cúmez',      'pablo.xicay@gmail.com',      '50255555555', '$2b$12$abc123hashcliente15');

-- ============================================
-- CLIENTES (id_usuario 1–10, 16–20)
-- ============================================
INSERT INTO cliente (id_usuario, tipo_cliente) VALUES
(1,  'vip'),
(2,  'regular'),
(3,  'vip'),
(4,  'regular'),
(5,  'vip'),
(6,  'regular'),
(7,  'vip'),
(8,  'regular'),
(9,  'vip'),
(10, 'regular'),
(16, 'regular'),
(17, 'vip'),
(18, 'regular'),
(19, 'regular'),
(20, 'vip');

-- ============================================
-- EMPLEADOS (id_usuario 11–15)
-- ============================================
INSERT INTO empleado (id_usuario, tipo_empleado) VALUES
(11, 'guia'),
(12, 'asesor'),
(13, 'reclutador'),
(14, 'guia'),
(15, 'asesor');

-- ============================================
-- DIRECCIONES
-- ============================================
INSERT INTO direccion (id_cliente, detalle, nombre, ciudad, estado, codigo_postal, pais) VALUES
(1,  'Zona 10, Calle Reforma 12-34', 'Casa principal',      'Ciudad de Guatemala', 'Guatemala',     '01010', 'Guatemala'),
(1,  'Av. Brickell 1234, Apt 5B',    'Residencia Miami',    'Miami',               'Florida',        '33131', 'Estados Unidos'),
(2,  'Zona 14, Av. Las Américas 5-6','Casa',                'Ciudad de Guatemala', 'Guatemala',     '01014', 'Guatemala'),
(3,  'Zona 15, Vista Hermosa 3-45',  'Casa principal',      'Ciudad de Guatemala', 'Guatemala',     '01015', 'Guatemala'),
(3,  'Calle Serrano 45, Piso 3',     'Apartamento Madrid',  'Madrid',              'Comunidad de Madrid','28001','España'),
(4,  'Zona 1, 6ta Av 7-89',         'Casa',                'Ciudad de Guatemala', 'Guatemala',     '01001', 'Guatemala'),
(5,  'Zona 16, Cayalá Torre A',     'Apartamento',         'Ciudad de Guatemala', 'Guatemala',     '01016', 'Guatemala'),
(6,  'Antigua Guatemala, 3ra Calle', 'Casa colonial',       'Antigua Guatemala',   'Sacatepéquez',  '03001', 'Guatemala'),
(7,  'Zona 10, Blvd Los Próceres',  'Casa',                'Ciudad de Guatemala', 'Guatemala',     '01010', 'Guatemala'),
(7,  'Via Condotti 12, Int 4',      'Apartamento Roma',    'Roma',                'Lazio',          '00187', 'Italia'),
(8,  'Zona 11, Mariscal 23-45',     'Casa',                'Ciudad de Guatemala', 'Guatemala',     '01011', 'Guatemala'),
(9,  'Zona 10, Oakland Mall área',  'Penthouse',           'Ciudad de Guatemala', 'Guatemala',     '01010', 'Guatemala'),
(10, 'Xela, Zona 1, 5ta Calle',     'Casa',                'Quetzaltenango',      'Quetzaltenango','09001', 'Guatemala'),
(11, 'Zona 12, Colonia La Florida', 'Casa',                'Ciudad de Guatemala', 'Guatemala',     '01012', 'Guatemala'),
(12, 'Zona 7, Tikal 2 Sector B',    'Casa',                'Ciudad de Guatemala', 'Guatemala',     '01007', 'Guatemala');

-- ============================================
-- ARTISTAS
-- ============================================
INSERT INTO artista (nombre, apellido, nacionalidad, id_reclutador) VALUES
('Diego',    'Rivera Barrientos',  'Mexicana',     3),
('Frida',    'Xicay Ajú',          'Guatemalteca', 3),
('Remedios', 'Varo Cifuentes',     'Española',     3),
('Joaquín',  'Torres García',      'Uruguaya',     3),
('Rufino',   'Tamayo Pop',         'Mexicana',     3),
('Carlos',   'Mérida Ajquijay',    'Guatemalteca', 3),
('Rosa',     'Cabrera Tzul',       'Guatemalteca', 3),
('Luis',     'González Morales',   'Guatemalteca', 3),
('Elena',    'Quispe Mamani',      'Peruana',      3),
('Mateo',    'Rodríguez Sajché',   'Colombiana',   3);

-- ============================================
-- COLECCIONES
-- ============================================
INSERT INTO coleccion (nombre, descripcion, exclusiva, fecha_lanzamiento) VALUES
('Raíces Mayas',         'Colección que celebra la herencia maya guatemalteca con técnicas contemporáneas.',     TRUE,  '2024-01-15'),
('Paisajes del Altiplano','Obras que capturan la majestuosidad de los paisajes guatemaltecos.',                   FALSE, '2024-03-01'),
('Arte Abstracto GT',    'Exploración del arte abstracto por artistas guatemaltecos emergentes.',                 FALSE, '2024-05-10'),
('Mujeres de Colores',   'Colección exclusiva dedicada a retratos de mujeres indígenas guatemaltecas.',          TRUE,  '2024-07-20'),
('Naturaleza Viva',      'Pinturas de flora y fauna de Guatemala en técnicas mixtas.',                            FALSE, '2024-09-05');

-- ============================================
-- TÉCNICAS
-- ============================================
INSERT INTO tecnica (nombre, descripcion) VALUES
('Óleo sobre lienzo',    'Técnica tradicional usando pigmentos mezclados con aceite sobre lienzo de algodón.'),
('Acuarela',             'Pintura transparente a base de agua, ideal para paisajes y naturaleza.'),
('Acrílico',             'Pintura de secado rápido a base de polímeros, versátil y duradera.'),
('Técnica mixta',        'Combinación de múltiples medios y materiales en una sola obra.'),
('Encáustica',           'Técnica antigua que usa cera de abeja mezclada con pigmentos de colores.'),
('Gouache',              'Similar a la acuarela pero opaca, con mayor cuerpo y luminosidad.'),
('Pastel',               'Barras de pigmento puro comprimido, permite texturas suaves y difuminados.');

-- ============================================
-- PINTURAS
-- ============================================
INSERT INTO pintura (id_artista, titulo, descripcion, precio, fecha_creacion, imagen_path, imagen_tipo, imagen_nombre, exclusiva, id_coleccion) VALUES
(1, 'Quetzal Dorado',         'Representación del quetzal nacional en tonos dorados y verdes.',         15000.00, '2023-06-01', '/uploads/pinturas/quetzal_dorado.jpg',        'image/jpeg', 'quetzal_dorado.jpg',        TRUE,  1),
(2, 'Mercado de Chichicastenango','Escena vibrante del mercado indígena más famoso de Guatemala.',      8500.00,  '2023-08-15', '/uploads/pinturas/mercado_chichi.jpg',         'image/jpeg', 'mercado_chichi.jpg',         FALSE, 2),
(3, 'Volcán de Agua',         'Vista majestuosa del Volcán de Agua desde Antigua Guatemala.',           12000.00, '2023-09-20', '/uploads/pinturas/volcan_agua.jpg',           'image/jpeg', 'volcan_agua.jpg',           FALSE, 2),
(4, 'Tejedoras de Sololá',    'Mujeres indígenas tejiendo trajes típicos en el altiplano.',             18000.00, '2023-10-05', '/uploads/pinturas/tejedoras_solola.jpg',       'image/jpeg', 'tejedoras_solola.jpg',       TRUE,  4),
(5, 'Lago Atitlán al Amanecer','Amanecer sobre el Lago Atitlán con los tres volcanes de fondo.',        9500.00,  '2023-11-12', '/uploads/pinturas/atitlan_amanecer.jpg',       'image/jpeg', 'atitlan_amanecer.jpg',       FALSE, 2),
(6, 'Geometría Maya',         'Patrones geométricos mayas reinterpretados en arte contemporáneo.',      22000.00, '2023-12-01', '/uploads/pinturas/geometria_maya.jpg',        'image/jpeg', 'geometria_maya.jpg',        TRUE,  1),
(7, 'Flores de Guatemala',    'Composición floral con orquídeas y flores tropicales guatemaltecas.',    6000.00,  '2024-01-10', '/uploads/pinturas/flores_guatemala.jpg',       'image/jpeg', 'flores_guatemala.jpg',       FALSE, 5),
(8, 'Abstracción Urbana',     'Interpretación abstracta de la ciudad de Guatemala moderna.',            11000.00, '2024-02-20', '/uploads/pinturas/abstraccion_urbana.jpg',     'image/jpeg', 'abstraccion_urbana.jpg',     FALSE, 3),
(9, 'Danza del Venado',       'Representación de la danza folklórica guatemalteca del venado.',         14000.00, '2024-03-15', '/uploads/pinturas/danza_venado.jpg',          'image/jpeg', 'danza_venado.jpg',          TRUE,  1),
(10,'Semana Santa Antigua',   'Procesión de Semana Santa en las calles de Antigua Guatemala.',          7500.00,  '2024-04-01', '/uploads/pinturas/semana_santa.jpg',          'image/jpeg', 'semana_santa.jpg',          FALSE, NULL),
(1, 'Cacao Sagrado',          'El cacao como elemento sagrado en la cultura maya prehispánica.',        19000.00, '2024-04-20', '/uploads/pinturas/cacao_sagrado.jpg',         'image/jpeg', 'cacao_sagrado.jpg',         TRUE,  1),
(2, 'Niña de Todos Santos',   'Retrato de niña con traje típico de Todos Santos Cuchumatán.',           13000.00, '2024-05-05', '/uploads/pinturas/nina_todos_santos.jpg',      'image/jpeg', 'nina_todos_santos.jpg',      TRUE,  4),
(3, 'Ceiba Sagrada',          'La ceiba, árbol nacional de Guatemala, en técnica mixta.',               8000.00,  '2024-05-25', '/uploads/pinturas/ceiba_sagrada.jpg',         'image/jpeg', 'ceiba_sagrada.jpg',         FALSE, 5),
(4, 'Huipil de Nebaj',        'Detalle de los intrincados bordados del huipil de Nebaj.',               16000.00, '2024-06-10', '/uploads/pinturas/huipil_nebaj.jpg',          'image/jpeg', 'huipil_nebaj.jpg',          TRUE,  4),
(5, 'Río Motagua',            'Paisaje del río Motagua con vegetación tropical.',                       5500.00,  '2024-06-30', '/uploads/pinturas/rio_motagua.jpg',           'image/jpeg', 'rio_motagua.jpg',           FALSE, 5),
(6, 'Códice Contemporáneo',   'Reinterpretación de los códices mayas en formato contemporáneo.',        25000.00, '2024-07-15', '/uploads/pinturas/codice_contemporaneo.jpg',  'image/jpeg', 'codice_contemporaneo.jpg',  TRUE,  1),
(7, 'Orquídea Morada',        'Monja blanca, flor nacional de Guatemala, en acuarela.',                 4500.00,  '2024-07-30', '/uploads/pinturas/orquidea_morada.jpg',       'image/jpeg', 'orquidea_morada.jpg',       FALSE, 5),
(8, 'Ciudad Fragmentada',     'Visión abstracta de la fragmentación urbana en Guatemala.',              10000.00, '2024-08-10', '/uploads/pinturas/ciudad_fragmentada.jpg',    'image/jpeg', 'ciudad_fragmentada.jpg',    FALSE, 3),
(9, 'Marimba',                'La marimba, instrumento nacional, representada en colores vibrantes.',   9000.00,  '2024-08-25', '/uploads/pinturas/marimba.jpg',               'image/jpeg', 'marimba.jpg',               FALSE, NULL),
(10,'Xelajú Eterno',          'Panorámica de Quetzaltenango desde el cerro El Baúl.',                  7000.00,  '2024-09-05', '/uploads/pinturas/xelaju_eterno.jpg',         'image/jpeg', 'xelaju_eterno.jpg',         FALSE, 2),
(1, 'Popol Vuh',              'Escenas del libro sagrado maya quiché en técnica encáustica.',           30000.00, '2024-09-20', '/uploads/pinturas/popol_vuh.jpg',             'image/jpeg', 'popol_vuh.jpg',             TRUE,  1),
(2, 'Abuela Tejedora',        'Retrato íntimo de anciana indígena tejiendo en telar de cintura.',       17000.00, '2024-10-01', '/uploads/pinturas/abuela_tejedora.jpg',       'image/jpeg', 'abuela_tejedora.jpg',       TRUE,  4),
(3, 'Selva Petén',            'Densa selva del Petén con fauna autóctona guatemalteca.',                8500.00,  '2024-10-15', '/uploads/pinturas/selva_peten.jpg',           'image/jpeg', 'selva_peten.jpg',           FALSE, 5),
(4, 'Espíritu del Maíz',      'El maíz como elemento central de la cosmovisión maya.',                  20000.00, '2024-10-30', '/uploads/pinturas/espiritu_maiz.jpg',         'image/jpeg', 'espiritu_maiz.jpg',         TRUE,  1),
(5, 'Tikal al Atardecer',     'Templo I de Tikal emergiendo de la selva en el atardecer.',              11500.00, '2024-11-10', '/uploads/pinturas/tikal_atardecer.jpg',       'image/jpeg', 'tikal_atardecer.jpg',       FALSE, 2);

-- ============================================
-- PINTURA_TECNICA
-- ============================================
INSERT INTO pintura_tecnica (id_pintura, id_tecnica) VALUES
(1,  1), (1,  3),
(2,  1),
(3,  2),
(4,  1), (4,  4),
(5,  2),
(6,  3), (6,  4),
(7,  2),
(8,  3), (8,  4),
(9,  1),
(10, 1),
(11, 5),
(12, 1),
(13, 4),
(14, 1), (14, 6),
(15, 2),
(16, 4), (16, 5),
(17, 2),
(18, 3),
(19, 3),
(20, 1),
(21, 5),
(22, 1),
(23, 2), (23, 4),
(24, 1), (24, 3),
(25, 1);

-- ============================================
-- TOURS
-- ============================================
INSERT INTO tour (id_guia, nombre, descripcion, fecha_inicio, fecha_fin, horario, precio) VALUES
(1, 'Arte Maya Contemporáneo',  'Recorrido por las obras de arte inspiradas en la cultura maya.',         '2025-01-10', '2025-01-10', '10:00 - 12:00', 150.00),
(4, 'Técnicas Tradicionales',   'Visita guiada enfocada en técnicas de pintura tradicionales.',           '2025-01-17', '2025-01-17', '14:00 - 16:00', 120.00),
(1, 'Colecciones Exclusivas',   'Tour privado por las colecciones exclusivas de la galería.',             '2025-02-07', '2025-02-07', '09:00 - 11:00', 300.00),
(4, 'Paisajes Guatemaltecos',   'Recorrido temático por pinturas de paisajes nacionales.',                '2025-02-14', '2025-02-14', '15:00 - 17:00', 100.00),
(1, 'Artistas Emergentes GT',   'Presentación de los artistas guatemaltecos más prometedores.',           '2025-03-07', '2025-03-07', '11:00 - 13:00', 130.00),
(4, 'Noche de Arte',            'Tour nocturno especial con cóctel de bienvenida incluido.',              '2025-03-21', '2025-03-21', '19:00 - 21:00', 250.00),
(1, 'Arte y Naturaleza',        'Conexión entre las pinturas de naturaleza y el entorno guatemalteco.',   '2025-04-04', '2025-04-04', '10:00 - 12:00', 110.00),
(4, 'Mujeres Artistas',         'Tour dedicado a las obras de artistas femeninas de la galería.',         '2025-04-18', '2025-04-18', '14:00 - 16:00', 140.00),
(1, 'Historia del Arte Maya',   'Conferencia y recorrido sobre la historia del arte maya guatemalteco.',  '2025-05-02', '2025-05-02', '09:00 - 12:00', 180.00),
(4, 'Arte Abstracto',           'Introducción al arte abstracto guatemalteco contemporáneo.',             '2025-05-16', '2025-05-16', '15:00 - 17:00', 120.00);

-- ============================================
-- CLIENTE_TOUR
-- ============================================
INSERT INTO cliente_tour (id_cliente, id_tour, fecha_reserva) VALUES
(1,  1, '2025-01-05'),
(2,  1, '2025-01-06'),
(3,  1, '2025-01-07'),
(4,  2, '2025-01-10'),
(5,  2, '2025-01-11'),
(6,  3, '2025-01-20'),
(7,  3, '2025-01-21'),
(1,  3, '2025-01-22'),
(8,  4, '2025-02-01'),
(9,  4, '2025-02-02'),
(10, 5, '2025-02-20'),
(11, 5, '2025-02-21'),
(12, 6, '2025-03-01'),
(13, 6, '2025-03-02'),
(3,  6, '2025-03-03'),
(5,  7, '2025-03-15'),
(7,  7, '2025-03-16'),
(14, 8, '2025-04-01'),
(15, 8, '2025-04-02'),
(2,  9, '2025-04-15'),
(4,  9, '2025-04-16'),
(6,  9, '2025-04-17'),
(8,  10,'2025-05-01'),
(9,  10,'2025-05-02'),
(1,  10,'2025-05-03');

-- ============================================
-- VENTAS
-- ============================================
INSERT INTO venta (id_cliente, id_empleado, fecha_venta, precio) VALUES
(1,  2, '2025-01-15', 15000.00),
(3,  5, '2025-01-22', 18000.00),
(5,  2, '2025-02-03', 22000.00),
(7,  5, '2025-02-14', 8500.00),
(9,  2, '2025-02-28', 12000.00),
(1,  5, '2025-03-05', 19000.00),
(12, 2, '2025-03-12', 9500.00),
(3,  5, '2025-03-20', 13000.00),
(5,  2, '2025-04-02', 25000.00),
(7,  5, '2025-04-10', 14000.00),
(15, 2, '2025-04-18', 6000.00),
(9,  5, '2025-04-25', 30000.00),
(1,  2, '2025-05-01', 17000.00),
(3,  5, '2025-05-08', 8000.00),
(5,  2, '2025-05-15', 20000.00),
(7,  5, '2025-05-22', 11500.00),
(9,  2, '2025-06-01', 7500.00),
(12, 5, '2025-06-08', 16000.00),
(15, 2, '2025-06-15', 9000.00),
(1,  5, '2025-06-22', 5500.00),
(3,  2, '2025-07-01', 11000.00),
(5,  5, '2025-07-08', 7000.00),
(7,  2, '2025-07-15', 4500.00),
(9,  5, '2025-07-22', 10000.00),
(12, 2, '2025-07-29', 8500.00);

-- ============================================
-- DETALLE_VENTA
-- ============================================
INSERT INTO detalle_venta (id_venta, id_pintura, cantidad, precio_unitario) VALUES
(1,  1,  1, 15000.00),
(2,  4,  1, 18000.00),
(3,  6,  1, 22000.00),
(4,  2,  1, 8500.00),
(5,  3,  1, 12000.00),
(6,  11, 1, 19000.00),
(7,  5,  1, 9500.00),
(8,  12, 1, 13000.00),
(9,  16, 1, 25000.00),
(10, 9,  1, 14000.00),
(11, 7,  1, 6000.00),
(12, 21, 1, 30000.00),
(13, 22, 1, 17000.00),
(14, 13, 1, 8000.00),
(15, 24, 1, 20000.00),
(16, 25, 1, 11500.00),
(17, 10, 1, 7500.00),
(18, 14, 1, 16000.00),
(19, 19, 1, 9000.00),
(20, 15, 1, 5500.00),
(21, 8,  1, 11000.00),
(22, 20, 1, 7000.00),
(23, 17, 1, 4500.00),
(24, 18, 1, 10000.00),
(25, 23, 1, 8500.00);

-- ============================================
-- ENVÍOS
-- ============================================
INSERT INTO envio (id_venta, direccion_envio, fecha_envio, estado_envio) VALUES
(1,  'Zona 10, Calle Reforma 12-34, Ciudad de Guatemala',          '2025-01-18', 'entregado'),
(2,  'Zona 15, Vista Hermosa 3-45, Ciudad de Guatemala',           '2025-01-25', 'entregado'),
(3,  'Zona 16, Cayalá Torre A, Ciudad de Guatemala',               '2025-02-06', 'entregado'),
(4,  'Zona 10, Blvd Los Próceres, Ciudad de Guatemala',            '2025-02-17', 'entregado'),
(5,  'Zona 10, Oakland Mall área, Ciudad de Guatemala',            '2025-03-03', 'entregado'),
(6,  'Av. Brickell 1234, Apt 5B, Miami, Florida, USA',             '2025-03-08', 'entregado'),
(7,  'Zona 12, Colonia La Florida, Ciudad de Guatemala',           '2025-03-15', 'entregado'),
(8,  'Calle Serrano 45, Piso 3, Madrid, España',                   '2025-03-25', 'en tránsito'),
(9,  'Zona 16, Cayalá Torre A, Ciudad de Guatemala',               '2025-04-05', 'entregado'),
(10, 'Via Condotti 12, Int 4, Roma, Italia',                       '2025-04-15', 'entregado'),
(12, 'Zona 10, Oakland Mall área, Ciudad de Guatemala',            '2025-04-28', 'entregado'),
(13, 'Zona 10, Calle Reforma 12-34, Ciudad de Guatemala',          '2025-05-04', 'entregado'),
(15, 'Zona 16, Cayalá Torre A, Ciudad de Guatemala',               '2025-05-18', 'en tránsito'),
(16, 'Zona 10, Blvd Los Próceres, Ciudad de Guatemala',            '2025-05-25', 'entregado'),
(18, 'Zona 7, Tikal 2 Sector B, Ciudad de Guatemala',              '2025-06-11', 'entregado');
