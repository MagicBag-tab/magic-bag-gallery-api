-- ============================================
-- DML - MAGIC BAG GALLERY
-- Datos de prueba realistas con artistas modernos
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
('Ana',       'Solís Gramajo',    'ana.solis@magicbag.gt',      '50244411111', '$2b$12$abc123hashempleado1'),
('Roberto',   'Lima Barrios',     'roberto.lima@magicbag.gt',   '50244422222', '$2b$12$abc123hashempleado2'),
('Patricia',  'Aguilar Choc',     'patricia.aguilar@magicbag.gt','50244433333','$2b$12$abc123hashempleado3'),
('Miguel',    'Ramos Tzul',       'miguel.ramos@magicbag.gt',   '50244444444', '$2b$12$abc123hashempleado4'),
('Carmen',    'Cifuentes Pop',    'carmen.cifuentes@magicbag.gt','50244455555','$2b$12$abc123hashempleado5'),
('Jorge',     'Bolaños Méndez',   'jorge.bolanos@gmail.com',    '50255511111', '$2b$12$abc123hashcliente11'),
('Daniela',   'Monterroso Pac',   'daniela.monterroso@gmail.com','50255522222','$2b$12$abc123hashcliente12'),
('Ricardo',   'Orellana Sajché',  'ricardo.orellana@gmail.com', '50255533333', '$2b$12$abc123hashcliente13'),
('Gabriela',  'Ajú Toj',          'gabriela.aju@gmail.com',     '50255544444', '$2b$12$abc123hashcliente14'),
('Pablo',     'Xicay Cúmez',      'pablo.xicay@gmail.com',      '50255555555', '$2b$12$abc123hashcliente15');

-- ============================================
-- CLIENTES
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
-- EMPLEADOS
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
(1,  'Zona 10, Calle Reforma 12-34',  'Casa principal',     'Ciudad de Guatemala', 'Guatemala',          '01010', 'Guatemala'),
(1,  'Av. Brickell 1234, Apt 5B',     'Residencia Miami',   'Miami',               'Florida',            '33131', 'Estados Unidos'),
(2,  'Zona 14, Av. Las Américas 5-6', 'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01014', 'Guatemala'),
(3,  'Zona 15, Vista Hermosa 3-45',   'Casa principal',     'Ciudad de Guatemala', 'Guatemala',          '01015', 'Guatemala'),
(3,  'Calle Serrano 45, Piso 3',      'Apartamento Madrid', 'Madrid',              'Comunidad de Madrid','28001', 'España'),
(4,  'Zona 1, 6ta Av 7-89',          'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01001', 'Guatemala'),
(5,  'Zona 16, Cayalá Torre A',       'Apartamento',        'Ciudad de Guatemala', 'Guatemala',          '01016', 'Guatemala'),
(6,  'Antigua Guatemala, 3ra Calle',  'Casa colonial',      'Antigua Guatemala',   'Sacatepéquez',       '03001', 'Guatemala'),
(7,  'Zona 10, Blvd Los Próceres',   'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01010', 'Guatemala'),
(7,  'Via Condotti 12, Int 4',        'Apartamento Roma',   'Roma',                'Lazio',              '00187', 'Italia'),
(8,  'Zona 11, Mariscal 23-45',       'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01011', 'Guatemala'),
(9,  'Zona 10, Oakland Mall área',    'Penthouse',          'Ciudad de Guatemala', 'Guatemala',          '01010', 'Guatemala'),
(10, 'Xela, Zona 1, 5ta Calle',       'Casa',               'Quetzaltenango',      'Quetzaltenango',     '09001', 'Guatemala'),
(11, 'Zona 12, Colonia La Florida',   'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01012', 'Guatemala'),
(12, 'Zona 7, Tikal 2 Sector B',      'Casa',               'Ciudad de Guatemala', 'Guatemala',          '01007', 'Guatemala');

-- ============================================
-- ARTISTAS (nombre completo en un solo campo)
-- ============================================
INSERT INTO artista (nombre, nacionalidad, id_reclutador) VALUES
('Jean-Michel Basquiat', 'Estadounidense', 3),
('Banksy',               'Británica',      3),
('Yayoi Kusama',         'Japonesa',       3),
('Jeff Koons',           'Estadounidense', 3),
('Damien Hirst',         'Británica',      3),
('Takashi Murakami',     'Japonesa',       3),
('KAWS',                 'Estadounidense', 3),
('Shepard Fairey',       'Estadounidense', 3),
('Kerry James Marshall', 'Estadounidense', 3),
('Kehinde Wiley',        'Estadounidense', 3);

