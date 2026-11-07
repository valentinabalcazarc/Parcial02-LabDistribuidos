package dto

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Objeto de Transferencia de Datos para Ruido Blanco.
 * 
 * Contiene los datos que son transferidos hacia o desde las capas externas.
 */
type MetadataRuidoBlancoDTO struct {
	ID                  int              `json:"id"`
	Tipo                entity.TipoAudio `json:"tipo"`
	TipoSonido          string           `json:"tipo_sonido"`
	FuenteAudio         string           `json:"fuente_audio"`
	UsoSugerido         string           `json:"uso_sugerido"`
	Proveedor           string           `json:"proveedor"`
	Duracion            int              `json:"duracion"`
	FrecuenciaDominante string           `json:"frecuencia_dominante"`
}

/**
 * @brief DTO para la respuesta de una consulta de un solo ruido blanco.
 */
type RespuestaMetadataRuidoBlancoDTO struct {
	ObjRuidoBlanco MetadataRuidoBlancoDTO `json:"objRuidoBlanco"`
	Codigo         int                    `json:"codigo"`
	Mensaje        string                 `json:"mensaje"`
}

/**
 * @brief DTO para la respuesta de un listado de ruidos blancos.
 */
type RespuestaListaRuidoBlancoDTO struct {
	RuidosBlancos []MetadataRuidoBlancoDTO `json:"ruidosBlancos"`
	Codigo        int                      `json:"codigo"`
	Mensaje       string                   `json:"mensaje"`
}
