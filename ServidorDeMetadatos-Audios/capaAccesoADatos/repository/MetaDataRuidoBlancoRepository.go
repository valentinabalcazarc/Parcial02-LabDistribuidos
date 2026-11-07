package repository

import "servidorMetadatos/capaAccesoADatos/entity"

/**
 * @brief Repositorio que mantiene en memoria el slice de ruido blanco.
 * 
 * Expone las operaciones de búsqueda y carga de ruido blanco.
 */
type MetadataRuidoBlancoRepository struct {
	vectorMetadataRuidoBlanco []entity.MetaDataRuidoBlanco
}

/**
 * @brief Crea el repositorio y lo precarga con audios de ruido blanco.
 * 
 * @return *MetadataRuidoBlancoRepository Instancia del repositorio.
 */
func NewMetadataRuidoBlancoRepository() *MetadataRuidoBlancoRepository {
	this := &MetadataRuidoBlancoRepository{}
	this.CargarMetadataRuidoBlanco()
	return this
}

/**
 * @brief Inicializa el vector con 2 ruidos blancos de ejemplo.
 */
func (this *MetadataRuidoBlancoRepository) CargarMetadataRuidoBlanco() {
	var objRuido1, objRuido2 entity.MetaDataRuidoBlanco

	tipoRuidoBlanco := entity.NewTipoAudio(2, "Ruido Blanco")

	// Ruido 1: Ruido Blanco / Lluvia
	objRuido1.SetId(7)
	objRuido1.SetTipo(tipoRuidoBlanco)
	objRuido1.SetTipoSonido("Ruido Blanco")
	objRuido1.SetFuenteAudio("Lluvia")
	objRuido1.SetUsoSugerido("Dormir")
	objRuido1.SetProveedor("SleepSounds Studio")
	objRuido1.SetDuracion(60) // Duración del bucle en segundos (1 min)
	objRuido1.SetFrecuenciaDominante("Agudos")

	// Ruido 2: Ruido Marrón / Ventilador
	objRuido2.SetId(8)
	objRuido2.SetTipo(tipoRuidoBlanco)
	objRuido2.SetTipoSonido("Ruido Marrón")
	objRuido2.SetFuenteAudio("Ventilador")
	objRuido2.SetUsoSugerido("Concentración")
	objRuido2.SetProveedor("FocusLoop Channel")
	objRuido2.SetDuracion(120) // Duración del bucle en segundos (2 min)
	objRuido2.SetFrecuenciaDominante("Graves")

	// Asignación al slice del repositorio
	this.vectorMetadataRuidoBlanco = []entity.MetaDataRuidoBlanco{
		objRuido1,
		objRuido2,
	}
}

/**
 * @brief Busca un registro de ruido blanco por su ID.
 * 
 * @param id Identificador del ruido blanco a buscar.
 * @return entity.MetaDataRuidoBlanco Objeto encontrado.
 * @return bool Indica si la búsqueda tuvo éxito.
 */
func (this *MetadataRuidoBlancoRepository) BuscarRuidoBlanco(id int) (entity.MetaDataRuidoBlanco, bool) {
	for _, ruido := range this.vectorMetadataRuidoBlanco {
		if ruido.GetId() == id {
			return ruido, true
		}
	}

	return entity.MetaDataRuidoBlanco{}, false
}

/**
 * @brief Retorna todos los registros de ruido blanco almacenados.
 * 
 * @return []entity.MetaDataRuidoBlanco Arreglo con los ruidos blancos.
 */
func (this *MetadataRuidoBlancoRepository) ListarRuidoBlanco() []entity.MetaDataRuidoBlanco {
	return this.vectorMetadataRuidoBlanco
}
