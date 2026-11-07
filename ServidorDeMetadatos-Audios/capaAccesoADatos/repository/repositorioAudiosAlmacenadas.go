/**
 * @file repositorioAudiosAlmacenadas.go
 * @brief Definición del repositorio para almacenar y recuperar archivos de audio físicos.
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
 * @brief Estructura que representa el repositorio de audios.
 *
 * Gestiona el acceso concurrente al almacenamiento físico de los audios
 * utilizando un Mutex para evitar condiciones de carrera.
 */
type RepositorioAudios struct {
	mu sync.Mutex
}

var (
	instancia *RepositorioAudios
	once      sync.Once
)

/**
 * @brief Obtiene la instancia única del repositorio de audios (Patrón Singleton).
 *
 * Garantiza que solo exista una instancia de RepositorioAudios en toda la aplicación.
 *
 * @return *RepositorioAudios Puntero a la instancia única del repositorio.
 */
func GetRepositorioAudios() *RepositorioAudios {
	once.Do(func() {
		instancia = &RepositorioAudios{}
	})
	return instancia
}

/**
 * @brief Guarda un archivo de audio físico en formato .mp3 en la carpeta 'audios'.
 *
 * Crea el directorio 'audios' si este no existe y escribe los bytes proporcionados
 * en un nuevo archivo dentro de este directorio.
 *
 * @param nombreArchivo Nombre del archivo a guardar (incluyendo extensión, ej. 'audio.mp3').
 * @param data Array de bytes que representa el contenido del archivo de audio.
 * @return error Retorna un error si hubo problemas al crear el directorio o al escribir el archivo, nil en caso de éxito.
 */
func (r *RepositorioAudios) GuardarAudio(nombreArchivo string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Crear carpeta si no existe
	if err := os.MkdirAll("audios", os.ModePerm); err != nil {
		return fmt.Errorf("error al crear directorio audios: %v", err)
	}

	filePath := filepath.Join("audios", nombreArchivo)

	// Guardar archivo fisico
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar el archivo: %v", err)
	}

	return nil
}

/**
 * @brief Lee y devuelve el contenido de un archivo de audio físico desde la carpeta 'audios'.
 *
 * @param nombreArchivo Nombre del archivo de audio a leer (incluyendo extensión).
 * @return []byte Array de bytes con el contenido del archivo de audio.
 * @return error Retorna un error si el archivo no pudo ser leído o encontrado, nil en caso de éxito.
 */
func (r *RepositorioAudios) ObtenerAudio(nombreArchivo string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join("audios", nombreArchivo)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo %s: %v", nombreArchivo, err)
	}

	return data, nil
}
