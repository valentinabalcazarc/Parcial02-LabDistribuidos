package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pasarela-pagos/services"
	"pasarela-pagos/services/dto"
)

type PagoController struct {
	servicio *services.PagoService
}

func NuevoPagoController(servicio *services.PagoService) *PagoController {
	return &PagoController{servicio: servicio}
}

func (ctl *PagoController) Registrar(rg *gin.RouterGroup) {
	rg.POST("/tokens", ctl.generarToken)
	rg.POST("/pagos", ctl.realizarPago)
	rg.GET("/transacciones/:cedula", ctl.consultarTransacciones)
}

// POST /api/tokens — Servicio REST 1: generar token de idempotencia.
func (ctl *PagoController) generarToken(c *gin.Context) {
	var solicitud dto.SolicitudTokenDTO
	if err := c.ShouldBindJSON(&solicitud); err != nil {
		responderCuerpoInvalido(c)
		return
	}
	respuesta, err := ctl.servicio.GenerarToken(solicitud)
	if err != nil {
		responderError(c, err)
		return
	}
	c.JSON(http.StatusCreated, respuesta)
}

// POST /api/pagos — Servicio REST 2: realizar pago.
func (ctl *PagoController) realizarPago(c *gin.Context) {
	var sol dto.SolicitudPagoDTO
	if err := c.ShouldBindJSON(&sol); err != nil {
		responderCuerpoInvalido(c)
		return
	}
	respuesta, err := ctl.servicio.RealizarPago(sol)
	if err != nil {
		responderError(c, err)
		return
	}
	c.JSON(http.StatusOK, respuesta)
}

// GET /api/transacciones/:cedula — Servicio REST 3: consultar transacciones.
func (ctl *PagoController) consultarTransacciones(c *gin.Context) {
	respuesta, err := ctl.servicio.ConsultarTransacciones(c.Param("cedula"))
	if err != nil {
		responderError(c, err)
		return
	}
	c.JSON(http.StatusOK, respuesta)
}
