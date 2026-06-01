package models

type UsuarioEntity struct {
	IDUsuario         int    `gorm:"column:id_usuario;primaryKey;autoIncrement"`
	Nombre            string `gorm:"column:nombre"`
	Apellido          string `gorm:"column:apellido"`
	CorreoElectronico string `gorm:"column:correo_electronico"`
	Telefono          string `gorm:"column:telefono"`
	Contrasena        string `gorm:"column:contrasena"`
}

func (UsuarioEntity) TableName() string { return "usuario" }

type ClienteEntity struct {
	IDCliente   int    `gorm:"column:id_cliente;primaryKey;autoIncrement"`
	IDUsuario   int    `gorm:"column:id_usuario"`
	TipoCliente string `gorm:"column:tipo_cliente"`
}

func (ClienteEntity) TableName() string { return "cliente" }

type EmpleadoEntity struct {
	IDEmpleado   int    `gorm:"column:id_empleado;primaryKey;autoIncrement"`
	IDUsuario    int    `gorm:"column:id_usuario"`
	TipoEmpleado string `gorm:"column:tipo_empleado"`
}

func (EmpleadoEntity) TableName() string { return "empleado" }

type ArtistaEntity struct {
	IDArtista      int    `gorm:"column:id_artista;primaryKey;autoIncrement"`
	NombreCompleto string `gorm:"column:nombre_completo"`
	Nacionalidad   string `gorm:"column:nacionalidad"`
	IDReclutador   int    `gorm:"column:id_reclutador"`
}

func (ArtistaEntity) TableName() string { return "artista" }

type ColeccionEntity struct {
	IDColeccion      int    `gorm:"column:id_coleccion;primaryKey;autoIncrement"`
	Nombre           string `gorm:"column:nombre"`
	Descripcion      string `gorm:"column:descripcion"`
	Exclusiva        bool   `gorm:"column:exclusiva"`
	FechaLanzamiento string `gorm:"column:fecha_lanzamiento"`
}

func (ColeccionEntity) TableName() string { return "coleccion" }

type PinturaEntity struct {
	IDPintura     int     `gorm:"column:id_pintura;primaryKey;autoIncrement"`
	IDArtista     int     `gorm:"column:id_artista"`
	Titulo        string  `gorm:"column:titulo"`
	Descripcion   string  `gorm:"column:descripcion"`
	Precio        float64 `gorm:"column:precio"`
	FechaCreacion string  `gorm:"column:fecha_creacion"`
	ImagenPath    string  `gorm:"column:imagen_path"`
	ImagenTipo    string  `gorm:"column:imagen_tipo"`
	ImagenNombre  string  `gorm:"column:imagen_nombre"`
	Exclusiva     bool    `gorm:"column:exclusiva"`
	IDColeccion   *int    `gorm:"column:id_coleccion"`
}

func (PinturaEntity) TableName() string { return "pintura" }

type TecnicaEntity struct {
	IDTecnica   int    `gorm:"column:id_tecnica;primaryKey;autoIncrement"`
	Nombre      string `gorm:"column:nombre"`
	Descripcion string `gorm:"column:descripcion"`
}

func (TecnicaEntity) TableName() string { return "tecnica" }

type PinturaTecnicaEntity struct {
	IDPinturaTecnica int `gorm:"column:id_pintura_tecnica;primaryKey;autoIncrement"`
	IDPintura        int `gorm:"column:id_pintura"`
	IDTecnica        int `gorm:"column:id_tecnica"`
}

func (PinturaTecnicaEntity) TableName() string { return "pintura_tecnica" }

type TourEntity struct {
	IDTour      int    `gorm:"column:id_tour;primaryKey;autoIncrement"`
	IDGuia      int    `gorm:"column:id_guia"`
	Nombre      string `gorm:"column:nombre"`
	Descripcion string `gorm:"column:descripcion"`
	FechaInicio string `gorm:"column:fecha_inicio"`
	FechaFin    string `gorm:"column:fecha_fin"`
	Horario     string `gorm:"column:horario"`
	Precio      string `gorm:"column:precio"`
}

func (TourEntity) TableName() string { return "tour" }

type ClienteTourEntity struct {
	IDClienteTour int    `gorm:"column:id_cliente_tour;primaryKey;autoIncrement"`
	IDCliente     int    `gorm:"column:id_cliente"`
	IDTour        int    `gorm:"column:id_tour"`
	FechaReserva  string `gorm:"column:fecha_reserva"`
}

func (ClienteTourEntity) TableName() string { return "cliente_tour" }

type VentaEntity struct {
	IDVenta    int    `gorm:"column:id_venta;primaryKey;autoIncrement"`
	IDCliente  int    `gorm:"column:id_cliente"`
	IDEmpleado int    `gorm:"column:id_empleado"`
	FechaVenta string `gorm:"column:fecha_venta"`
	Precio     string `gorm:"column:precio"`
}

func (VentaEntity) TableName() string { return "venta" }

type DetalleVentaEntity struct {
	IDDetalleVenta int    `gorm:"column:id_detalle_venta;primaryKey;autoIncrement"`
	IDVenta        int    `gorm:"column:id_venta"`
	IDPintura      int    `gorm:"column:id_pintura"`
	Cantidad       int    `gorm:"column:cantidad"`
	PrecioUnitario string `gorm:"column:precio_unitario"`
}

func (DetalleVentaEntity) TableName() string { return "detalle_venta" }

type EnvioEntity struct {
	IDEnvio        int    `gorm:"column:id_envio;primaryKey;autoIncrement"`
	IDVenta        int    `gorm:"column:id_venta"`
	DireccionEnvio string `gorm:"column:direccion_envio"`
	FechaEnvio     string `gorm:"column:fecha_envio"`
	EstadoEnvio    string `gorm:"column:estado_envio"`
}

func (EnvioEntity) TableName() string { return "envio" }
