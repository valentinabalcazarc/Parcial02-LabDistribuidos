/**
 * @file Audio_service.go
 * @brief Definición del servicio para la lógica de almacenamiento de archivos de audio.
 */
package services

import (
	"fmt"
	capaaccesoadatos "servidorMetadatos/capaAccesoADatos/repository"
	dtos "servidorMetadatos/capaServices/dto"
)

/**
 * @struct AudioService
 * @brief Servicio centralizado que gestiona el almacenamiento de archivos de audio en disco.
 * 
 * Permite que los controladores y servicios de metadatos deleguen el guardado de archivos
 * físicos sin duplicar código.
 */
type AudioService struct {
	repo *capaaccesoadatos.RepositorioAudios
}

/**
 * @brief Crea una nueva instancia de AudioService.
 *
 * @return *AudioService Puntero a la nueva instancia del servicio.
 */
func NewAudioService() *AudioService {
	fmt.Println("Inicializando service de almacenamiento de audios")
	return &AudioService{
		repo: capaaccesoadatos.GetRepositorioAudios(),
	}
}

/**
 * @brief Delega la acción de guardar un audio al repositorio mediante un DTO.
 *
 * @param objAudio Objeto de transferencia de datos con la información del audio.
 * @param data Array de bytes que representa el contenido del audio.
 * @return error Retorna un error si ocurre un problema al guardar, de lo contrario nil.
 */
func (thisS *AudioService) GuardarAudio(objAudio dtos.AudioDTO, data []byte) error {
	return thisS.repo.GuardarAudio(objAudio.NombreArchivo, data)
}

/**
 * @brief Sobrecarga conveniente para guardar un archivo pasando directamente su nombre y bytes.
 *
 * @param nombreArchivo Nombre del archivo a guardar (ej. "cancion1.mp3").
 * @param data Contenido binario del archivo.
 * @return error Error si no se pudo guardar.
 */
func (thisS *AudioService) GuardarArchivo(nombreArchivo string, data []byte) error {
	return thisS.repo.GuardarAudio(nombreArchivo, data)
}

/**
 * @brief Obtiene los bytes de un archivo físico para uso interno del servidor.
 *
 * @param nombreArchivo Nombre del archivo de audio a obtener.
 * @return []byte Bytes leídos del archivo.
 * @return error Error si el archivo no existe.
 */
func (thisS *AudioService) ObtenerAudio(nombreArchivo string) ([]byte, error) {
	return thisS.repo.ObtenerAudio(nombreArchivo)
}