-- ============================================
-- COLECCIONES
-- ============================================
INSERT INTO coleccion (nombre, descripcion, exclusiva, fecha_lanzamiento) VALUES
('Neo-Expresionismo',  'Obras clave del movimiento neo-expresionista de finales del siglo XX.',           TRUE,  '2024-01-15'),
('Arte Callejero',     'Street art y arte urbano de los artistas más influyentes del mundo.',             FALSE, '2024-03-01'),
('Arte Contemporáneo', 'Selección de obras del arte contemporáneo global más cotizado.',                  FALSE, '2024-05-10'),
('Pop Art Moderno',    'Colección exclusiva que fusiona el pop art clásico con tendencias actuales.',     TRUE,  '2024-07-20'),
('Identidad y Cultura','Obras que exploran la identidad racial, cultural y social en el arte moderno.',   FALSE, '2024-09-05');

-- ============================================
-- TÉCNICAS
-- ============================================
INSERT INTO tecnica (nombre, descripcion) VALUES
('Óleo sobre lienzo', 'Técnica tradicional usando pigmentos mezclados con aceite sobre lienzo de algodón.'),
('Acuarela',          'Pintura transparente a base de agua, ideal para paisajes y naturaleza.'),
('Acrílico',          'Pintura de secado rápido a base de polímeros, versátil y duradera.'),
('Técnica mixta',     'Combinación de múltiples medios y materiales en una sola obra.'),
('Serigrafía',        'Técnica de impresión que usa una malla para transferir tinta sobre la superficie.'),
('Spray sobre muro',  'Pintura en aerosol aplicada sobre superficies urbanas, característica del street art.'),
('Instalación',       'Obra artística tridimensional que transforma un espacio completo.');

