/**
 * @file Audio_service.go
 * @brief Definición del servicio para la lógica de negocio de los audios.
 */
package services

import (
	"fmt"
	capaaccesoadatos "servidorMetadatos/capaAccesoADatos/repository"
	dtos "servidorMetadatos/capaServices/dto"
)

/**
 * @struct AudioService
 * @brief Servicio que gestiona la lógica de almacenamiento y recuperación de audios.
 */
type AudioService struct {
	repo *capaaccesoadatos.RepositorioAudios
}

/**
 * @brief Crea una nueva instancia de AudioService.
 *
 * Inicializa el servicio obteniendo la instancia única del repositorio de audios.
 *
 * @return *AudioService Puntero a la nueva instancia del servicio.
 */
func NewAudioService() *AudioService {
	fmt.Println("Inicializando service de almacenamiento")
	return &AudioService{
		repo: capaaccesoadatos.GetRepositorioAudios(),
	}
}

/**
 * @brief Delega la acción de guardar un audio al repositorio.
 *
 * @param objAudio Objeto de transferencia de datos con la información del audio.
 * @param data Array de bytes que representa el contenido del audio.
 * @return error Retorna un error si ocurre un problema al guardar, de lo contrario nil.
 */
func (thisS *AudioService) GuardarAudio(objAudio dtos.AudioDTO, data []byte) error {
	return thisS.repo.GuardarAudio(objAudio.NombreArchivo, data)
}

/**
 * @brief Obtiene los bytes del archivo mp3 por su nombre desde el repositorio.
 *
 * @param nombreArchivo Nombre del archivo de audio a obtener.
 * @return []byte Array de bytes con el contenido del archivo de audio.
 * @return error Retorna un error si el archivo no es encontrado, de lo contrario nil.
 */
func (s *AudioService) ObtenerAudio(nombreArchivo string) ([]byte, error) {
	return s.repo.ObtenerAudio(nombreArchivo)
}
