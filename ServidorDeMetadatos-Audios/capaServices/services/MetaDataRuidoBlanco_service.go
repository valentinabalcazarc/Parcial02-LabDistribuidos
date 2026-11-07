package services

import (
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de Ruido Blanco.
 * 
 * Es la fachada (facade) que expone al controlador las operaciones de negocio,
 * ocultando el acceso al repositorio y la conversión entre Entity y DTO.
 */
type MetadataRuidoBlancoService struct {
	repository *repository.MetadataRuidoBlancoRepository
}

/**
 * @brief Crea una nueva instancia del servicio de MetadataRuidoBlanco.
 * 
 * @param repository Repositorio de MetadataRuidoBlanco inyectado.
 * @return *MetadataRuidoBlancoService Instancia del servicio.
 */
func NewMetadataRuidoBlancoService(repository *repository.MetadataRuidoBlancoRepository) *MetadataRuidoBlancoService {
	return &MetadataRuidoBlancoService{repository: repository}
}

/**
 * @brief Consulta un ruido blanco por su ID.
 * 
 * Recibe un ID, busca el Entity en el repositorio y lo convierte a
 * RespuestaMetadataRuidoBlancoDTO con el código y mensaje según el resultado.
 * 
 * @param id Identificador del ruido blanco a buscar.
 * @return dto.RespuestaMetadataRuidoBlancoDTO DTO con la respuesta de la consulta.
 */
func (this *MetadataRuidoBlancoService) ConsultarRuidoBlanco(id int) dto.RespuestaMetadataRuidoBlancoDTO {
	var respuesta dto.RespuestaMetadataRuidoBlancoDTO

	ruidoBlanco, encontrado := this.repository.BuscarRuidoBlanco(id)

	if encontrado {
		var ruidoBlancoDTO dto.MetadataRuidoBlancoDTO
		ruidoBlancoDTO.ID = ruidoBlanco.GetId()
		ruidoBlancoDTO.Tipo = ruidoBlanco.GetTipo()
		ruidoBlancoDTO.TipoSonido = ruidoBlanco.GetTipoSonido()
		ruidoBlancoDTO.FuenteAudio = ruidoBlanco.GetFuenteAudio()
		ruidoBlancoDTO.UsoSugerido = ruidoBlanco.GetUsoSugerido()
		ruidoBlancoDTO.Proveedor = ruidoBlanco.GetProveedor()
		ruidoBlancoDTO.Duracion = ruidoBlanco.GetDuracion()
		ruidoBlancoDTO.FrecuenciaDominante = ruidoBlanco.GetFrecuenciaDominante()

		respuesta.ObjRuidoBlanco = ruidoBlancoDTO
		respuesta.Codigo = 200
		respuesta.Mensaje = "Métadata del ruidoBlanco encontrada"
	} else {
		respuesta.Codigo = 400
		respuesta.Mensaje = "La métadata del ruidoBlanco no se encontró"
	}
	return respuesta
}

/**
 * @brief Obtiene la lista de todos los ruidos blancos.
 * 
 * @return dto.RespuestaListaRuidoBlancoDTO DTO con la lista y código de respuesta.
 */
func (this *MetadataRuidoBlancoService) ListarRuidoBlanco() dto.RespuestaListaRuidoBlancoDTO {
	lista := this.repository.ListarRuidoBlanco()
	var dtos []dto.MetadataRuidoBlancoDTO

	for _, ruidoBlanco := range lista {
		var d dto.MetadataRuidoBlancoDTO
		d.ID = ruidoBlanco.GetId()
		d.Tipo = ruidoBlanco.GetTipo()
		d.TipoSonido = ruidoBlanco.GetTipoSonido()
		d.FuenteAudio = ruidoBlanco.GetFuenteAudio()
		d.UsoSugerido = ruidoBlanco.GetUsoSugerido()
		d.Proveedor = ruidoBlanco.GetProveedor()
		d.Duracion = ruidoBlanco.GetDuracion()
		d.FrecuenciaDominante = ruidoBlanco.GetFrecuenciaDominante()

		dtos = append(dtos, d)
	}

	return dto.RespuestaListaRuidoBlancoDTO{
		RuidosBlancos: dtos,
		Codigo:        200,
		Mensaje:       "Lista de ruido blanco obtenida correctamente",
	}
}
