/**
 * @file Audio_controller.go
 * @brief Definición del controlador para manejar las peticiones HTTP de audios.
 */
package capacontrollers

import (
	"fmt"
	"io"
	"net/http"
	dtos "servidorMetadatos/capaServices/dto"
	service "servidorMetadatos/capaServices/services"
)

/**
 * @struct AudioController
 * @brief Controlador que maneja las solicitudes HTTP relacionadas con los audios.
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
 * @brief Manejador HTTP para guardar un archivo de audio.
 *
 * Este método recibe una petición POST multipart/form-data con el archivo de audio
 * y su nombre, y delega la tarea de guardado al servicio.
 *
 * @param w http.ResponseWriter para enviar la respuesta HTTP.
 * @param r *http.Request petición HTTP recibida.
 */
func (thisC *AudioController) GuardarAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(50 << 20)
	file, _, err := r.FormFile("Archivo")
	if err != nil {
		http.Error(w, "Error leyendo el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	dto := dtos.AudioDTO{
		NombreArchivo: r.FormValue("nombre_archivo"),
	}

	err = thisC.service.GuardarAudio(dto, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al guardar: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Audio guardado exitosamente"))
}

/**
 * @brief Manejador HTTP para obtener un archivo de audio.
 *
 * Este método recibe una petición GET con el nombre del archivo de audio como parámetro de consulta
 * y lo devuelve como respuesta HTTP.
 *
 * @param w http.ResponseWriter para enviar la respuesta HTTP (el archivo de audio).
 * @param r *http.Request petición HTTP recibida.
 */
func (thisC *AudioController) ObtenerAudio(w http.ResponseWriter, r *http.Request) {
	nombreArchivo := r.URL.Query().Get("nombre")
	if nombreArchivo == "" {
		http.Error(w, "Parámetro 'nombre' requerido", http.StatusBadRequest)
		return
	}

	bytesAudio, err := thisC.service.ObtenerAudio(nombreArchivo)
	if err != nil {
		http.Error(w, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(bytesAudio)
}
