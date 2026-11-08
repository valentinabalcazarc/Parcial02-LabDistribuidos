package config

// Ruta asocia un prefijo público del gateway con el nombre de un servicio en Consul.
type Ruta struct {
	Prefijo  string // ruta pública en el gateway, p. ej. /api/productos
	Servicio string // nombre con el que el microservicio se registró en Consul
	Destino  string // ruta interna en el microservicio, p. ej. /productos
}

// Las rutas que el Gateway enruta a los microservicios REST registrados en Consul.
// NOTA: El streaming gRPC va directo por Envoy, no necesita ruta aquí.
var Rutas = []Ruta{
	// ServidorDeMetadatos-Audios (puerto :8081)
	// Registrado en Consul como "ms-servidorDeMetadatos-Audios"
	{Prefijo: "/api/metadatosaudios", Servicio: "ms-servidorDeMetadatos-Audios", Destino: "/metadatosaudios"},

	// ServidorDeEstadisticas (puerto :8082)
	// Registrado en Consul como "ms-servidorDeEstadisticas"
	{Prefijo: "/api/estadisticas", Servicio: "ms-servidorDeEstadisticas", Destino: "/estadisticas"},

	// ServidorDeReacciones (puerto :5000)
	// Registrado en Consul como "ms-servidorDeReacciones"
	{Prefijo: "/api/reacciones", Servicio: "ms-servidorDeReacciones", Destino: "/reacciones"},

	// ServidorSuscripciones (puerto :8083)
	// Registrado en Consul como "ms-servidorSuscripciones"
	{Prefijo: "/api/pasarela", Servicio: "ms-servidorSuscripciones", Destino: "/api"},
}
