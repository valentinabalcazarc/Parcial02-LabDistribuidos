/**
 * @file Audio_controller.go
 * @brief Definición del controlador para manejar las peticiones HTTP de audios usando Gin.
 */
package capacontrollers

import (
	"io"
	"net/http"
	dtos "servidorMetadatos/capaServices/dto"
	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

/**
 * @struct AudioController
 * @brief Controlador que maneja las solicitudes HTTP relacionadas con el almacenamiento de audios.
 */
type AudioController struct {
	service *service.AudioService
}

/**
 * @brief Crea una nueva instancia de AudioController.
 *
 * @param service Puntero a la instancia de AudioService que utilizará el controlador.
 * @return *AudioController Puntero al nuevo controlador instanciado.
 */
func NewAudioController(service *service.AudioService) *AudioController {
	return &AudioController{service: service}
}

/**
 * @brief Manejador Gin para guardar un archivo de audio físico.
 *
 * Este método recibe una petición POST multipart/form-data con el archivo de audio
 * y su nombre opcional, delegando la tarea de guardado físico al servicio de audio.
 *
 * @param ctx Contexto de Gin.
 */
func (thisC *AudioController) GuardarAudio(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("archivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"codigo":  http.StatusBadRequest,
			"mensaje": "Error leyendo el archivo del formulario (campo 'archivo'): " + err.Error(),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"codigo":  http.StatusInternalServerError,
			"mensaje": "Error al abrir el archivo recibido: " + err.Error(),
		})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"codigo":  http.StatusInternalServerError,
			"mensaje": "Error al leer los bytes del archivo: " + err.Error(),
		})
		return
	}

	nombreArchivo := ctx.DefaultPostForm("nombre_archivo", fileHeader.Filename)
	dto := dtos.AudioDTO{
		NombreArchivo: nombreArchivo,
	}

	err = thisC.service.GuardarAudio(dto, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"codigo":  http.StatusInternalServerError,
			"mensaje": "Error al guardar el archivo en disco: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"codigo":         http.StatusOK,
		"mensaje":        "Audio guardado exitosamente",
		"nombre_archivo": nombreArchivo,
	})
}
