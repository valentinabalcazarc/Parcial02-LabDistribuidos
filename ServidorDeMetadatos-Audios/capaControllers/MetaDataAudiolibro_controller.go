package capacontrollers

import (
	"net/http"
	"strconv"

	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @brief Controlador para la metadata de audiolibros.
 * 
 * Expone los servicios REST de audiolibro usando Gin.
 */
type MetadataAudiolibroController struct {
	service *service.MetadataAudiolibroService
}

/**
 * @brief Constructor del controlador de MetadataAudiolibro.
 * 
 * @param service Servicio de MetadataAudiolibro que será inyectado.
 * @return *MetadataAudiolibroController Nueva instancia del controlador.
 */
func NewMetadataAudiolibroController(service *service.MetadataAudiolibroService) *MetadataAudiolibroController {
	return &MetadataAudiolibroController{service: service}
}

/**
 * @brief Consulta un audiolibro por su ID.
 * 
 * Expone el endpoint GET /audiolibro/:id.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataAudiolibroController) ConsultarAudiolibro(ctx *gin.Context) {
	idParam := ctx.Param("id")

	// Convertir el id de string a int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  400,
			"mensaje": "El ID debe ser un número entero válido",
		})
		return
	}

	respuesta := this.service.ConsultarAudiolibro(id)
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Lista todos los audiolibros disponibles.
 * 
 * Expone el endpoint GET /audiolibros.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataAudiolibroController) ListarAudiolibros(ctx *gin.Context) {
	respuesta := this.service.ListarAudiolibros()
	ctx.JSON(respuesta.Codigo, respuesta)
}
