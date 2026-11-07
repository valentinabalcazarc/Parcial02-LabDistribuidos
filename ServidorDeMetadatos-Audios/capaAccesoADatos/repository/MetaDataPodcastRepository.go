package repository

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Repositorio que mantiene en memoria el slice de podcasts.
 * 
 * Expone las operaciones de búsqueda y carga de podcasts.
 */
type MetadataPodcastRepository struct {
	vectorMetadataPodcasts []entity.MetaDataPodcast
}

/**
 * @brief Crea el repositorio y lo precarga con episodios de podcast.
 * 
 * @return *MetadataPodcastRepository Instancia del repositorio.
 */
func NewMetadataPodcastRepository() *MetadataPodcastRepository {
	this := &MetadataPodcastRepository{}
	this.CargarMetadataPodcasts()
	return this
}

/**
 * @brief Inicializa el vector con 2 episodios de ejemplo.
 */
func (this *MetadataPodcastRepository) CargarMetadataPodcasts() {
	var objPod1, objPod2 entity.MetaDataPodcast

	tipoPodcast := entity.NewTipoAudio(3, "Podcast")

	// Podcast 1
	objPod1.SetId(5)
	objPod1.SetTipo(tipoPodcast)
	objPod1.SetNombre("Radio Ambulante")
	objPod1.SetTituloEpisodio("El Polizón")
	objPod1.SetNumeroTemporada(12)
	objPod1.SetNotasShow("Historia de un viaje inesperado hacia el norte.")
	objPod1.SetClasificacionContenido("Explicito")

	// Podcast 2
	objPod2.SetId(6)
	objPod2.SetTipo(tipoPodcast)
	objPod2.SetNombre("The Daily")
	objPod2.SetTituloEpisodio("Global Economy Trends")
	objPod2.SetNumeroTemporada(8)
	objPod2.SetNotasShow("Análisis semanal de los mercados internacionales.")
	objPod2.SetClasificacionContenido("Para toda la familia")

	// Asignación al slice del repositorio
	this.vectorMetadataPodcasts = []entity.MetaDataPodcast{
		objPod1,
		objPod2,
	}
}

/**
 * @brief Busca un podcast por su ID.
 * 
 * @param id Identificador del podcast a buscar.
 * @return entity.MetaDataPodcast Objeto encontrado.
 * @return bool Indica si la búsqueda tuvo éxito.
 */
func (this *MetadataPodcastRepository) BuscarPodcast(id int) (entity.MetaDataPodcast, bool) {
	for _, podcast := range this.vectorMetadataPodcasts {
		if podcast.GetId() == id {
			return podcast, true
		}
	}

	return entity.MetaDataPodcast{}, false
}

/**
 * @brief Retorna todos los podcasts almacenados.
 * 
 * @return []entity.MetaDataPodcast Arreglo con los podcasts.
 */
func (this *MetadataPodcastRepository) ListarPodcasts() []entity.MetaDataPodcast {
	return this.vectorMetadataPodcasts
}
