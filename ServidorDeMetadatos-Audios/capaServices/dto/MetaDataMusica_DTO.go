package dto

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Objeto de Transferencia de Datos para Música.
 * 
 * Contiene los datos que son transferidos hacia o desde las capas externas.
 */
type MetadataMusicaDTO struct {
	ID                int              `json:"id"`
	Tipo              entity.TipoAudio `json:"tipo"`
	ArtistaPrincipal  string           `json:"artista_principal"`
	Album             string           `json:"album"`
	Genero            string           `json:"genero"`
	TituloCancion     string           `json:"titulo_cancion"`
	SelloDiscografico string           `json:"sello_discografico"`
	AnioLanzamiento   string           `json:"anio_lanzamiento"`
}

/**
 * @brief DTO para la respuesta de una consulta de una sola música.
 */
type RespuestaMetadataMusicaDTO struct {
	ObjMusica MetadataMusicaDTO `json:"objMusica"`
	Codigo    int               `json:"codigo"`
	Mensaje   string            `json:"mensaje"`
}

/**
 * @brief DTO para la respuesta de un listado de músicas.
 */
type RespuestaListaMusicaDTO struct {
	Musica  []MetadataMusicaDTO `json:"musica"`
	Codigo  int                 `json:"codigo"`
	Mensaje string              `json:"mensaje"`
}
