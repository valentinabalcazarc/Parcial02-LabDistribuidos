/**
 * @file repositorioStreaming.go
 * @brief Capa de acceso a datos para el servidor de streaming.
 */
package capaaccesoadatos

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

/**
 * @struct RepositorioStreaming
 * @brief Estructura que representa el repositorio de acceso a archivos de streaming.
 */
type RepositorioStreaming struct{}

var (
	instancia *RepositorioStreaming
	once      sync.Once
)

/**
 * @brief Retorna la instancia única (Singleton) del repositorio de streaming.
 * @return *RepositorioStreaming Puntero a la instancia del repositorio.
 */
func GetRepositorioStreaming() *RepositorioStreaming {
	once.Do(func() {
		instancia = &RepositorioStreaming{}
	})
	return instancia
}

/**
 * @brief Abre el archivo físico .mp3 desde la carpeta 'audios'.
 * @param nombreArchivo Nombre del archivo a abrir.
 * @return *os.File Puntero al archivo abierto.
 * @return error Error en caso de que no se pueda abrir el archivo.
 */
func (r *RepositorioStreaming) AbrirArchivo(nombreArchivo string) (*os.File, error) {
	// Se busca en la carpeta audios del proyecto
	ruta := filepath.Join("canciones", nombreArchivo)

	file, err := os.Open(ruta)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo de audio %s: %v", nombreArchivo, err)
	}

	return file, nil
}
