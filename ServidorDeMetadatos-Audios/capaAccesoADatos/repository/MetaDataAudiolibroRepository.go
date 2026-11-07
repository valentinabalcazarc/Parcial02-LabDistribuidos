package repository

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Repositorio que mantiene en memoria el slice de audiolibros.
 * 
 * Expone las operaciones de búsqueda y registro para los audiolibros.
 */
type MetadataAudiolibroRepository struct {
	vectorMetadataAudiolibros []entity.MetaDataAudiolibro
}

/**
 * @brief Crea el repositorio y lo precarga con audiolibros de ejemplo.
 * 
 * @return *MetadataAudiolibroRepository Instancia del repositorio.
 */
func NewMetadataAudiolibroRepository() *MetadataAudiolibroRepository {
	this := &MetadataAudiolibroRepository{}
	this.CargarMetadataAudiolibros()
	return this
}

/**
 * @brief Inicializa el vector con 2 audiolibros de ejemplo.
 */
func (this *MetadataAudiolibroRepository) CargarMetadataAudiolibros() {
	var objAudio1, objAudio2 entity.MetaDataAudiolibro

	tipoAudiolibro := entity.NewTipoAudio(4, "Audiolibro")

	// Audiolibro 1
	objAudio1.SetId(1)
	objAudio1.SetTipo(tipoAudiolibro)
	objAudio1.SetTituloLibro("Cien Años de Soledad")
	objAudio1.SetAutor("Gabriel García Márquez")
	objAudio1.SetNarrador("Gustavo Bonfigli")
	objAudio1.SetEditorial("Penguin Random House")
	objAudio1.SetISBN(978030747)
	objAudio1.SetCapitulo(1)

	// Audiolibro 2
	objAudio2.SetId(2)
	objAudio2.SetTipo(tipoAudiolibro)
	objAudio2.SetTituloLibro("El Principito")
	objAudio2.SetAutor("Antoine de Saint-Exupéry")
	objAudio2.SetNarrador("José María Carrascal")
	objAudio2.SetEditorial("Salamandra")
	objAudio2.SetISBN(978847888)
	objAudio2.SetCapitulo(3)

	// Asignación al slice del repositorio
	this.vectorMetadataAudiolibros = []entity.MetaDataAudiolibro{
		objAudio1,
		objAudio2,
	}
}

/**
 * @brief Busca un audiolibro por su ID.
 * 
 * @param id Identificador del audiolibro a buscar.
 * @return entity.MetaDataAudiolibro Objeto encontrado.
 * @return bool Indica si la búsqueda tuvo éxito.
 */
func (this *MetadataAudiolibroRepository) BuscarAudiolibro(id int) (entity.MetaDataAudiolibro, bool) {
	for _, audio := range this.vectorMetadataAudiolibros {
		if audio.GetId() == id {
			return audio, true
		}
	}

	return entity.MetaDataAudiolibro{}, false
}

/**
 * @brief Retorna todos los audiolibros almacenados.
 * 
 * @return []entity.MetaDataAudiolibro Arreglo con los audiolibros.
 */
func (this *MetadataAudiolibroRepository) ListarAudiolibros() []entity.MetaDataAudiolibro {
	return this.vectorMetadataAudiolibros
}
