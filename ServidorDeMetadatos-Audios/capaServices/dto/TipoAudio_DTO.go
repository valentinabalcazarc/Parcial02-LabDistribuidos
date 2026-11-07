package dto

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief DTO de respuesta para la lista de tipos de audio.
 * 
 * Contiene el código, mensaje y los tipos de audio.
 */
type RespuestaTiposAudioDTO struct {
	Tipos   []entity.TipoAudio `json:"tipos"`
	Codigo  int                `json:"codigo"`
	Mensaje string             `json:"mensaje"`
}
