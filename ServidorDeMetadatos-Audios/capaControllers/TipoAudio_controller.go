package capacontrollers

import (
	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @brief Controlador para los tipos de audio.
 * 
 * Gestiona las peticiones HTTP relacionadas con la entidad TipoAudio.
 */
type TipoAudioController struct {
	service *service.TipoAudioService
}

/**
 * @brief Constructor del controlador de TipoAudio.
 * 
 * @param service Servicio de TipoAudio que será inyectado.
 * @return *TipoAudioController Nueva instancia del controlador.
 */
func NewTipoAudioController(service *service.TipoAudioService) *TipoAudioController {
	return &TipoAudioController{service: service}
}

/**
 * @brief Lista todos los tipos de audio.
 * 
 * Expone el endpoint GET /tipos-audio.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición.
 */
func (c *TipoAudioController) ListarTipos(ctx *gin.Context) {
	respuesta := c.service.ListarTipos()
	ctx.JSON(respuesta.Codigo, respuesta)
}