-- ============================================
-- PINTURAS
-- ============================================
INSERT INTO pintura (id_artista, titulo, descripcion, precio, fecha_creacion, imagen_path, imagen_tipo, imagen_nombre, exclusiva, id_coleccion) VALUES
(1, 'Untitled (Skull)',                       'Icónico cráneo neo-expresionista de Basquiat, símbolo de su obra.',                    95000.00, '1981-01-01', '/uploads/pinturas/basquiat_skull.jpg',           'image/jpeg', 'basquiat_skull.jpg',           TRUE,  1),
(1, 'Hollywood Africans',                     'Crítica social sobre la representación de los afroamericanos en Hollywood.',            78000.00, '1983-01-01', '/uploads/pinturas/basquiat_hollywood.jpg',       'image/jpeg', 'basquiat_hollywood.jpg',       TRUE,  1),
(1, 'Warrior',                                'Figura guerrera que mezcla iconografía africana con el grafiti urbano.',                65000.00, '1982-01-01', '/uploads/pinturas/basquiat_warrior.jpg',         'image/jpeg', 'basquiat_warrior.jpg',         FALSE, 1),
(2, 'Girl with Balloon',                      'La icónica niña con globo rojo, símbolo de esperanza e inocencia.',                    88000.00, '2002-01-01', '/uploads/pinturas/banksy_balloon_girl.jpg',       'image/jpeg', 'banksy_balloon_girl.jpg',       FALSE, 2),
(2, 'Flower Thrower',                         'Manifestante lanzando un ramo de flores en lugar de una bomba.',                       72000.00, '2003-01-01', '/uploads/pinturas/banksy_flower_thrower.jpg',     'image/jpeg', 'banksy_flower_thrower.jpg',     FALSE, 2),
(2, 'Napalm',                                 'Mickey Mouse y Ronald McDonald tomando de la mano a la niña de napalm.',               55000.00, '2004-01-01', '/uploads/pinturas/banksy_napalm.jpg',             'image/jpeg', 'banksy_napalm.jpg',             TRUE,  2),
(3, 'Infinity Nets',                          'Patrón infinito de redes que cubre todo el lienzo, obra obsesiva de Kusama.',          90000.00, '1958-01-01', '/uploads/pinturas/kusama_infinity_nets.jpg',      'image/jpeg', 'kusama_infinity_nets.jpg',      TRUE,  3),
(3, 'Pumpkin Yellow',                         'La icónica calabaza amarilla con puntos negros de Kusama.',                            68000.00, '1994-01-01', '/uploads/pinturas/kusama_pumpkin.jpg',            'image/jpeg', 'kusama_pumpkin.jpg',            FALSE, 3),
(3, 'Flowers that Bloom at Midnight',         'Flores oníricas en colores vibrantes características del estilo de Kusama.',           45000.00, '2012-01-01', '/uploads/pinturas/kusama_midnight_flowers.jpg',   'image/jpeg', 'kusama_midnight_flowers.jpg',   FALSE, 3),
(4, 'Balloon Dog Blue',                       'El famoso perro globo azul de Jeff Koons en su versión pictórica.',                   120000.00, '1994-01-01', '/uploads/pinturas/koons_balloon_dog.jpg',         'image/jpeg', 'koons_balloon_dog.jpg',         TRUE,  4),
(4, 'Michael Jackson and Bubbles',            'Retrato del icónico artista pop con su chimpancé en estilo porcelana.',                85000.00, '1988-01-01', '/uploads/pinturas/koons_michael_jackson.jpg',     'image/jpeg', 'koons_michael_jackson.jpg',     TRUE,  4),
(5, 'The Physical Impossibility of Death',    'Representación del tiburón en formol, obra más icónica de Hirst.',                    99000.00, '1991-01-01', '/uploads/pinturas/hirst_shark.jpg',               'image/jpeg', 'hirst_shark.jpg',               TRUE,  3),
(5, 'Beautiful Inside My Head Forever',       'Pintura de manchas de colores vibrantes sobre lienzo circular.',                       42000.00, '2008-01-01', '/uploads/pinturas/hirst_spin_painting.jpg',       'image/jpeg', 'hirst_spin_painting.jpg',       FALSE, 3),
(6, 'Superflat Monogram',                     'Fusión de la cultura pop japonesa con el arte contemporáneo occidental.',              75000.00, '2003-01-01', '/uploads/pinturas/murakami_superflat.jpg',        'image/jpeg', 'murakami_superflat.jpg',        FALSE, 4),
(6, 'My Lonesome Cowboy',                     'Figura anime en estilo neo-pop que mezcla manga con arte contemporáneo.',              82000.00, '1998-01-01', '/uploads/pinturas/murakami_cowboy.jpg',           'image/jpeg', 'murakami_cowboy.jpg',           TRUE,  4),
(6, 'In the Land of the Dead',                'Escena colorida inspirada en el folclore japonés y la cultura otaku.',                 58000.00, '2014-01-01', '/uploads/pinturas/murakami_land_dead.jpg',        'image/jpeg', 'murakami_land_dead.jpg',        FALSE, 4),
(7, 'COMPANION Passing Through',              'El icónico personaje COMPANION con ojos en X en pose reflexiva.',                     67000.00, '2013-01-01', '/uploads/pinturas/kaws_companion.jpg',            'image/jpeg', 'kaws_companion.jpg',            FALSE, 4),
(7, 'SHARE',                                  'Obra de KAWS explorando la conexión humana a través de sus personajes.',               54000.00, '2020-01-01', '/uploads/pinturas/kaws_share.jpg',                'image/jpeg', 'kaws_share.jpg',                FALSE, 3),
(8, 'Obama Hope Poster',                      'El famoso retrato en rojo, blanco y azul de Barack Obama.',                           38000.00, '2008-01-01', '/uploads/pinturas/fairey_obama_hope.jpg',         'image/jpeg', 'fairey_obama_hope.jpg',         FALSE, 2),
(8, 'Andre the Giant Has a Posse',            'El icónico sticker OBEY que marcó el inicio del arte callejero moderno.',             32000.00, '1989-01-01', '/uploads/pinturas/fairey_obey.jpg',               'image/jpeg', 'fairey_obey.jpg',               FALSE, 2),
(9, 'Past Times',                             'Escena de la vida afroamericana en un parque suburbano americano.',                   88000.00, '1997-01-01', '/uploads/pinturas/marshall_past_times.jpg',       'image/jpeg', 'marshall_past_times.jpg',       TRUE,  5),
(9, 'School of Beauty School of Culture',     'Salón de belleza afroamericano como espacio de identidad y cultura.',                 76000.00, '2012-01-01', '/uploads/pinturas/marshall_beauty_school.jpg',    'image/jpeg', 'marshall_beauty_school.jpg',    TRUE,  5),
(10,'Napoleon Leading the Army over the Alps','Reinterpretación del clásico napoleónico con un joven afroamericano.',                92000.00, '2005-01-01', '/uploads/pinturas/wiley_napoleon.jpg',            'image/jpeg', 'wiley_napoleon.jpg',            TRUE,  5),
(10,'Saint Jerome Hearing the Trumpet',       'Figura afroamericana en pose clásica de San Jerónimo.',                              71000.00, '2017-01-01', '/uploads/pinturas/wiley_saint_jerome.jpg',        'image/jpeg', 'wiley_saint_jerome.jpg',        TRUE,  5),
(10,'Equestrian Portrait of King Philip II',  'Reinterpretación del retrato ecuestre clásico con protagonista moderno.',             63000.00, '2006-01-01', '/uploads/pinturas/wiley_equestrian.jpg',          'image/jpeg', 'wiley_equestrian.jpg',          FALSE, 5);

