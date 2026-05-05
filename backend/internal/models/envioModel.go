package models

type Envio struct {
	ID             int    `json:"id_envio"`
	IDVenta        int    `json:"id_venta"`
	DireccionEnvio string `json:"direccion_envio"`
	FechaEnvio     string `json:"fecha_envio"`
	EstadoEnvio    string `json:"estado_envio"`
}

type EnvioRequest struct {
	IDVenta        int    `json:"id_venta"`
	DireccionEnvio string `json:"direccion_envio"`
	FechaEnvio     string `json:"fecha_envio"`
	EstadoEnvio    string `json:"estado_envio"`
}
