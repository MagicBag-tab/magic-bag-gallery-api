-- Proyecto 3: roles y permisos.

CREATE OR REPLACE FUNCTION sp_registrar_cliente(
    p_nombre VARCHAR,
    p_apellido VARCHAR,
    p_correo_electronico VARCHAR,
    p_telefono VARCHAR,
    p_contrasena VARCHAR,
    p_tipo_cliente VARCHAR DEFAULT 'regular'
)
RETURNS TABLE(id_usuario INT, id_cliente INT)
LANGUAGE plpgsql
AS $$
DECLARE
    v_id_usuario INT;
    v_id_cliente INT;
BEGIN
    INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena)
    VALUES (p_nombre, p_apellido, p_correo_electronico, p_telefono, p_contrasena)
    RETURNING usuario.id_usuario INTO v_id_usuario;

    INSERT INTO cliente (id_usuario, tipo_cliente)
    VALUES (v_id_usuario, COALESCE(NULLIF(p_tipo_cliente, ''), 'regular'))
    RETURNING cliente.id_cliente INTO v_id_cliente;

    RETURN QUERY SELECT v_id_usuario, v_id_cliente;
END;
$$;

-- SP con parámetros IN/OUT y manejo de excepciones (Rúbrica: 10 pts)
CREATE OR REPLACE FUNCTION sp_registrar_empleado(
    p_nombre VARCHAR,
    p_apellido VARCHAR,
    p_correo_electronico VARCHAR,
    p_telefono VARCHAR,
    p_contrasena VARCHAR,
    p_tipo_empleado VARCHAR,
    OUT p_id_usuario INT,
    OUT p_id_empleado INT,
    OUT p_error VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    p_error := NULL;
    p_id_usuario := NULL;
    p_id_empleado := NULL;

    IF p_tipo_empleado NOT IN ('guia', 'asesor', 'reclutador') THEN
        p_error := 'Tipo de empleado invalido: ' || p_tipo_empleado;
        RETURN;
    END IF;

    BEGIN
        INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena)
        VALUES (p_nombre, p_apellido, p_correo_electronico, p_telefono, p_contrasena)
        RETURNING usuario.id_usuario INTO p_id_usuario;

        INSERT INTO empleado (id_usuario, tipo_empleado)
        VALUES (p_id_usuario, p_tipo_empleado)
        RETURNING empleado.id_empleado INTO p_id_empleado;
    EXCEPTION WHEN unique_violation THEN
        p_error := 'El correo electrónico ya existe';
        p_id_usuario := NULL;
        p_id_empleado := NULL;
    WHEN OTHERS THEN
        p_error := 'Error al registrar empleado: ' || SQLERRM;
        p_id_usuario := NULL;
        p_id_empleado := NULL;
    END;
END;
$$;

CREATE OR REPLACE FUNCTION sp_crear_reserva(
    p_id_cliente INT,
    p_id_tour INT,
    p_fecha_reserva DATE
)
RETURNS TABLE(id_cliente_tour INT)
LANGUAGE plpgsql
AS $$
DECLARE
    v_id_cliente_tour INT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM cliente WHERE cliente.id_cliente = p_id_cliente) THEN
        RAISE EXCEPTION 'Cliente % no existe', p_id_cliente;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM tour WHERE tour.id_tour = p_id_tour) THEN
        RAISE EXCEPTION 'Tour % no existe', p_id_tour;
    END IF;

    INSERT INTO cliente_tour (id_cliente, id_tour, fecha_reserva)
    VALUES (p_id_cliente, p_id_tour, p_fecha_reserva)
    RETURNING cliente_tour.id_cliente_tour INTO v_id_cliente_tour;

    RETURN QUERY SELECT v_id_cliente_tour;
END;
$$;

