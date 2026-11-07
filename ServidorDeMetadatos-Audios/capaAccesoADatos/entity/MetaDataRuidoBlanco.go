package entity

import "encoding/json"

/**
 * @brief Estructura que representa la metadata de ruido blanco.
 * 
 * Almacena atributos específicos de un audio de ruido blanco,
 * como el tipo de sonido, fuente, uso sugerido y duración.
 */
type MetaDataRuidoBlanco struct {
	id                  int
	tipo                TipoAudio
	tipoSonido          string
	fuenteAudio         string
	usoSugerido         string
	proveedor           string
	duracion            int
	frecuenciaDominante string
}

// --- Getters ---

func (r *MetaDataRuidoBlanco) GetId() int {
	return r.id
}

func (r *MetaDataRuidoBlanco) GetTipo() TipoAudio {
	return r.tipo
}

func (r *MetaDataRuidoBlanco) GetTipoSonido() string {
	return r.tipoSonido
}

func (r *MetaDataRuidoBlanco) GetFuenteAudio() string {
	return r.fuenteAudio
}

func (r *MetaDataRuidoBlanco) GetUsoSugerido() string {
	return r.usoSugerido
}

func (r *MetaDataRuidoBlanco) GetProveedor() string {
	return r.proveedor
}

func (r *MetaDataRuidoBlanco) GetDuracion() int {
	return r.duracion
}

func (r *MetaDataRuidoBlanco) GetFrecuenciaDominante() string {
	return r.frecuenciaDominante
}

// --- Setters ---

func (r *MetaDataRuidoBlanco) SetId(id int) {
	r.id = id
}

func (r *MetaDataRuidoBlanco) SetTipo(tipo TipoAudio) {
	r.tipo = tipo
}

func (r *MetaDataRuidoBlanco) SetTipoSonido(tipoSonido string) {
	r.tipoSonido = tipoSonido
}

func (r *MetaDataRuidoBlanco) SetFuenteAudio(fuenteAudio string) {
	r.fuenteAudio = fuenteAudio
}

func (r *MetaDataRuidoBlanco) SetUsoSugerido(usoSugerido string) {
	r.usoSugerido = usoSugerido
}

func (r *MetaDataRuidoBlanco) SetProveedor(proveedor string) {
	r.proveedor = proveedor
}

func (r *MetaDataRuidoBlanco) SetDuracion(duracion int) {
	r.duracion = duracion
}

func (r *MetaDataRuidoBlanco) SetFrecuenciaDominante(frecuenciaDominante string) {
	r.frecuenciaDominante = frecuenciaDominante
}

// --- JSON Marshaling ---

func (r MetaDataRuidoBlanco) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                  int       `json:"id"`
		Tipo                TipoAudio `json:"tipo"`
		TipoSonido          string    `json:"tipo_sonido"`
		FuenteAudio         string    `json:"fuente_audio"`
		UsoSugerido         string    `json:"uso_sugerido"`
		Proveedor           string    `json:"proveedor"`
		Duracion            int       `json:"duracion"`
		FrecuenciaDominante string    `json:"frecuencia_dominante"`
	}{
		ID:                  r.id,
		Tipo:                r.tipo,
		TipoSonido:          r.tipoSonido,
		FuenteAudio:         r.fuenteAudio,
		UsoSugerido:         r.usoSugerido,
		Proveedor:           r.proveedor,
		Duracion:            r.duracion,
		FrecuenciaDominante: r.frecuenciaDominante,
	})
}

func (r *MetaDataRuidoBlanco) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID                  int       `json:"id"`
		Tipo                TipoAudio `json:"tipo"`
		TipoSonido          string    `json:"tipo_sonido"`
		FuenteAudio         string    `json:"fuente_audio"`
		UsoSugerido         string    `json:"uso_sugerido"`
		Proveedor           string    `json:"proveedor"`
		Duracion            int       `json:"duracion"`
		FrecuenciaDominante string    `json:"frecuencia_dominante"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.id = aux.ID
	r.tipo = aux.Tipo
	r.tipoSonido = aux.TipoSonido
	r.fuenteAudio = aux.FuenteAudio
	r.usoSugerido = aux.UsoSugerido
	r.proveedor = aux.Proveedor
	r.duracion = aux.Duracion
	r.frecuenciaDominante = aux.FrecuenciaDominante

	return nil
}
