/**
 * @file estadisticas_service.go
 * @brief Capa de servicio para la gestión y procesamiento de estadísticas de reproducción.
 */
package capaservices

import (
	"fmt"
	dtos "servidorEstadisticas/DTOs"
	capaaccesoadatos "servidorEstadisticas/capaAccesoADatos"
)

/**
 * @struct EstadisticasService
 * @brief Servicio que procesa eventos de la cola y provee consultas estadísticas.
 */
type EstadisticasService struct {
	repo *capaaccesoadatos.RepositorioEstadisticas
}

/**
 * @brief Crea una nueva instancia de EstadisticasService.
 */
func NewEstadisticasService(repo *capaaccesoadatos.RepositorioEstadisticas) *EstadisticasService {
	return &EstadisticasService{
		repo: repo,
	}
}

/**
 * @brief Procesa una notificación recibida de RabbitMQ y la persiste en el repositorio.
 */
func (s *EstadisticasService) ProcesarNotificacion(notif dtos.NotificacionReproduccion) {
	s.repo.RegistrarReproduccion(notif)
	fmt.Printf("[EstadisticasService] Reproducción registrada: Cliente='%s', Audio='%s'\n", notif.Cliente, notif.Titulo)
}

/**
 * @brief Responde a la consulta de audios reproducidos por un cliente (Requerimiento d).
 */
func (s *EstadisticasService) ConsultarReproduccionesCliente(cliente string) dtos.RespuestaAPI {
	if cliente == "" {
		return dtos.RespuestaAPI{
			Codigo:  400,
			Mensaje: "El parámetro de cliente es requerido",
		}
	}

	resumen := s.repo.ObtenerReproduccionesCliente(cliente)
	return dtos.RespuestaAPI{
		Codigo:  200,
		Mensaje: fmt.Sprintf("Estadísticas del cliente '%s' obtenidas exitosamente", cliente),
		Datos:   resumen,
	}
}

/**
 * @brief Obtiene el ranking general de audios más reproducidos.
 */
func (s *EstadisticasService) ConsultarRankingGlobal() dtos.RespuestaAPI {
	ranking := s.repo.ObtenerRankingGlobal()
	return dtos.RespuestaAPI{
		Codigo:  200,
		Mensaje: "Ranking global de reproducciones obtenido exitosamente",
		Datos:   ranking,
	}
}

/**
 * @brief Obtiene los eventos recientes de reproducción.
 */
func (s *EstadisticasService) ConsultarHistorial(limite int) dtos.RespuestaAPI {
	historial := s.repo.ObtenerHistorial(limite)
	return dtos.RespuestaAPI{
		Codigo:  200,
		Mensaje: "Historial de reproducciones obtenido exitosamente",
		Datos:   historial,
	}
}

