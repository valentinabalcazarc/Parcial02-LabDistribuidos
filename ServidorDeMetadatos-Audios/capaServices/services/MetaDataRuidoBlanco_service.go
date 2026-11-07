package services

import (
	"fmt"
	"servidorMetadatos/capaAccesoADatos/entity"
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Ruido Blanco.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y coordinando internamente con AudioService
 * el guardado del archivo físico.
 */
type MetadataRuidoBlancoService struct {
	repository   *repository.MetadataRuidoBlancoRepository
	audioService *AudioService
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataRuidoBlanco.
 * 
 * @param repository Repositorio de MetadataRuidoBlanco inyectado.
 * @param audioService Servicio centralizado de almacenamiento de archivos de audio.
 * @return *MetadataRuidoBlancoService Instancia del servicio.
 */
func NewMetadataRuidoBlancoService(repository *repository.MetadataRuidoBlancoRepository, audioService *AudioService) *MetadataRuidoBlancoService {
	return &MetadataRuidoBlancoService{
		repository:   repository,
		audioService: audioService,
	}
}

/**
 * @brief Consulta un registro de ruido blanco por su ID.
 * 
 * @param id Identificador del ruido blanco a buscar.
 * @return dto.RespuestaMetadataRuidoBlancoDTO DTO con la respuesta de la consulta.
 */
func (this *MetadataRuidoBlancoService) ConsultarRuidoBlanco(id int) dto.RespuestaMetadataRuidoBlancoDTO {
	var respuesta dto.RespuestaMetadataRuidoBlancoDTO

	ruido, encontrado := this.repository.BuscarRuidoBlanco(id)

	if encontrado {
		var ruidoDTO dto.MetadataRuidoBlancoDTO
		ruidoDTO.ID = ruido.GetId()
		ruidoDTO.Tipo = ruido.GetTipo()
		ruidoDTO.TipoSonido = ruido.GetTipoSonido()
		ruidoDTO.FuenteAudio = ruido.GetFuenteAudio()
		ruidoDTO.UsoSugerido = ruido.GetUsoSugerido()
		ruidoDTO.Proveedor = ruido.GetProveedor()
		ruidoDTO.Duracion = ruido.GetDuracion()
		ruidoDTO.FrecuenciaDominante = ruido.GetFrecuenciaDominante()

		respuesta.ObjRuidoBlanco = ruidoDTO
		respuesta.Codigo = 200
		respuesta.Mensaje = "Métadata del ruido blanco encontrada"
	} else {
		respuesta.Codigo = 404
		respuesta.Mensaje = "La métadata del ruido blanco no se encontró"
	}
	return respuesta
}

/**
 * @brief Obtiene la lista de todos los registros de ruido blanco.
 * 
 * @return dto.RespuestaListaRuidoBlancoDTO DTO con la lista y código de respuesta.
 */
func (this *MetadataRuidoBlancoService) ListarRuidoBlanco() dto.RespuestaListaRuidoBlancoDTO {
	lista := this.repository.ListarRuidoBlanco()
	var dtos []dto.MetadataRuidoBlancoDTO

	for _, ruido := range lista {
		var d dto.MetadataRuidoBlancoDTO
		d.ID = ruido.GetId()
		d.Tipo = ruido.GetTipo()
		d.TipoSonido = ruido.GetTipoSonido()
		d.FuenteAudio = ruido.GetFuenteAudio()
		d.UsoSugerido = ruido.GetUsoSugerido()
		d.Proveedor = ruido.GetProveedor()
		d.Duracion = ruido.GetDuracion()
		d.FrecuenciaDominante = ruido.GetFrecuenciaDominante()

		dtos = append(dtos, d)
	}

	return dto.RespuestaListaRuidoBlancoDTO{
		RuidosBlancos: dtos,
		Codigo:        200,
		Mensaje:       "Lista de ruido blanco obtenida correctamente",
	}
}

/**
 * @brief Registra un nuevo ruido blanco y coordina internamente el guardado del archivo físico.
 * 
 * @param ruidoDTO DTO con los datos del ruido blanco.
 * @param archivoBytes Contenido binario del archivo de audio (opcional).
 * @param nombreArchivo Nombre del archivo en disco.
 * @return dto.RespuestaMetadataRuidoBlancoDTO DTO con la respuesta y entidad creada.
 */
func (this *MetadataRuidoBlancoService) RegistrarRuidoBlanco(ruidoDTO dto.MetadataRuidoBlancoDTO, archivoBytes []byte, nombreArchivo string) dto.RespuestaMetadataRuidoBlancoDTO {
	if len(archivoBytes) > 0 && this.audioService != nil {
		if nombreArchivo == "" {
			nombreArchivo = fmt.Sprintf("%s.mp3", ruidoDTO.FuenteAudio)
		}
		err := this.audioService.GuardarArchivo(nombreArchivo, archivoBytes)
		if err != nil {
			return dto.RespuestaMetadataRuidoBlancoDTO{
				Codigo:  500,
				Mensaje: "Error al guardar el archivo físico de audio: " + err.Error(),
			}
		}
	}

	var entidad entity.MetaDataRuidoBlanco
	entidad.SetId(ruidoDTO.ID)
	if ruidoDTO.Tipo.GetId() > 0 {
		entidad.SetTipo(ruidoDTO.Tipo)
	} else {
		entidad.SetTipo(entity.NewTipoAudio(2, "Ruido Blanco"))
	}
	entidad.SetTipoSonido(ruidoDTO.TipoSonido)
	entidad.SetFuenteAudio(ruidoDTO.FuenteAudio)
	entidad.SetUsoSugerido(ruidoDTO.UsoSugerido)
	entidad.SetProveedor(ruidoDTO.Proveedor)
	entidad.SetDuracion(ruidoDTO.Duracion)
	entidad.SetFrecuenciaDominante(ruidoDTO.FrecuenciaDominante)

	guardado := this.repository.GuardarRuidoBlanco(entidad)

	ruidoDTO.ID = guardado.GetId()
	ruidoDTO.Tipo = guardado.GetTipo()

	return dto.RespuestaMetadataRuidoBlancoDTO{
		ObjRuidoBlanco: ruidoDTO,
		Codigo:         201,
		Mensaje:        "Ruido blanco y archivo registrados exitosamente",
	}
}
