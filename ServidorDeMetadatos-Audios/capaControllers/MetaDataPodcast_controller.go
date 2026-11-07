package capacontrollers

import (
	"net/http"
	"strconv"

	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @brief Controlador para la metadata de Podcast.
 * 
 * Expone los servicios REST de Podcast usando Gin.
 */
type MetadataPodcastController struct {
	service *service.MetadataPodcastService
}

/**
 * @brief Constructor del controlador de MetadataPodcast.
 * 
 * @param service Servicio de MetadataPodcast que será inyectado.
 * @return *MetadataPodcastController Nueva instancia del controlador.
 */
func NewMetadataPodcastController(service *service.MetadataPodcastService) *MetadataPodcastController {
	return &MetadataPodcastController{service: service}
}

/**
 * @brief Consulta un podcast por su ID.
 * 
 * Expone el endpoint GET /podcast/:id.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataPodcastController) ConsultarPodcast(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  400,
			"mensaje": "El ID debe ser un número entero válido",
		})
		return
	}

	respuesta := this.service.ConsultarPodcast(id)
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Lista todos los podcasts disponibles.
 * 
 * Expone el endpoint GET /podcasts.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataPodcastController) ListarPodcasts(ctx *gin.Context) {
	respuesta := this.service.ListarPodcasts()
	ctx.JSON(respuesta.Codigo, respuesta)
}
