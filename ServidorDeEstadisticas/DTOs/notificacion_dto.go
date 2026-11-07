/**
 * @file notificacion_dto.go
 * @brief DTOs para el procesamiento y consulta de estadísticas de reproducción.
 */
package dtos

/**
 * @struct NotificacionReproduccion
 * @brief DTO recibido a través de RabbitMQ cuando se reproduce un audio.
 */
type NotificacionReproduccion struct {
	Cliente   string `json:"cliente"`               /**< @brief Identificador del cliente (usuario o IP). */
	Titulo    string `json:"titulo"`                /**< @brief Título del audio reproducido. */
	TipoAudio string `json:"tipo_audio,omitempty"`  /**< @brief Categoría (Música, Podcast, etc.). */
	FechaHora string `json:"fecha_hora,omitempty"`  /**< @brief Fecha y hora del evento. */
	Mensaje   string `json:"mensaje"`               /**< @brief Mensaje descriptivo. */
}

/**
 * @struct EstadisticaAudioDTO
 * @brief Resumen de reproducciones de un audio para un cliente específico (Requerimiento d).
 */
type EstadisticaAudioDTO struct {
	Titulo             string `json:"titulo"`
	TipoAudio          string `json:"tipo_audio,omitempty"`
	CantidadVeces      int    `json:"cantidad_veces"`
	UltimaReproduccion string `json:"ultima_reproduccion"`
}

/**
 * @struct ResumenClienteDTO
 * @brief DTO con la lista de audios reproducidos por un cliente y sus conteos.
 */
type ResumenClienteDTO struct {
	Cliente             string                `json:"cliente"`
	TotalReproducciones int                   `json:"total_reproducciones"`
	Audios              []EstadisticaAudioDTO `json:"audios"`
}

/**
 * @struct RankingAudioDTO
 * @brief Estadística global de reproducciones de un audio en la plataforma.
 */
type RankingAudioDTO struct {
	Titulo              string `json:"titulo"`
	TotalReproducciones int    `json:"total_reproducciones"`
	UltimaReproduccion  string `json:"ultima_reproduccion"`
}

/**
 * @struct RespuestaAPI
 * @brief Estructura genérica de respuesta HTTP para la API de estadísticas.
 */
type RespuestaAPI struct {
	Codigo  int         `json:"codigo"`
	Mensaje string      `json:"mensaje"`
	Datos   interface{} `json:"datos,omitempty"`
}
