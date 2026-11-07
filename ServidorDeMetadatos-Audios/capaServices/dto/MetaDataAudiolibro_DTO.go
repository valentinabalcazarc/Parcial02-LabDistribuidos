package dto

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Objeto de Transferencia de Datos para Audiolibro.
 * 
 * Contiene los datos que son transferidos hacia o desde las capas externas.
 */
type MetadataAudiolibroDTO struct {
	ID          int              `json:"id"`
	Tipo        entity.TipoAudio `json:"tipo"`
	TituloLibro string           `json:"titulo_libro"`
	Autor       string           `json:"autor"`
	Narrador    string           `json:"narrador"`
	Editorial   string           `json:"editorial"`
	ISBN        int              `json:"isbn"`
	Capitulo    int              `json:"capitulo"`
}

/**
 * @brief DTO para la respuesta de una consulta de un solo audiolibro.
 */
type RespuestaMetadataAudiolibroDTO struct {
	ObjAudiolibro MetadataAudiolibroDTO `json:"objAudiolibro"`
	Codigo        int                   `json:"codigo"`
	Mensaje       string                `json:"mensaje"`
}

/**
 * @brief DTO para la respuesta de un listado de audiolibros.
 */
type RespuestaListaAudiolibroDTO struct {
	Audiolibros []MetadataAudiolibroDTO `json:"audiolibros"`
	Codigo      int                     `json:"codigo"`
	Mensaje     string                  `json:"mensaje"`
}
