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

func (c *MetadataMusicaController) RegistrarRutas(r gin.IRouter) {
	r.GET("/metadatosaudios/musica/:id", c.ConsultarMusica)
	r.GET("/metadatosaudios/musicas", c.ListarMusica)
	r.POST("/metadatosaudios/musica", c.RegistrarMusica)
}

/**
 * @brief Consulta una música por su ID.
 *
 * Expone el endpoint GET /metadatosaudios/musica/:id.
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
 * Expone el endpoint GET /metadatosaudios/musicas.
 *
 * @param ctx Contexto de Gin que contiene la información de la petición HTTP.
 */
func (this *MetadataMusicaController) ListarMusica(ctx *gin.Context) {
	respuesta := this.service.ListarMusica()
	ctx.JSON(respuesta.Codigo, respuesta)
}

/**
 * @brief Registra una nueva música junto con su archivo de audio.
 *
 * Expone el endpoint POST /metadatosaudios/musica.
 * Soporta tanto peticiones JSON puras como multipart/form-data (con archivo binario).
 *
 * @param ctx Contexto de Gin.
 */
func (this *MetadataMusicaController) RegistrarMusica(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")

	var musicaDTO dto.MetadataMusicaDTO
	var archivoBytes []byte
	var nombreArchivo string

	if strings.Contains(contentType, "multipart/form-data") {
		// 1. Extraer campos del formulario
		musicaDTO.TituloCancion = ctx.PostForm("titulo_cancion")
		musicaDTO.ArtistaPrincipal = ctx.PostForm("artista_principal")
		musicaDTO.Album = ctx.PostForm("album")
		musicaDTO.Genero = ctx.PostForm("genero")
		musicaDTO.SelloDiscografico = ctx.PostForm("sello_discografico")
		musicaDTO.AnioLanzamiento = ctx.PostForm("anio_lanzamiento")

		// 2. Extraer archivo de audio si se adjuntó
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
		// Petición JSON
		if err := ctx.ShouldBindJSON(&musicaDTO); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"codigo":  http.StatusBadRequest,
				"mensaje": "Cuerpo JSON inválido: " + err.Error(),
			})
			return
		}
	}

	if musicaDTO.TituloCancion == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  http.StatusBadRequest,
			"mensaje": "El campo 'titulo_cancion' es obligatorio",
		})
		return
	}

	respuesta := this.service.RegistrarMusica(musicaDTO, archivoBytes, nombreArchivo)
	ctx.JSON(respuesta.Codigo, respuesta)
}
