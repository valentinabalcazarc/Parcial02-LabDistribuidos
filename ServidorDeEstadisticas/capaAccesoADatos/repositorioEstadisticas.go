/**
 * @file repositorioEstadisticas.go
 * @brief Repositorio en memoria thread-safe para almacenar las estadísticas de reproducción.
 */
package capaaccesoadatos

import (
	"sort"
	"sync"

	dtos "servidorEstadisticas/DTOs"
)

/**
 * @struct RepositorioEstadisticas
 * @brief Gestiona el almacenamiento en memoria de las reproducciones por cliente y globales.
 */
type RepositorioEstadisticas struct {
	mu sync.RWMutex
	// Mapa: clienteId -> (tituloAudio -> *EstadisticaAudioDTO)
	reproduccionesPorCliente map[string]map[string]*dtos.EstadisticaAudioDTO
	// Mapa: tituloAudio -> *RankingAudioDTO
	reproduccionesGlobales map[string]*dtos.RankingAudioDTO
	// Historial cronológico de notificaciones
	historialEventos []dtos.NotificacionReproduccion
}

var (
	instancia *RepositorioEstadisticas
	once      sync.Once
)

/**
 * @brief Obtiene la instancia única (Singleton) del repositorio de estadísticas.
 */
func GetRepositorioEstadisticas() *RepositorioEstadisticas {
	once.Do(func() {
		instancia = &RepositorioEstadisticas{
			reproduccionesPorCliente: make(map[string]map[string]*dtos.EstadisticaAudioDTO),
			reproduccionesGlobales:   make(map[string]*dtos.RankingAudioDTO),
			historialEventos:         make([]dtos.NotificacionReproduccion, 0),
		}
	})
	return instancia
}

/**
 * @brief Registra un evento de reproducción recibido desde la cola.
 *
 * @param notif Información del evento de reproducción.
 */
func (r *RepositorioEstadisticas) RegistrarReproduccion(notif dtos.NotificacionReproduccion) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cliente := notif.Cliente
	if cliente == "" {
		cliente = "anonimo"
	}

	titulo := notif.Titulo
	if titulo == "" {
		titulo = "audio_sin_titulo"
	}

	// 1. Registrar para el cliente específico
	if _, ok := r.reproduccionesPorCliente[cliente]; !ok {
		r.reproduccionesPorCliente[cliente] = make(map[string]*dtos.EstadisticaAudioDTO)
	}

	if est, ok := r.reproduccionesPorCliente[cliente][titulo]; ok {
		est.CantidadVeces++
		est.UltimaReproduccion = notif.FechaHora
		if notif.TipoAudio != "" {
			est.TipoAudio = notif.TipoAudio
		}
	} else {
		r.reproduccionesPorCliente[cliente][titulo] = &dtos.EstadisticaAudioDTO{
			Titulo:             titulo,
			TipoAudio:          notif.TipoAudio,
			CantidadVeces:      1,
			UltimaReproduccion: notif.FechaHora,
		}
	}

	// 2. Registrar en acumulador global
	if global, ok := r.reproduccionesGlobales[titulo]; ok {
		global.TotalReproducciones++
		global.UltimaReproduccion = notif.FechaHora
	} else {
		r.reproduccionesGlobales[titulo] = &dtos.RankingAudioDTO{
			Titulo:              titulo,
			TotalReproducciones: 1,
			UltimaReproduccion:  notif.FechaHora,
		}
	}

	// 3. Añadir al historial (máximo los últimos 200 eventos)
	r.historialEventos = append(r.historialEventos, notif)
	if len(r.historialEventos) > 200 {
		r.historialEventos = r.historialEventos[len(r.historialEventos)-200:]
	}
}

/**
 * @brief Obtiene las reproducciones de un cliente.
 *
 * @param cliente Identificador del cliente a consultar.
 * @return dtos.ResumenClienteDTO Resumen con los audios y cuántas veces se reprodujeron.
 */
func (r *RepositorioEstadisticas) ObtenerReproduccionesCliente(cliente string) dtos.ResumenClienteDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resumen := dtos.ResumenClienteDTO{
		Cliente:             cliente,
		TotalReproducciones: 0,
		Audios:              make([]dtos.EstadisticaAudioDTO, 0),
	}

	audiosMap, ok := r.reproduccionesPorCliente[cliente]
	if !ok {
		return resumen
	}

	for _, audio := range audiosMap {
		resumen.TotalReproducciones += audio.CantidadVeces
		resumen.Audios = append(resumen.Audios, *audio)
	}

	// Ordenar los audios del cliente por mayor número de reproducciones
	sort.Slice(resumen.Audios, func(i, j int) bool {
		return resumen.Audios[i].CantidadVeces > resumen.Audios[j].CantidadVeces
	})

	return resumen
}

/**
 * @brief Obtiene el listado global de audios con sus reproducciones ordenados por popularidad.
 */
func (r *RepositorioEstadisticas) ObtenerRankingGlobal() []dtos.RankingAudioDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ranking := make([]dtos.RankingAudioDTO, 0, len(r.reproduccionesGlobales))
	for _, stat := range r.reproduccionesGlobales {
		ranking = append(ranking, *stat)
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].TotalReproducciones > ranking[j].TotalReproducciones
	})

	return ranking
}

/**
 * @brief Obtiene los eventos de reproducción más recientes.
 */
func (r *RepositorioEstadisticas) ObtenerHistorial(limite int) []dtos.NotificacionReproduccion {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.historialEventos)
	if limite <= 0 || limite > total {
		limite = total
	}

	resultado := make([]dtos.NotificacionReproduccion, limite)
	for i := 0; i < limite; i++ {
		// En orden cronológico inverso (el más reciente primero)
		resultado[i] = r.historialEventos[total-1-i]
	}

	return resultado
}
