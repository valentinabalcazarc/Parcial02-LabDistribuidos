package dto

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Objeto de Transferencia de Datos para Podcast.
 * 
 * Contiene los datos que son transferidos hacia o desde las capas externas.
 */
type MetadataPodcastDTO struct {
	ID                     int              `json:"id"`
	Tipo                   entity.TipoAudio `json:"tipo"`
	Nombre                 string           `json:"nombre"`
	TituloEpisodio         string           `json:"titulo_episodio"`
	NumeroTemporada        int              `json:"numero_temporada"`
	NotasShow              string           `json:"notas_show"`
	ClasificacionContenido string           `json:"clasificacion_contenido"`
}

/**
 * @brief DTO para la respuesta de una consulta de un solo podcast.
 */
type RespuestaMetadataPodcastDTO struct {
	ObjPodcast MetadataPodcastDTO `json:"objPodcast"`
	Codigo     int                `json:"codigo"`
	Mensaje    string             `json:"mensaje"`
}

/**
 * @brief DTO para la respuesta de un listado de podcasts.
 */
type RespuestaListaPodcastDTO struct {
	Podcasts []MetadataPodcastDTO `json:"podcasts"`
	Codigo   int                  `json:"codigo"`
	Mensaje  string               `json:"mensaje"`
}
