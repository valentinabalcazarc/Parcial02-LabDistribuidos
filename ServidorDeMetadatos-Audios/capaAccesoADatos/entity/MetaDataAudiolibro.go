package entity

import "encoding/json"

/**
 * @brief Estructura que representa la metadata de un audiolibro.
 * 
 * Contiene información detallada sobre un audiolibro, incluyendo su
 * título, autor, narrador, editorial, etc.
 */
type MetaDataAudiolibro struct {
	id          int
	tipo        TipoAudio
	tituloLibro string
	autor       string
	narrador    string
	editorial   string
	isbn        int
	capitulo    int
}

// --- Getters ---

func (a *MetaDataAudiolibro) GetId() int {
	return a.id
}

func (a *MetaDataAudiolibro) GetTipo() TipoAudio {
	return a.tipo
}

func (a *MetaDataAudiolibro) GetTituloLibro() string {
	return a.tituloLibro
}

func (a *MetaDataAudiolibro) GetAutor() string {
	return a.autor
}

func (a *MetaDataAudiolibro) GetNarrador() string {
	return a.narrador
}

func (a *MetaDataAudiolibro) GetEditorial() string {
	return a.editorial
}

func (a *MetaDataAudiolibro) GetISBN() int {
	return a.isbn
}

func (a *MetaDataAudiolibro) GetCapitulo() int {
	return a.capitulo
}

// --- Setters ---

func (a *MetaDataAudiolibro) SetId(id int) {
	a.id = id
}

func (a *MetaDataAudiolibro) SetTipo(tipo TipoAudio) {
	a.tipo = tipo
}

func (a *MetaDataAudiolibro) SetTituloLibro(tituloLibro string) {
	a.tituloLibro = tituloLibro
}

func (a *MetaDataAudiolibro) SetAutor(autor string) {
	a.autor = autor
}

func (a *MetaDataAudiolibro) SetNarrador(narrador string) {
	a.narrador = narrador
}

func (a *MetaDataAudiolibro) SetEditorial(editorial string) {
	a.editorial = editorial
}

func (a *MetaDataAudiolibro) SetISBN(isbn int) {
	a.isbn = isbn
}

func (a *MetaDataAudiolibro) SetCapitulo(capitulo int) {
	a.capitulo = capitulo
}

// --- JSON Marshaling ---

func (a MetaDataAudiolibro) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          int       `json:"id"`
		Tipo        TipoAudio `json:"tipo"`
		TituloLibro string    `json:"titulo_libro"`
		Autor       string    `json:"autor"`
		Narrador    string    `json:"narrador"`
		Editorial   string    `json:"editorial"`
		ISBN        int       `json:"isbn"`
		Capitulo    int       `json:"capitulo"`
	}{
		ID:          a.id,
		Tipo:        a.tipo,
		TituloLibro: a.tituloLibro,
		Autor:       a.autor,
		Narrador:    a.narrador,
		Editorial:   a.editorial,
		ISBN:        a.isbn,
		Capitulo:    a.capitulo,
	})
}

func (a *MetaDataAudiolibro) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID          int       `json:"id"`
		Tipo        TipoAudio `json:"tipo"`
		TituloLibro string    `json:"titulo_libro"`
		Autor       string    `json:"autor"`
		Narrador    string    `json:"narrador"`
		Editorial   string    `json:"editorial"`
		ISBN        int       `json:"isbn"`
		Capitulo    int       `json:"capitulo"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.id = aux.ID
	a.tipo = aux.Tipo
	a.tituloLibro = aux.TituloLibro
	a.autor = aux.Autor
	a.narrador = aux.Narrador
	a.editorial = aux.Editorial
	a.isbn = aux.ISBN
	a.capitulo = aux.Capitulo

	return nil
}