-- ============================================
-- PINTURA_TECNICA
-- ============================================
INSERT INTO pintura_tecnica (id_pintura, id_tecnica) VALUES
(1,  4), 
(1,  1),
(2,  4),
(3,  4), 
(3,  3),
(4,  6),
(5,  6),
(6,  6), 
(6,  4),
(7,  1),
(8,  3),
(9,  3),
(10, 3),
(11, 3), 
(11, 4),
(12, 4),
(13, 3),
(14, 3), 
(14, 5),
(15, 3),
(16, 3), 
(16, 4),
(17, 3), 
(17, 4),
(18, 3),
(19, 5), 
(19, 6),
(20, 5), 
(20, 6),
(21, 1), 
(21, 3),
(22, 1),
(23, 1),
(24, 1),
(25, 1);

-- ============================================
-- TOURS
-- ============================================
INSERT INTO tour (id_guia, nombre, descripcion, fecha_inicio, fecha_fin, horario, precio) VALUES
(1, 'Neo-Expresionismo y Basquiat', 'Recorrido por las obras neo-expresionistas y el legado de Basquiat.',           '2025-01-10', '2025-01-10', '10:00 - 12:00', 150.00),
(4, 'Arte Callejero Global',        'Visita guiada por las obras de Banksy y Shepard Fairey.',                       '2025-01-17', '2025-01-17', '14:00 - 16:00', 120.00),
(1, 'Colecciones Exclusivas',       'Tour privado por las colecciones exclusivas más cotizadas de la galería.',      '2025-02-07', '2025-02-07', '09:00 - 11:00', 300.00),
(4, 'Arte Japonés Contemporáneo',   'Recorrido por las obras de Kusama y Murakami.',                                 '2025-02-14', '2025-02-14', '15:00 - 17:00', 180.00),
(1, 'Identidad y Cultura',          'Exploración de obras que redefinen la identidad cultural afroamericana.',       '2025-03-07', '2025-03-07', '11:00 - 13:00', 160.00),
(4, 'Noche de Arte Moderno',        'Tour nocturno especial con cóctel de bienvenida incluido.',                     '2025-03-21', '2025-03-21', '19:00 - 21:00', 250.00),
(1, 'Pop Art y Cultura Popular',    'Conexión entre el pop art clásico y el arte de Koons y KAWS.',                  '2025-04-04', '2025-04-04', '10:00 - 12:00', 140.00),
(4, 'Arte que Rompe Récords',       'Tour dedicado a las obras más costosas y controversiales del arte moderno.',    '2025-04-18', '2025-04-18', '14:00 - 16:00', 200.00),
(1, 'Historia del Arte Urbano',     'Conferencia y recorrido sobre la evolución del street art al arte de galería.', '2025-05-02', '2025-05-02', '09:00 - 12:00', 180.00),
(4, 'Arte Abstracto Contemporáneo', 'Introducción al arte abstracto con obras de Kusama y Hirst.',                   '2025-05-16', '2025-05-16', '15:00 - 17:00', 130.00);

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
(1,  2, '2025-01-15', 95000.00),
(3,  5, '2025-01-22', 90000.00),
(5,  2, '2025-02-03', 120000.00),
(7,  5, '2025-02-14', 88000.00),
(9,  2, '2025-02-28', 72000.00),
(1,  5, '2025-03-05', 99000.00),
(12, 2, '2025-03-12', 68000.00),
(3,  5, '2025-03-20', 82000.00),
(5,  2, '2025-04-02', 75000.00),
(7,  5, '2025-04-10', 67000.00),
(15, 2, '2025-04-18', 45000.00),
(9,  5, '2025-04-25', 88000.00),
(1,  2, '2025-05-01', 76000.00),
(3,  5, '2025-05-08', 42000.00),
(5,  2, '2025-05-15', 92000.00),
(7,  5, '2025-05-22', 71000.00),
(9,  2, '2025-06-01', 54000.00),
(12, 5, '2025-06-08', 38000.00),
(15, 2, '2025-06-15', 32000.00),
(1,  5, '2025-06-22', 58000.00),
(3,  2, '2025-07-01', 65000.00),
(5,  5, '2025-07-08', 63000.00),
(7,  2, '2025-07-15', 78000.00),
(9,  5, '2025-07-22', 55000.00),
(12, 2, '2025-07-29', 85000.00);