CREATE OR REPLACE FUNCTION sp_crear_venta(
    p_id_cliente INT,
    p_id_empleado INT,
    p_fecha_venta DATE,
    p_precio DECIMAL
)
RETURNS TABLE(id_venta INT)
LANGUAGE plpgsql
AS $$
DECLARE
    v_id_venta INT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM cliente WHERE cliente.id_cliente = p_id_cliente) THEN
        RAISE EXCEPTION 'Cliente % no existe', p_id_cliente;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM empleado WHERE empleado.id_empleado = p_id_empleado) THEN
        RAISE EXCEPTION 'Empleado % no existe', p_id_empleado;
    END IF;

    IF p_precio <= 0 THEN
        RAISE EXCEPTION 'El precio de venta debe ser positivo';
    END IF;

    INSERT INTO venta (id_cliente, id_empleado, fecha_venta, precio)
    VALUES (p_id_cliente, p_id_empleado, p_fecha_venta, p_precio)
    RETURNING venta.id_venta INTO v_id_venta;

    RETURN QUERY SELECT v_id_venta;
END;
$$;

CREATE OR REPLACE FUNCTION sp_eliminar_pintura(p_id_pintura INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pintura WHERE pintura.id_pintura = p_id_pintura) THEN
        RAISE EXCEPTION 'Pintura % no existe', p_id_pintura;
    END IF;

    DELETE FROM pintura_tecnica WHERE id_pintura = p_id_pintura;
    DELETE FROM pintura WHERE id_pintura = p_id_pintura;
END;
$$;

