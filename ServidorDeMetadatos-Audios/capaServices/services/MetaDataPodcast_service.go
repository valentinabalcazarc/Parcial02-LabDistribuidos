package services

import (
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Podcast.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y la conversión entre Entity y DTO.
 */
type MetadataPodcastService struct {
	repository *repository.MetadataPodcastRepository
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataPodcast.
 * 
 * @param repository Repositorio de MetadataPodcast inyectado.
 * @return *MetadataPodcastService Instancia del servicio.
 */
func NewMetadataPodcastService(repository *repository.MetadataPodcastRepository) *MetadataPodcastService {
	return &MetadataPodcastService{repository: repository}
}

/**
 * @brief Consulta un podcast por su ID.
 * 
 * Recibe un ID, busca el Entity en el repositorio y lo convierte a
 * RespuestaMetadataPodcastDTO con el código y mensaje según el resultado.
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
		respuesta.Codigo = 400
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