-- ============================================
-- DETALLE_VENTA
-- ============================================
INSERT INTO detalle_venta (id_venta, id_pintura, cantidad, precio_unitario) VALUES
(1,  1,  1, 95000.00),
(2,  7,  1, 90000.00),
(3,  10, 1, 120000.00),
(4,  4,  1, 88000.00),
(5,  5,  1, 72000.00),
(6,  12, 1, 99000.00),
(7,  8,  1, 68000.00),
(8,  15, 1, 82000.00),
(9,  14, 1, 75000.00),
(10, 17, 1, 67000.00),
(11, 9,  1, 45000.00),
(12, 21, 1, 88000.00),
(13, 22, 1, 76000.00),
(14, 13, 1, 42000.00),
(15, 23, 1, 92000.00),
(16, 24, 1, 71000.00),
(17, 18, 1, 54000.00),
(18, 19, 1, 38000.00),
(19, 20, 1, 32000.00),
(20, 16, 1, 58000.00),
(21, 3,  1, 65000.00),
(22, 25, 1, 63000.00),
(23, 2,  1, 78000.00),
(24, 11, 1, 55000.00),
(25, 6,  1, 85000.00);

-- ============================================
-- ENVÍOS
-- ============================================
INSERT INTO envio (id_venta, direccion_envio, fecha_envio, estado_envio) VALUES
(1,  'Zona 10, Calle Reforma 12-34, Ciudad de Guatemala',   '2025-01-18', 'entregado'),
(2,  'Zona 15, Vista Hermosa 3-45, Ciudad de Guatemala',    '2025-01-25', 'entregado'),
(3,  'Zona 16, Cayalá Torre A, Ciudad de Guatemala',        '2025-02-06', 'entregado'),
(4,  'Zona 10, Blvd Los Próceres, Ciudad de Guatemala',     '2025-02-17', 'entregado'),
(5,  'Zona 10, Oakland Mall área, Ciudad de Guatemala',     '2025-03-03', 'entregado'),
(6,  'Av. Brickell 1234, Apt 5B, Miami, Florida, USA',      '2025-03-08', 'entregado'),
(7,  'Zona 12, Colonia La Florida, Ciudad de Guatemala',    '2025-03-15', 'entregado'),
(8,  'Calle Serrano 45, Piso 3, Madrid, España',            '2025-03-25', 'en tránsito'),
(9,  'Zona 16, Cayalá Torre A, Ciudad de Guatemala',        '2025-04-05', 'entregado'),
(10, 'Via Condotti 12, Int 4, Roma, Italia',                '2025-04-15', 'entregado'),
(12, 'Zona 10, Oakland Mall área, Ciudad de Guatemala',     '2025-04-28', 'entregado'),
(13, 'Zona 10, Calle Reforma 12-34, Ciudad de Guatemala',   '2025-05-04', 'entregado'),
(15, 'Zona 16, Cayalá Torre A, Ciudad de Guatemala',        '2025-05-18', 'en tránsito'),
(16, 'Zona 10, Blvd Los Próceres, Ciudad de Guatemala',     '2025-05-25', 'entregado'),
(18, 'Zona 7, Tikal 2 Sector B, Ciudad de Guatemala',       '2025-06-11', 'entregado');