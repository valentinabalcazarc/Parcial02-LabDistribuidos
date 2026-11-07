package repository

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Repositorio para la gestión de TipoAudio.
 * 
 * Mantiene en memoria un vector con los tipos de audio predefinidos.
 */
type TipoAudioRepository struct {
	vectorTipos []entity.TipoAudio
}

/**
 * @brief Crea una nueva instancia de TipoAudioRepository.
 * 
 * @return *TipoAudioRepository Instancia del repositorio con los tipos cargados.
 */
func NewTipoAudioRepository() *TipoAudioRepository {
	repo := &TipoAudioRepository{}
	repo.CargarTipos()
	return repo
}

/**
 * @brief Carga en memoria los tipos de audio iniciales.
 */
func (r *TipoAudioRepository) CargarTipos() {
	r.vectorTipos = []entity.TipoAudio{
		entity.NewTipoAudio(1, "Música"),
		entity.NewTipoAudio(2, "Ruido Blanco"),
		entity.NewTipoAudio(3, "Podcast"),
		entity.NewTipoAudio(4, "Audiolibro"),
	}
}

/**
 * @brief Retorna todos los tipos de audio almacenados.
 * 
 * @return []entity.TipoAudio Arreglo con los tipos de audio.
 */
func (r *TipoAudioRepository) ListarTipos() []entity.TipoAudio {
	return r.vectorTipos
}
