/**
 * @file estadisticas_controller.go
 * @brief Controlador HTTP Gin para exponer los endpoints de consulta de estadísticas.
 */
package capacontrollers

import (
	"net/http"
	"strconv"

	capaservices "servidorEstadisticas/capaServices"

	"github.com/gin-gonic/gin"
)

/**
 * @struct EstadisticasController
 * @brief Controlador que maneja las solicitudes REST de estadísticas.
 */
type EstadisticasController struct {
	service *capaservices.EstadisticasService
}

/**
 * @brief Crea una nueva instancia de EstadisticasController.
 */
func NewEstadisticasController(service *capaservices.EstadisticasService) *EstadisticasController {
	return &EstadisticasController{
		service: service,
	}
}

/**
 * @brief Obtiene los audios reproducidos por un cliente y la cantidad de veces.
 *
 * GET /estadisticas/cliente/:cliente
 */
func (c *EstadisticasController) ConsultarReproduccionesCliente(ctx *gin.Context) {
	cliente := ctx.Param("cliente")
	respuesta := c.service.ConsultarReproduccionesCliente(cliente)
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Obtiene el ranking global de reproducciones de audios en la plataforma.
 *
 * GET /estadisticas/audios
 */
func (c *EstadisticasController) ConsultarRankingGlobal(ctx *gin.Context) {
	respuesta := c.service.ConsultarRankingGlobal()
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Obtiene los eventos recientes de reproducción.
 *
 * GET /estadisticas/historial?limite=20
 */
func (c *EstadisticasController) ConsultarHistorial(ctx *gin.Context) {
	limiteStr := ctx.DefaultQuery("limite", "50")
	limite, err := strconv.Atoi(limiteStr)
	if err != nil {
		limite = 50
	}

	respuesta := c.service.ConsultarHistorial(limite)
	ctx.JSON(http.StatusOK, respuesta)
}
