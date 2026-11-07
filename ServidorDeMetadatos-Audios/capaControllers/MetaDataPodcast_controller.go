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

/**
 * @brief Registra un nuevo podcast junto con su archivo de audio.
 * 
 * Expone el endpoint POST /podcast.
 * 
 * @param ctx Contexto de Gin.
 */
func (this *MetadataPodcastController) RegistrarPodcast(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")

	var podcastDTO dto.MetadataPodcastDTO
	var archivoBytes []byte
	var nombreArchivo string

	if strings.Contains(contentType, "multipart/form-data") {
		podcastDTO.Nombre = ctx.PostForm("nombre")
		podcastDTO.TituloEpisodio = ctx.PostForm("titulo_episodio")
		if temp, err := strconv.Atoi(ctx.PostForm("numero_temporada")); err == nil {
			podcastDTO.NumeroTemporada = temp
		}
		podcastDTO.NotasShow = ctx.PostForm("notas_show")
		podcastDTO.ClasificacionContenido = ctx.PostForm("clasificacion_contenido")

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
		if err := ctx.ShouldBindJSON(&podcastDTO); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"codigo":  http.StatusBadRequest,
				"mensaje": "Cuerpo JSON inválido: " + err.Error(),
			})
			return
		}
	}

	if podcastDTO.TituloEpisodio == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  http.StatusBadRequest,
			"mensaje": "El campo 'titulo_episodio' es obligatorio",
		})
		return
	}

	respuesta := this.service.RegistrarPodcast(podcastDTO, archivoBytes, nombreArchivo)
	ctx.JSON(respuesta.Codigo, respuesta)
}
