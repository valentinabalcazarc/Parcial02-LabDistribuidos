package services

import (
	"fmt"
	"servidorMetadatos/capaAccesoADatos/entity"
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Música.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio, la conversión entre Entity y DTO, y
 * coordinando internamente con AudioService para el almacenamiento físico.
 */
type MetadataMusicaService struct {
	repository   *repository.MetadataMusicaRepository
	audioService *AudioService
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataMusica.
 * 
 * @param repository Repositorio de MetadataMusica inyectado.
 * @param audioService Servicio centralizado de almacenamiento de archivos de audio.
 * @return *MetadataMusicaService Instancia del servicio.
 */
func NewMetadataMusicaService(repository *repository.MetadataMusicaRepository, audioService *AudioService) *MetadataMusicaService {
	return &MetadataMusicaService{
		repository:   repository,
		audioService: audioService,
	}
}

/**
 * @brief Consulta una música por su ID.
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
		respuesta.Codigo = 404
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

/**
 * @brief Registra una nueva canción (metadatos) y coordina internamente el guardado del archivo físico.
 * 
 * Si se incluye el archivo binario, llama por dentro a AudioService para persistirlo en disco,
 * evitando duplicar la lógica de I/O en los controladores.
 * 
 * @param musicaDTO DTO con los datos de la canción.
 * @param archivoBytes Contenido binario del archivo de audio (opcional).
 * @param nombreArchivo Nombre del archivo para guardarlo en disco (ej. "cancion1.mp3").
 * @return dto.RespuestaMetadataMusicaDTO DTO con la respuesta y entidad creada.
 */
func (this *MetadataMusicaService) RegistrarMusica(musicaDTO dto.MetadataMusicaDTO, archivoBytes []byte, nombreArchivo string) dto.RespuestaMetadataMusicaDTO {
	// 1. Guardar archivo físico por dentro si viene adjunto
	if len(archivoBytes) > 0 && this.audioService != nil {
		if nombreArchivo == "" {
			nombreArchivo = fmt.Sprintf("%s.mp3", musicaDTO.TituloCancion)
		}
		err := this.audioService.GuardarArchivo(nombreArchivo, archivoBytes)
		if err != nil {
			return dto.RespuestaMetadataMusicaDTO{
				Codigo:  500,
				Mensaje: "Error al guardar el archivo físico de audio: " + err.Error(),
			}
		}
	}

	// 2. Mapear DTO a Entity
	var entidad entity.MetaDataMusica
	entidad.SetId(musicaDTO.ID)
	if musicaDTO.Tipo.GetId() > 0 {
		entidad.SetTipo(musicaDTO.Tipo)
	} else {
		entidad.SetTipo(entity.NewTipoAudio(1, "Música"))
	}
	entidad.SetArtistaPrincipal(musicaDTO.ArtistaPrincipal)
	entidad.SetAlbum(musicaDTO.Album)
	entidad.SetGenero(musicaDTO.Genero)
	entidad.SetTituloCancion(musicaDTO.TituloCancion)
	entidad.SetSelloDiscografico(musicaDTO.SelloDiscografico)
	entidad.SetAnioLanzamiento(musicaDTO.AnioLanzamiento)

	// 3. Guardar en repositorio de metadatos
	guardada := this.repository.GuardarMusica(entidad)

	musicaDTO.ID = guardada.GetId()
	musicaDTO.Tipo = guardada.GetTipo()

	return dto.RespuestaMetadataMusicaDTO{
		ObjMusica: musicaDTO,
		Codigo:    201,
		Mensaje:   "Música y archivo registrados exitosamente",
	}
}
