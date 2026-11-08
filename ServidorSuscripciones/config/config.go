package config

import (
	"os"
	"strings"
)

// Config agrupa la configuración del servidor leída desde variables de entorno.
type Config struct {
	Puerto          string
	DirectorioDatos string
}

// Cargar lee las variables de entorno y aplica valores por defecto.
//
//	PUERTO            puerto HTTP (por defecto 8080)
//	DIRECTORIO_DATOS  carpeta donde se persisten tokens y transacciones (por defecto ./datos)
func Cargar() Config {
	return Config{
		Puerto:          valorOPorDefecto("PUERTO", "8080"),
		DirectorioDatos: valorOPorDefecto("DIRECTORIO_DATOS", "./datos"),
	}
}

func valorOPorDefecto(clave, porDefecto string) string {
	if v := strings.TrimSpace(os.Getenv(clave)); v != "" {
		return v
	}
	return porDefecto
}