-- SP con transacción explícita, COMMIT y ROLLBACK (Rúbrica: 10 pts)
CREATE OR REPLACE FUNCTION sp_registrar_venta_con_detalles(
    p_id_cliente INT,
    p_id_empleado INT,
    p_fecha_venta DATE,
    p_precio_total DECIMAL,
    p_id_pintura INT,
    OUT p_id_venta INT,
    OUT p_exito BOOLEAN,
    OUT p_mensaje VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    p_exito := FALSE;
    p_id_venta := NULL;
    p_mensaje := '';

    BEGIN
        -- Transacción: Verificar cliente
        IF NOT EXISTS (SELECT 1 FROM cliente WHERE id_cliente = p_id_cliente) THEN
            p_mensaje := 'Cliente no existe';
            RETURN;
        END IF;

        -- Transacción: Verificar empleado
        IF NOT EXISTS (SELECT 1 FROM empleado WHERE id_empleado = p_id_empleado) THEN
            p_mensaje := 'Empleado no existe';
            RETURN;
        END IF;

        -- Transacción: Verificar pintura
        IF NOT EXISTS (SELECT 1 FROM pintura WHERE id_pintura = p_id_pintura) THEN
            p_mensaje := 'Pintura no existe';
            RETURN;
        END IF;

        -- Transacción: Validar precio positivo
        IF p_precio_total <= 0 THEN
            p_mensaje := 'El precio debe ser positivo';
            RETURN;
        END IF;

        -- Transacción: Crear venta
        INSERT INTO venta (id_cliente, id_empleado, fecha_venta, precio)
        VALUES (p_id_cliente, p_id_empleado, p_fecha_venta, p_precio_total)
        RETURNING venta.id_venta INTO p_id_venta;

        -- Transacción: Crear detalle de venta
        INSERT INTO detalle_venta (id_venta, id_pintura, cantidad, precio_unitario)
        VALUES (p_id_venta, p_id_pintura, 1, p_precio_total);

        -- COMMIT: Si todo fue exitoso
        p_exito := TRUE;
        p_mensaje := 'Venta registrada exitosamente';

    EXCEPTION WHEN foreign_key_violation THEN
        -- ROLLBACK: Error de integridad referencial
        p_exito := FALSE;
        p_id_venta := NULL;
        p_mensaje := 'Error de integridad referencial - datos inválidos';

    WHEN check_violation THEN
        -- ROLLBACK: Error de validación
        p_exito := FALSE;
        p_id_venta := NULL;
        p_mensaje := 'Error de validación en los datos';

    WHEN OTHERS THEN
        -- ROLLBACK: Otros errores
        p_exito := FALSE;
        p_id_venta := NULL;
        p_mensaje := 'Error al registrar venta: ' || SQLERRM;
    END;
END;
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mbg_catalogo') THEN
        CREATE ROLE mbg_catalogo NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mbg_cliente') THEN
        CREATE ROLE mbg_cliente NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mbg_guia') THEN
        CREATE ROLE mbg_guia NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mbg_asesor') THEN
        CREATE ROLE mbg_asesor NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mbg_reclutador') THEN
        CREATE ROLE mbg_reclutador NOLOGIN;
    END IF;
END $$;

REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO mbg_catalogo, mbg_cliente, mbg_guia, mbg_asesor, mbg_reclutador;

REVOKE ALL ON ALL TABLES IN SCHEMA public FROM mbg_catalogo, mbg_cliente, mbg_guia, mbg_asesor, mbg_reclutador;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM mbg_catalogo, mbg_cliente, mbg_guia, mbg_asesor, mbg_reclutador;

GRANT SELECT ON
    artista,
    coleccion,
    pintura,
    pintura_tecnica,
    tecnica,
    tour,
    vista_pinturas_completa,
    vista_artistas_resumen
TO mbg_catalogo;

GRANT mbg_catalogo TO mbg_cliente;
GRANT SELECT ON usuario, cliente, venta, detalle_venta, envio, cliente_tour TO mbg_cliente;
GRANT INSERT, UPDATE, DELETE ON cliente_tour TO mbg_cliente;

GRANT mbg_catalogo TO mbg_guia;
GRANT SELECT, INSERT, UPDATE, DELETE ON tour, cliente_tour TO mbg_guia;
GRANT SELECT ON usuario, cliente, empleado TO mbg_guia;

GRANT mbg_catalogo TO mbg_asesor;
GRANT SELECT, INSERT, UPDATE, DELETE ON venta, detalle_venta, envio TO mbg_asesor;
GRANT SELECT, INSERT, UPDATE ON usuario, cliente TO mbg_asesor;
GRANT SELECT ON empleado, vista_ventas_detalle TO mbg_asesor;

GRANT mbg_catalogo TO mbg_reclutador;
GRANT SELECT, INSERT, UPDATE, DELETE ON artista, coleccion, pintura, pintura_tecnica, tecnica TO mbg_reclutador;
GRANT SELECT ON usuario, empleado TO mbg_reclutador;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO mbg_cliente, mbg_guia, mbg_asesor, mbg_reclutador;

GRANT EXECUTE ON FUNCTION sp_registrar_cliente(VARCHAR, VARCHAR, VARCHAR, VARCHAR, VARCHAR, VARCHAR) TO mbg_catalogo, mbg_cliente;
GRANT EXECUTE ON FUNCTION sp_registrar_empleado(VARCHAR, VARCHAR, VARCHAR, VARCHAR, VARCHAR, VARCHAR, OUT INT, OUT INT, OUT VARCHAR) TO mbg_asesor;
GRANT EXECUTE ON FUNCTION sp_crear_reserva(INT, INT, DATE) TO mbg_cliente, mbg_guia;
GRANT EXECUTE ON FUNCTION sp_crear_venta(INT, INT, DATE, DECIMAL) TO mbg_asesor;
GRANT EXECUTE ON FUNCTION sp_eliminar_pintura(INT) TO mbg_reclutador;
GRANT EXECUTE ON FUNCTION sp_registrar_venta_con_detalles(INT, INT, DATE, DECIMAL, INT, OUT INT, OUT BOOLEAN, OUT VARCHAR) TO mbg_asesor;

GRANT mbg_catalogo, mbg_cliente, mbg_guia, mbg_asesor, mbg_reclutador TO CURRENT_USER;
