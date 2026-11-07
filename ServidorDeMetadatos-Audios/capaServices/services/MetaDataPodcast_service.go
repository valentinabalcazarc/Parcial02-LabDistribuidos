package services

import (
	"fmt"
	"servidorMetadatos/capaAccesoADatos/entity"
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Podcast.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y coordinando internamente con AudioService
 * el guardado del archivo físico.
 */
type MetadataPodcastService struct {
	repository   *repository.MetadataPodcastRepository
	audioService *AudioService
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataPodcast.
 * 
 * @param repository Repositorio de MetadataPodcast inyectado.
 * @param audioService Servicio centralizado de almacenamiento de archivos de audio.
 * @return *MetadataPodcastService Instancia del servicio.
 */
func NewMetadataPodcastService(repository *repository.MetadataPodcastRepository, audioService *AudioService) *MetadataPodcastService {
	return &MetadataPodcastService{
		repository:   repository,
		audioService: audioService,
	}
}

/**
 * @brief Consulta un podcast por su ID.
 * 
 * @param id Identificador del podcast a buscar.
 * @return dto.RespuestaMetadataPodcastDTO DTO con la respuesta de la consulta.
 */
func (this *MetadataPodcastService) ConsultarPodcast(id int) dto.RespuestaMetadataPodcastDTO {
	var respuesta dto.RespuestaMetadataPodcastDTO

	podcast, encontrado := this.repository.BuscarPodcast(id)

	if encontrado {
		var podcastDTO dto.MetadataPodcastDTO
		podcastDTO.ID = podcast.GetId()
		podcastDTO.Tipo = podcast.GetTipo()
		podcastDTO.Nombre = podcast.GetNombre()
		podcastDTO.TituloEpisodio = podcast.GetTituloEpisodio()
		podcastDTO.NumeroTemporada = podcast.GetNumeroTemporada()
		podcastDTO.NotasShow = podcast.GetNotasShow()
		podcastDTO.ClasificacionContenido = podcast.GetClasificacionContenido()

		respuesta.ObjPodcast = podcastDTO
		respuesta.Codigo = 200
		respuesta.Mensaje = "Métadata del podcast encontrada"
	} else {
		respuesta.Codigo = 404
		respuesta.Mensaje = "La métadata del podcast no se encontró"
	}
	return respuesta
}

/**
 * @brief Obtiene la lista de todos los podcasts.
 * 
 * @return dto.RespuestaListaPodcastDTO DTO con la lista y código de respuesta.
 */
func (this *MetadataPodcastService) ListarPodcasts() dto.RespuestaListaPodcastDTO {
	lista := this.repository.ListarPodcasts()
	var dtos []dto.MetadataPodcastDTO

	for _, podcast := range lista {
		var d dto.MetadataPodcastDTO
		d.ID = podcast.GetId()
		d.Tipo = podcast.GetTipo()
		d.Nombre = podcast.GetNombre()
		d.TituloEpisodio = podcast.GetTituloEpisodio()
		d.NumeroTemporada = podcast.GetNumeroTemporada()
		d.NotasShow = podcast.GetNotasShow()
		d.ClasificacionContenido = podcast.GetClasificacionContenido()

		dtos = append(dtos, d)
	}

	return dto.RespuestaListaPodcastDTO{
		Podcasts: dtos,
		Codigo:   200,
		Mensaje:  "Lista de podcasts obtenida correctamente",
	}
}

/**
 * @brief Registra un nuevo podcast y coordina internamente el guardado del archivo físico.
 * 
 * @param podcastDTO DTO con los datos del podcast.
 * @param archivoBytes Contenido binario del archivo de audio (opcional).
 * @param nombreArchivo Nombre del archivo en disco.
 * @return dto.RespuestaMetadataPodcastDTO DTO con la respuesta y entidad creada.
 */
func (this *MetadataPodcastService) RegistrarPodcast(podcastDTO dto.MetadataPodcastDTO, archivoBytes []byte, nombreArchivo string) dto.RespuestaMetadataPodcastDTO {
	if len(archivoBytes) > 0 && this.audioService != nil {
		if nombreArchivo == "" {
			nombreArchivo = fmt.Sprintf("%s.mp3", podcastDTO.TituloEpisodio)
		}
		err := this.audioService.GuardarArchivo(nombreArchivo, archivoBytes)
		if err != nil {
			return dto.RespuestaMetadataPodcastDTO{
				Codigo:  500,
				Mensaje: "Error al guardar el archivo físico de audio: " + err.Error(),
			}
		}
	}

	var entidad entity.MetaDataPodcast
	entidad.SetId(podcastDTO.ID)
	if podcastDTO.Tipo.GetId() > 0 {
		entidad.SetTipo(podcastDTO.Tipo)
	} else {
		entidad.SetTipo(entity.NewTipoAudio(3, "Podcast"))
	}
	entidad.SetNombre(podcastDTO.Nombre)
	entidad.SetTituloEpisodio(podcastDTO.TituloEpisodio)
	entidad.SetNumeroTemporada(podcastDTO.NumeroTemporada)
	entidad.SetNotasShow(podcastDTO.NotasShow)
	entidad.SetClasificacionContenido(podcastDTO.ClasificacionContenido)

	guardado := this.repository.GuardarPodcast(entidad)

	podcastDTO.ID = guardado.GetId()
	podcastDTO.Tipo = guardado.GetTipo()

	return dto.RespuestaMetadataPodcastDTO{
		ObjPodcast: podcastDTO,
		Codigo:     201,
		Mensaje:    "Podcast y archivo registrados exitosamente",
	}
}
