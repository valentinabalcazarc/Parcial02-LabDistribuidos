package entity

import "encoding/json"

/**
 * @brief Entidad que representa un Tipo de Audio.
 *
 * Estructura fundamental para categorizar los diferentes tipos de audio
 * en el sistema.
 */
type TipoAudio struct {
	id     int
	nombre string
}

func NewTipoAudio(id int, nombre string) TipoAudio {
	return TipoAudio{
		id:     id,
		nombre: nombre,
	}
}

// Getters y Setters
func (t *TipoAudio) GetId() int {
	return t.id
}

func (t *TipoAudio) GetNombre() string {
	return t.nombre
}

func (t *TipoAudio) SetId(id int) {
	t.id = id
}

func (t *TipoAudio) SetNombre(nombre string) {
	t.nombre = nombre
}

// JSON Serialization
func (t TipoAudio) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}{
		ID:     t.id,
		Nombre: t.nombre,
	})
}

func (t *TipoAudio) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.id = aux.ID
	t.nombre = aux.Nombre
	return nil
}
