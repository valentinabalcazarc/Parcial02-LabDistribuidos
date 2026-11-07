package capacontrollers

import (
	"net/http"
	"strconv"

	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @brief Controlador para la metadata de Musica.
 * 
 * Expone los servicios REST de Musica usando Gin.
 */
type MetadataMusicaController struct {
	service *service.MetadataMusicaService
}

/**
 * @brief Constructor del controlador de MetadataMusica.
 * 
 * @param service Servicio de MetadataMusica que será inyectado.
 * @return *MetadataMusicaController Nueva instancia del controlador.
 */
func NewMetadataMusicaController(service *service.MetadataMusicaService) *MetadataMusicaController {
	return &MetadataMusicaController{service: service}
}

/**
 * @brief Consulta una música por su ID.
 * 
 * Expone el endpoint GET /musica/:id.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataMusicaController) ConsultarMusica(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  400,
			"mensaje": "El ID debe ser un número entero válido",
		})
		return
	}

	respuesta := this.service.ConsultarMusica(id)
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Lista todas las músicas disponibles.
 * 
 * Expone el endpoint GET /musicas.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataMusicaController) ListarMusica(ctx *gin.Context) {
	respuesta := this.service.ListarMusica()
	ctx.JSON(respuesta.Codigo, respuesta)
}
