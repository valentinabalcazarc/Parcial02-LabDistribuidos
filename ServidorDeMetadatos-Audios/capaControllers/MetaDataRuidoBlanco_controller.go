package capacontrollers

import (
	"net/http"
	"strconv"

	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @brief Controlador para la metadata de RuidoBlanco.
 * 
 * Expone los servicios REST de RuidoBlanco usando Gin.
 */
type MetadataRuidoBlancoController struct {
	service *service.MetadataRuidoBlancoService
}

/**
 * @brief Constructor del controlador de MetadataRuidoBlanco.
 * 
 * @param service Servicio de MetadataRuidoBlanco que será inyectado.
 * @return *MetadataRuidoBlancoController Nueva instancia del controlador.
 */
func NewMetadataRuidoBlancoController(service *service.MetadataRuidoBlancoService) *MetadataRuidoBlancoController {
	return &MetadataRuidoBlancoController{service: service}
}

/**
 * @brief Consulta un ruido blanco por su ID.
 * 
 * Expone el endpoint GET /ruidoblanco/:id.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataRuidoBlancoController) ConsultarRuidoBlanco(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  400,
			"mensaje": "El ID debe ser un número entero válido",
		})
		return
	}

	respuesta := this.service.ConsultarRuidoBlanco(id)
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Lista todos los ruidos blancos disponibles.
 * 
 * Expone el endpoint GET /ruidosblancos.
 * 
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataRuidoBlancoController) ListarRuidoBlanco(ctx *gin.Context) {
	respuesta := this.service.ListarRuidoBlanco()
	ctx.JSON(respuesta.Codigo, respuesta)
}
