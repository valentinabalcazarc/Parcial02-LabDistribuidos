// Descripción de cada tipo de audio soportado: endpoints REST y campos de metadatos.
// Los ids son los del servidor (TipoAudioRepository) y los mismos que usa el cliente de consola.
// Cada tipo responde con claves distintas, por eso se declaran aquí:
//   - rutaLista / claveLista:     GET de la lista y propiedad del JSON que contiene el arreglo.
//   - rutaDetalle / claveDetalle: GET de un audio y propiedad del JSON que contiene el objeto.
//   - titulo / subtitulo:         qué mostrar de cada audio en la lista y en la ficha.
//   - campos:                     [etiqueta, propiedad, formateador opcional] de la ficha técnica.

/** La duración del ruido blanco llega en segundos. */
function formatearDuracion(segundos) {
  if (typeof segundos !== "number") return segundos;
  const minutos = Math.floor(segundos / 60);
  const resto = segundos % 60;
  if (minutos === 0) return `${resto} s`;
  return resto === 0 ? `${minutos} min` : `${minutos} min ${resto} s`;
}

const TIPOS_AUDIO = {
  1: {
    icono: "🎵",
    rutaLista: "/musicas",
    claveLista: "musica",
    rutaDetalle: (id) => `/musica/${id}`,
    claveDetalle: "objMusica",
    titulo: (a) => a.titulo_cancion,
    subtitulo: (a) => a.artista_principal,
    campos: [
      ["Título de la canción", "titulo_cancion"],
      ["Artista principal", "artista_principal"],
      ["Álbum", "album"],
      ["Género", "genero"],
      ["Sello discográfico", "sello_discografico"],
      ["Año de lanzamiento", "anio_lanzamiento"],
    ],
  },

  2: {
    icono: "🌊",
    rutaLista: "/ruidosblancos",
    claveLista: "ruidosBlancos",
    rutaDetalle: (id) => `/ruidoblanco/${id}`,
    claveDetalle: "objRuidoBlanco",
    titulo: (a) => a.tipo_sonido,
    subtitulo: (a) => a.fuente_audio,
    campos: [
      ["Tipo de sonido", "tipo_sonido"],
      ["Fuente del audio", "fuente_audio"],
      ["Uso sugerido", "uso_sugerido"],
      ["Proveedor", "proveedor"],
      ["Duración", "duracion", formatearDuracion],
      ["Frecuencia dominante", "frecuencia_dominante"],
    ],
  },

  3: {
    icono: "🎙️",
    rutaLista: "/podcasts",
    claveLista: "podcasts",
    rutaDetalle: (id) => `/podcast/${id}`,
    claveDetalle: "objPodcast",
    titulo: (a) => a.nombre,
    subtitulo: (a) => a.titulo_episodio,
    campos: [
      ["Nombre del podcast", "nombre"],
      ["Título del episodio", "titulo_episodio"],
      ["Número de temporada", "numero_temporada"],
      ["Notas del show", "notas_show"],
      ["Clasificación del contenido", "clasificacion_contenido"],
    ],
  },

  4: {
    icono: "📖",
    rutaLista: "/audiolibros",
    claveLista: "audiolibros",
    rutaDetalle: (id) => `/audiolibro/${id}`,
    claveDetalle: "objAudiolibro",
    titulo: (a) => a.titulo_libro,
    subtitulo: (a) => a.autor,
    campos: [
      ["Título del libro", "titulo_libro"],
      ["Autor", "autor"],
      ["Narrador", "narrador"],
      ["Editorial", "editorial"],
      ["ISBN", "isbn"],
      ["Capítulo", "capitulo"],
    ],
  },
};
