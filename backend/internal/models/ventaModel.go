package models

type Venta struct {
	ID         int    `json:"id_venta"`
	IDCliente  int    `json:"id_cliente"`
	IDEmpleado int    `json:"id_empleado"`
	FechaVenta string `json:"fecha_venta"`
	Precio     string `json:"precio"`
}

type VentaRequest struct {
	IDCliente  int    `json:"id_cliente"`
	IDEmpleado int    `json:"id_empleado"`
	FechaVenta string `json:"fecha_venta"`
	Precio     string `json:"precio"`
}

type DetalleVenta struct {
	ID             int    `json:"id_detalle_venta"`
	IDVenta        int    `json:"id_venta"`
	IDPintura      int    `json:"id_pintura"`
	TituloPintura  string `json:"titulo_pintura,omitempty"`
	Cantidad       int    `json:"cantidad"`
	PrecioUnitario string `json:"precio_unitario"`
}

type DetalleVentaRequest struct {
	IDVenta        int    `json:"id_venta"`
	IDPintura      int    `json:"id_pintura"`
	Cantidad       int    `json:"cantidad"`
	PrecioUnitario string `json:"precio_unitario"`
}
