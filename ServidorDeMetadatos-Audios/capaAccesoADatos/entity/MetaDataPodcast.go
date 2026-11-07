package entity

import "encoding/json"

/**
 * @brief Estructura que representa la metadata de un podcast.
 * 
 * Contiene datos referentes a un episodio de podcast, incluyendo
 * su nombre, título, temporada y notas del show.
 */
type MetaDataPodcast struct {
	id                     int
	tipo                   TipoAudio
	nombre                 string
	tituloEpisodio         string
	numeroTemporada        int
	notasShow              string
	clasificacionContenido string
}

// --- Getters ---

func (p *MetaDataPodcast) GetId() int {
	return p.id
}

func (p *MetaDataPodcast) GetTipo() TipoAudio {
	return p.tipo
}

func (p *MetaDataPodcast) GetNombre() string {
	return p.nombre
}

func (p *MetaDataPodcast) GetTituloEpisodio() string {
	return p.tituloEpisodio
}

func (p *MetaDataPodcast) GetNumeroTemporada() int {
	return p.numeroTemporada
}

func (p *MetaDataPodcast) GetNotasShow() string {
	return p.notasShow
}

func (p *MetaDataPodcast) GetClasificacionContenido() string {
	return p.clasificacionContenido
}

// --- Setters ---

func (p *MetaDataPodcast) SetId(id int) {
	p.id = id
}

func (p *MetaDataPodcast) SetTipo(tipo TipoAudio) {
	p.tipo = tipo
}

func (p *MetaDataPodcast) SetNombre(nombre string) {
	p.nombre = nombre
}

func (p *MetaDataPodcast) SetTituloEpisodio(tituloEpisodio string) {
	p.tituloEpisodio = tituloEpisodio
}

func (p *MetaDataPodcast) SetNumeroTemporada(numeroTemporada int) {
	p.numeroTemporada = numeroTemporada
}

func (p *MetaDataPodcast) SetNotasShow(notasShow string) {
	p.notasShow = notasShow
}

func (p *MetaDataPodcast) SetClasificacionContenido(clasificacionContenido string) {
	p.clasificacionContenido = clasificacionContenido
}

// --- JSON Marshaling ---

func (p MetaDataPodcast) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                     int       `json:"id"`
		Tipo                   TipoAudio `json:"tipo"`
		Nombre                 string    `json:"nombre"`
		TituloEpisodio         string    `json:"titulo_episodio"`
		NumeroTemporada        int       `json:"numero_temporada"`
		NotasShow              string    `json:"notas_show"`
		ClasificacionContenido string    `json:"clasificacion_contenido"`
	}{
		ID:                     p.id,
		Tipo:                   p.tipo,
		Nombre:                 p.nombre,
		TituloEpisodio:         p.tituloEpisodio,
		NumeroTemporada:        p.numeroTemporada,
		NotasShow:              p.notasShow,
		ClasificacionContenido: p.clasificacionContenido,
	})
}

func (p *MetaDataPodcast) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID                     int       `json:"id"`
		Tipo                   TipoAudio `json:"tipo"`
		Nombre                 string    `json:"nombre"`
		TituloEpisodio         string    `json:"titulo_episodio"`
		NumeroTemporada        int       `json:"numero_temporada"`
		NotasShow              string    `json:"notas_show"`
		ClasificacionContenido string    `json:"clasificacion_contenido"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.id = aux.ID
	p.tipo = aux.Tipo
	p.nombre = aux.Nombre
	p.tituloEpisodio = aux.TituloEpisodio
	p.numeroTemporada = aux.NumeroTemporada
	p.notasShow = aux.NotasShow
	p.clasificacionContenido = aux.ClasificacionContenido

	return nil
}
