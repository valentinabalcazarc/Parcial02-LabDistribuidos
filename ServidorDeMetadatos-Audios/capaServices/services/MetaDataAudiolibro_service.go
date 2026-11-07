package services

import (
	"fmt"
	"servidorMetadatos/capaAccesoADatos/entity"
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Audiolibro.
 *
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y coordinando internamente con AudioService
 * el guardado del archivo físico.
 */
type MetadataAudiolibroService struct {
	repository   *repository.MetadataAudiolibroRepository
	audioService *AudioService
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataAudiolibro.
 *
 * @param repository Repositorio de MetadataAudiolibro inyectado.
 * @param audioService Servicio centralizado de almacenamiento de archivos de audio.
 * @return *MetadataAudiolibroService Instancia del servicio.
 */
func NewMetadataAudiolibroService(repository *repository.MetadataAudiolibroRepository, audioService *AudioService) *MetadataAudiolibroService {
	return &MetadataAudiolibroService{
		repository:   repository,
		audioService: audioService,
	}
}

/**
 * @brief Consulta un audiolibro por su ID.
 *
 * @param id Identificador del audiolibro a buscar.
 * @return dto.RespuestaMetadataAudiolibroDTO DTO con la respuesta de la consulta.
 */
func (this *MetadataAudiolibroService) ConsultarAudiolibro(id int) dto.RespuestaMetadataAudiolibroDTO {
	var respuesta dto.RespuestaMetadataAudiolibroDTO

	audiolibro, encontrado := this.repository.BuscarAudiolibro(id)

	if encontrado {
		var audiolibroDTO dto.MetadataAudiolibroDTO
		audiolibroDTO.ID = audiolibro.GetId()
		audiolibroDTO.Tipo = audiolibro.GetTipo()
		audiolibroDTO.TituloLibro = audiolibro.GetTituloLibro()
		audiolibroDTO.Autor = audiolibro.GetAutor()
		audiolibroDTO.Narrador = audiolibro.GetNarrador()
		audiolibroDTO.Editorial = audiolibro.GetEditorial()
		audiolibroDTO.ISBN = audiolibro.GetISBN()
		audiolibroDTO.Capitulo = audiolibro.GetCapitulo()

		respuesta.ObjAudiolibro = audiolibroDTO
		respuesta.Codigo = 200
		respuesta.Mensaje = "Métadata del audiolibro encontrada"
	} else {
		respuesta.Codigo = 404
		respuesta.Mensaje = "La métadata del audiolibro no se encontró"
	}
	return respuesta
}

/**
 * @brief Obtiene la lista de todos los audiolibros.
 *
 * @return dto.RespuestaListaAudiolibroDTO DTO con la lista y código de respuesta.
 */
func (this *MetadataAudiolibroService) ListarAudiolibros() dto.RespuestaListaAudiolibroDTO {
	lista := this.repository.ListarAudiolibros()
	var dtos []dto.MetadataAudiolibroDTO

	for _, audiolibro := range lista {
		var d dto.MetadataAudiolibroDTO
		d.ID = audiolibro.GetId()
		d.Tipo = audiolibro.GetTipo()
		d.TituloLibro = audiolibro.GetTituloLibro()
		d.Autor = audiolibro.GetAutor()
		d.Narrador = audiolibro.GetNarrador()
		d.Editorial = audiolibro.GetEditorial()
		d.ISBN = audiolibro.GetISBN()
		d.Capitulo = audiolibro.GetCapitulo()

		dtos = append(dtos, d)
	}

	return dto.RespuestaListaAudiolibroDTO{
		Audiolibros: dtos,
		Codigo:      200,
		Mensaje:     "Lista de audiolibros obtenida correctamente",
	}
}

/**
 * @brief Registra un nuevo audiolibro y coordina por dentro el guardado del archivo físico.
 *
 * @param audioDTO DTO con los datos del audiolibro.
 * @param archivoBytes Contenido binario del archivo de audio.
 * @param nombreArchivo Nombre del archivo en disco.
 * @return dto.RespuestaMetadataAudiolibroDTO DTO con la respuesta y entidad creada.
 */
func (this *MetadataAudiolibroService) RegistrarAudiolibro(audioDTO dto.MetadataAudiolibroDTO, archivoBytes []byte, nombreArchivo string) dto.RespuestaMetadataAudiolibroDTO {
	if len(archivoBytes) > 0 && this.audioService != nil {
		if nombreArchivo == "" {
			nombreArchivo = fmt.Sprintf("%s.mp3", audioDTO.TituloLibro)
		}
		err := this.audioService.GuardarArchivo(nombreArchivo, archivoBytes)
		if err != nil {
			return dto.RespuestaMetadataAudiolibroDTO{
				Codigo:  500,
				Mensaje: "Error al guardar el archivo físico de audio: " + err.Error(),
			}
		}
	}

	var entidad entity.MetaDataAudiolibro
	entidad.SetId(audioDTO.ID)
	if audioDTO.Tipo.GetId() > 0 {
		entidad.SetTipo(audioDTO.Tipo)
	} else {
		entidad.SetTipo(entity.NewTipoAudio(4, "Audiolibro"))
	}
	entidad.SetTituloLibro(audioDTO.TituloLibro)
	entidad.SetAutor(audioDTO.Autor)
	entidad.SetNarrador(audioDTO.Narrador)
	entidad.SetEditorial(audioDTO.Editorial)
	entidad.SetISBN(audioDTO.ISBN)
	entidad.SetCapitulo(audioDTO.Capitulo)

	guardado := this.repository.GuardarAudiolibro(entidad)

	audioDTO.ID = guardado.GetId()
	audioDTO.Tipo = guardado.GetTipo()

	return dto.RespuestaMetadataAudiolibroDTO{
		ObjAudiolibro: audioDTO,
		Codigo:        201,
		Mensaje:       "Audiolibro y archivo registrados exitosamente",
	}
}
