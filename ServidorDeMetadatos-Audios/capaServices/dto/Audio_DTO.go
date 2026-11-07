/**
 * @file Audio_DTO.go
 * @brief Definición del DTO para los metadatos de audios.
 */
package dto

/**
 * @struct AudioDTO
 * @brief Representa los metadatos básicos opcionales recibidos al subir o consultar un audio.
 */
type AudioDTO struct {
	NombreArchivo string `json:"nombre_archivo"`
}
