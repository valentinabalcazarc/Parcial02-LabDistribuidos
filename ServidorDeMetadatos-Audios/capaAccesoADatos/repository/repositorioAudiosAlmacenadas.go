/**
 * @file repositorioAudiosAlmacenadas.go
 * @brief Definición del repositorio para almacenar archivos de audio en la carpeta del servidor de streaming.
 */
package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

/**
 * @struct RepositorioAudios
 * @brief Estructura que gestiona el guardado de los archivos físicos de audio.
 */
type RepositorioAudios struct {
	mu         sync.Mutex
	directorio string
}

var (
	instancia *RepositorioAudios
	once      sync.Once
)

/**
 * @brief Obtiene la instancia única del repositorio de audios (Singleton).
 *
 * Apunta directamente a la carpeta 'canciones' del ServidorDeStreaming.
 */
func GetRepositorioAudios() *RepositorioAudios {
	once.Do(func() {
		instancia = &RepositorioAudios{
			directorio: "../ServidorDeStreaming/canciones",
		}
	})
	return instancia
}

/**
 * @brief Guarda el archivo .mp3 directamente en la carpeta de canciones del streaming.
 *
 * @param nombreArchivo Nombre del archivo
 * @param data Bytes del archivo de audio.
 * @return error Error si no se pudo escribir el archivo.
 */
func (r *RepositorioAudios) GuardarAudio(nombreArchivo string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Crear carpeta si no existe
	if err := os.MkdirAll(r.directorio, os.ModePerm); err != nil {
		return fmt.Errorf("error al crear directorio %s: %v", r.directorio, err)
	}

	nombreSeguro := filepath.Base(nombreArchivo)
	filePath := filepath.Join(r.directorio, nombreSeguro)

	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar el archivo: %v", err)
	}

	fmt.Printf("[RepositorioAudios] Archivo físico guardado en: %s\n", filePath)
	return nil
}

/**
 * @brief Lee un archivo físico del disco para uso interno.
 */
func (r *RepositorioAudios) ObtenerAudio(nombreArchivo string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	nombreSeguro := filepath.Base(nombreArchivo)
	filePath := filepath.Join(r.directorio, nombreSeguro)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo %s: %v", filePath, err)
	}

	return data, nil
}
