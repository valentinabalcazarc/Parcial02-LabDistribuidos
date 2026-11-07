package repository

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Repositorio que mantiene en memoria el slice de música.
 * 
 * Expone las operaciones de búsqueda y carga para música.
 */
type MetadataMusicaRepository struct {
	vectorMetadataMusica []entity.MetaDataMusica
}

/**
 * @brief Crea el repositorio y lo precarga con metadata de canciones.
 * 
 * @return *MetadataMusicaRepository Instancia del repositorio.
 */
func NewMetadataMusicaRepository() *MetadataMusicaRepository {
	this := &MetadataMusicaRepository{}
	this.CargarMetadataMusica()
	return this
}

/**
 * @brief Inicializa el vector con 2 canciones de ejemplo.
 */
func (this *MetadataMusicaRepository) CargarMetadataMusica() {
	var objMusica1, objMusica2 entity.MetaDataMusica

	tipoMusica := entity.NewTipoAudio(1, "Música")

	// Canción 1
	objMusica1.SetId(3)
	objMusica1.SetTipo(tipoMusica)
	objMusica1.SetArtistaPrincipal("Jarabe de Palo")
	objMusica1.SetAlbum("La Flaca")
	objMusica1.SetGenero("Rock")
	objMusica1.SetTituloCancion("La Flaca")
	objMusica1.SetSelloDiscografico("Virgin Records")
	objMusica1.SetAnioLanzamiento("1996")

	// Canción 2
	objMusica2.SetId(4)
	objMusica2.SetTipo(tipoMusica)
	objMusica2.SetArtistaPrincipal("Michael Jackson")
	objMusica2.SetAlbum("Thriller")
	objMusica2.SetGenero("Pop")
	objMusica2.SetTituloCancion("Billie Jean")
	objMusica2.SetSelloDiscografico("Epic Records")
	objMusica2.SetAnioLanzamiento("1982")

	// Asignación al slice del repositorio
	this.vectorMetadataMusica = []entity.MetaDataMusica{
		objMusica1,
		objMusica2,
	}
}

/**
 * @brief Busca una canción por su ID.
 * 
 * @param id Identificador de la música a buscar.
 * @return entity.MetaDataMusica Objeto encontrado.
 * @return bool Indica si la búsqueda tuvo éxito.
 */
func (this *MetadataMusicaRepository) BuscarMusica(id int) (entity.MetaDataMusica, bool) {
	for _, musica := range this.vectorMetadataMusica {
		if musica.GetId() == id {
			return musica, true
		}
	}

	return entity.MetaDataMusica{}, false
}

/**
 * @brief Retorna toda la música almacenada.
 * 
 * @return []entity.MetaDataMusica Arreglo con la música.
 */
func (this *MetadataMusicaRepository) ListarMusica() []entity.MetaDataMusica {
	return this.vectorMetadataMusica
}
