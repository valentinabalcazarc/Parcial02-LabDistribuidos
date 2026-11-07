package services

import (
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Música.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y la conversión entre Entity y DTO.
 */
type MetadataMusicaService struct {
	repository *repository.MetadataMusicaRepository
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataMusica.
 * 
 * @param repository Repositorio de MetadataMusica inyectado.
 * @return *MetadataMusicaService Instancia del servicio.
 */
func NewMetadataMusicaService(repository *repository.MetadataMusicaRepository) *MetadataMusicaService {
	return &MetadataMusicaService{repository: repository}
}

/**
 * @brief Consulta una música por su ID.
 * 
 * Recibe un ID, busca el Entity en el repositorio y lo convierte a
 * RespuestaMetadataMusicaDTO con el código y mensaje según el resultado.
 * 
 * @param id Identificador de la música a buscar.
 * @return dto.RespuestaMetadataMusicaDTO DTO con la respuesta de la consulta.
 */
func (this *MetadataMusicaService) ConsultarMusica(id int) dto.RespuestaMetadataMusicaDTO {
	var respuesta dto.RespuestaMetadataMusicaDTO

	musica, encontrado := this.repository.BuscarMusica(id)

	if encontrado {
		var musicaDTO dto.MetadataMusicaDTO
		musicaDTO.ID = musica.GetId()
		musicaDTO.Tipo = musica.GetTipo()
		musicaDTO.ArtistaPrincipal = musica.GetArtistaPrincipal()
		musicaDTO.Album = musica.GetAlbum()
		musicaDTO.Genero = musica.GetGenero()
		musicaDTO.TituloCancion = musica.GetTituloCancion()
		musicaDTO.SelloDiscografico = musica.GetSelloDiscografico()
		musicaDTO.AnioLanzamiento = musica.GetAnioLanzamiento()

		respuesta.ObjMusica = musicaDTO
		respuesta.Codigo = 200
		respuesta.Mensaje = "Métadata de la musica encontrada"
	} else {
		respuesta.Codigo = 400
		respuesta.Mensaje = "La métadata de la musica no se encontró"
	}
	return respuesta
}

/**
 * @brief Obtiene la lista de todas las músicas.
 * 
 * @return dto.RespuestaListaMusicaDTO DTO con la lista y código de respuesta.
 */
func (this *MetadataMusicaService) ListarMusica() dto.RespuestaListaMusicaDTO {
	lista := this.repository.ListarMusica()
	var dtos []dto.MetadataMusicaDTO

	for _, musica := range lista {
		var d dto.MetadataMusicaDTO
		d.ID = musica.GetId()
		d.Tipo = musica.GetTipo()
		d.ArtistaPrincipal = musica.GetArtistaPrincipal()
		d.Album = musica.GetAlbum()
		d.Genero = musica.GetGenero()
		d.TituloCancion = musica.GetTituloCancion()
		d.SelloDiscografico = musica.GetSelloDiscografico()
		d.AnioLanzamiento = musica.GetAnioLanzamiento()

		dtos = append(dtos, d)
	}

	return dto.RespuestaListaMusicaDTO{
		Musica:  dtos,
		Codigo:  200,
		Mensaje: "Lista de música obtenida correctamente",
	}
}
