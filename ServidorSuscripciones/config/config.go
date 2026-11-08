package config

import (
	"os"
	"strings"
)

// Config agrupa la configuración del servidor leída desde variables de entorno.
type Config struct {
	Puerto          string
	DirectorioDatos string
	ConsulAddr      string
	NombreServicio  string
	Host            string
}

// Cargar lee las variables de entorno y aplica valores por defecto.
//
//	PUERTO            puerto HTTP
//	DIRECTORIO_DATOS  carpeta donde se persisten tokens y transacciones (por defecto ./datos)
func Cargar() Config {
	return Config{
		Puerto:          valorOPorDefecto("PUERTO", "8083"),
		DirectorioDatos: valorOPorDefecto("DIRECTORIO_DATOS", "./datos"),
		ConsulAddr:      valorOPorDefecto("CONSUL_ADDR", "127.0.0.1:8500"),
		NombreServicio:  valorOPorDefecto("NOMBRE_SERVICIO", "ms-servidorSuscripciones"),
		Host:            valorOPorDefecto("HOST_SERVICIO", "127.0.0.1"),
	}
}

func valorOPorDefecto(clave, porDefecto string) string {
	if v := strings.TrimSpace(os.Getenv(clave)); v != "" {
		return v
	}
	return porDefecto
}
