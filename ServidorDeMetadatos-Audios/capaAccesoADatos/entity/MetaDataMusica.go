package entity

import "encoding/json"

/**
 * @brief Estructura que representa la metadata de música.
 * 
 * Contiene información sobre una pista musical, como el artista,
 * álbum, género, título de la canción, etc.
 */
type MetaDataMusica struct {
	id                int
	tipo              TipoAudio
	artistaPrincipal  string
	album             string
	genero            string
	tituloCancion     string
	selloDiscografico string
	anioLanzamiento   string
}

// --- Getters ---

func (m *MetaDataMusica) GetId() int {
	return m.id
}

func (m *MetaDataMusica) GetTipo() TipoAudio {
	return m.tipo
}

func (m *MetaDataMusica) GetArtistaPrincipal() string {
	return m.artistaPrincipal
}

func (m *MetaDataMusica) GetAlbum() string {
	return m.album
}

func (m *MetaDataMusica) GetGenero() string {
	return m.genero
}

func (m *MetaDataMusica) GetTituloCancion() string {
	return m.tituloCancion
}

func (m *MetaDataMusica) GetSelloDiscografico() string {
	return m.selloDiscografico
}

func (m *MetaDataMusica) GetAnioLanzamiento() string {
	return m.anioLanzamiento
}

// --- Setters ---

func (m *MetaDataMusica) SetId(id int) {
	m.id = id
}

func (m *MetaDataMusica) SetTipo(tipo TipoAudio) {
	m.tipo = tipo
}

func (m *MetaDataMusica) SetArtistaPrincipal(artistaPrincipal string) {
	m.artistaPrincipal = artistaPrincipal
}

func (m *MetaDataMusica) SetAlbum(album string) {
	m.album = album
}

func (m *MetaDataMusica) SetGenero(genero string) {
	m.genero = genero
}

func (m *MetaDataMusica) SetTituloCancion(tituloCancion string) {
	m.tituloCancion = tituloCancion
}

func (m *MetaDataMusica) SetSelloDiscografico(selloDiscografico string) {
	m.selloDiscografico = selloDiscografico
}

func (m *MetaDataMusica) SetAnioLanzamiento(anioLanzamiento string) {
	m.anioLanzamiento = anioLanzamiento
}

func (m MetaDataMusica) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                int       `json:"id"`
		Tipo              TipoAudio `json:"tipo"`
		ArtistaPrincipal  string    `json:"artista_principal"`
		Album             string    `json:"album"`
		Genero            string    `json:"genero"`
		TituloCancion     string    `json:"titulo_cancion"`
		SelloDiscografico string    `json:"sello_discografico"`
		AnioLanzamiento   string    `json:"anio_lanzamiento"`
	}{
		ID:                m.id,
		Tipo:              m.tipo,
		ArtistaPrincipal:  m.artistaPrincipal,
		Album:             m.album,
		Genero:            m.genero,
		TituloCancion:     m.tituloCancion,
		SelloDiscografico: m.selloDiscografico,
		AnioLanzamiento:   m.anioLanzamiento,
	})
}

func (m *MetaDataMusica) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID                int       `json:"id"`
		Tipo              TipoAudio `json:"tipo"`
		ArtistaPrincipal  string    `json:"artista_principal"`
		Album             string    `json:"album"`
		Genero            string    `json:"genero"`
		TituloCancion     string    `json:"titulo_cancion"`
		SelloDiscografico string    `json:"sello_discografico"`
		AnioLanzamiento   string    `json:"anio_lanzamiento"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.id = aux.ID
	m.tipo = aux.Tipo
	m.artistaPrincipal = aux.ArtistaPrincipal
	m.album = aux.Album
	m.genero = aux.Genero
	m.tituloCancion = aux.TituloCancion
	m.selloDiscografico = aux.SelloDiscografico
	m.anioLanzamiento = aux.AnioLanzamiento

	return nil
}
