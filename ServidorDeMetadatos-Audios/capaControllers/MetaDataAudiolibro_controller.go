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

/**
 * @brief Registra un nuevo audiolibro junto con su archivo de audio.
 * 
 * Expone el endpoint POST /audiolibro.
 * 
 * @param ctx Contexto de Gin.
 */
func (this *MetadataAudiolibroController) RegistrarAudiolibro(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")

	var audiolibroDTO dto.MetadataAudiolibroDTO
	var archivoBytes []byte
	var nombreArchivo string

	if strings.Contains(contentType, "multipart/form-data") {
		audiolibroDTO.TituloLibro = ctx.PostForm("titulo_libro")
		audiolibroDTO.Autor = ctx.PostForm("autor")
		audiolibroDTO.Narrador = ctx.PostForm("narrador")
		audiolibroDTO.Editorial = ctx.PostForm("editorial")
		if isbn, err := strconv.Atoi(ctx.PostForm("isbn")); err == nil {
			audiolibroDTO.ISBN = isbn
		}
		if cap, err := strconv.Atoi(ctx.PostForm("capitulo")); err == nil {
			audiolibroDTO.Capitulo = cap
		}

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
		if err := ctx.ShouldBindJSON(&audiolibroDTO); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"codigo":  http.StatusBadRequest,
				"mensaje": "Cuerpo JSON inválido: " + err.Error(),
			})
			return
		}
	}

	if audiolibroDTO.TituloLibro == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  http.StatusBadRequest,
			"mensaje": "El campo 'titulo_libro' es obligatorio",
		})
		return
	}

	respuesta := this.service.RegistrarAudiolibro(audiolibroDTO, archivoBytes, nombreArchivo)
	ctx.JSON(respuesta.Codigo, respuesta)
}
