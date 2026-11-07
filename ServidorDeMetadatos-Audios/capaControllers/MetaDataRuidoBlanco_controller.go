package capacontrollers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"servidorMetadatos/capaServices/dto"
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

/**
 * @brief Registra un nuevo ruido blanco junto con su archivo de audio.
 * 
 * Expone el endpoint POST /ruidoblanco.
 * 
 * @param ctx Contexto de Gin.
 */
func (this *MetadataRuidoBlancoController) RegistrarRuidoBlanco(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")

	var ruidoDTO dto.MetadataRuidoBlancoDTO
	var archivoBytes []byte
	var nombreArchivo string

	if strings.Contains(contentType, "multipart/form-data") {
		ruidoDTO.TipoSonido = ctx.PostForm("tipo_sonido")
		ruidoDTO.FuenteAudio = ctx.PostForm("fuente_audio")
		ruidoDTO.UsoSugerido = ctx.PostForm("uso_sugerido")
		ruidoDTO.Proveedor = ctx.PostForm("proveedor")
		if dur, err := strconv.Atoi(ctx.PostForm("duracion")); err == nil {
			ruidoDTO.Duracion = dur
		}
		ruidoDTO.FrecuenciaDominante = ctx.PostForm("frecuencia_dominante")

		fileHeader, err := ctx.FormFile("archivo")
		if err == nil && fileHeader != nil {
			file, err := fileHeader.Open()
			if err == nil {
				defer file.Close()
				archivoBytes, _ = io.ReadAll(file)
				nombreArchivo = fileHeader.Filename
			}
		}

		if customNombre := ctx.PostForm("nombre_archivo"); customNombre != "" {
			nombreArchivo = customNombre
		}
	} else {
		if err := ctx.ShouldBindJSON(&ruidoDTO); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"codigo":  http.StatusBadRequest,
				"mensaje": "Cuerpo JSON inválido: " + err.Error(),
			})
			return
		}
	}

	if ruidoDTO.TipoSonido == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  http.StatusBadRequest,
			"mensaje": "El campo 'tipo_sonido' es obligatorio",
		})
		return
	}

	respuesta := this.service.RegistrarRuidoBlanco(ruidoDTO, archivoBytes, nombreArchivo)
	ctx.JSON(respuesta.Codigo, respuesta)
}
