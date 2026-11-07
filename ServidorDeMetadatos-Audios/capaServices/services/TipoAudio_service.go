package services

import (
	"servidorMetadatos/capaAccesoADatos/repository"
	"servidorMetadatos/capaServices/dto"
)

/**
 * @brief Servicio para gestionar las operaciones de TipoAudio.
 * 
 * Actúa como intermediario entre el controlador y el repositorio.
 */
type TipoAudioService struct {
	repository *repository.TipoAudioRepository
}

/**
 * @brief Crea una nueva instancia del servicio de TipoAudio.
 * 
 * @param repository Repositorio de TipoAudio inyectado.
 * @return *TipoAudioService Instancia del servicio.
 */
func NewTipoAudioService(repository *repository.TipoAudioRepository) *TipoAudioService {
	return &TipoAudioService{repository: repository}
}

/**
 * @brief Obtiene la lista de todos los tipos de audio.
 * 
 * @return dto.RespuestaTiposAudioDTO DTO con la lista y código de respuesta.
 */
func (s *TipoAudioService) ListarTipos() dto.RespuestaTiposAudioDTO {
	tipos := s.repository.ListarTipos()
	return dto.RespuestaTiposAudioDTO{
		Tipos:   tipos,
		Codigo:  200,
		Mensaje: "Tipos de audio obtenidos correctamente",
	}
}
