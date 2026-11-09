
/** La duración del ruido blanco se guarda en segundos. */
function formatearDuracion(segundos) {
  if (typeof segundos !== "number") return segundos;
  const minutos = Math.floor(segundos / 60);
  const resto = segundos % 60;
  if (minutos === 0) return `${resto} s`;
  return resto === 0 ? `${minutos} min` : `${minutos} min ${resto} s`;
}

// ---------------------------------------------------------------- Clases de campo
// Cada clase sabe cómo validar el texto escrito y cómo convertirlo al valor que espera el DTO del servidor
// (ojo con los tipos: el año de lanzamiento es string, pero la temporada, el ISBN, el capítulo y la duración son int).

function validarEntero(texto, campo) {
  if (!/^\d+$/.test(texto)) return "Ingrese un número entero, sin signos ni decimales.";
  const numero = Number(texto);
  if (!Number.isSafeInteger(numero)) return "El número es demasiado grande.";
  const minimo = campo.min ?? 0;
  if (numero < minimo) return `Ingrese un número mayor o igual a ${minimo}.`;
  return null;
}

function validarAnio(texto) {
  if (!/^\d{4}$/.test(texto)) return "Ingrese un año de 4 dígitos (por ejemplo, 1982).";
  if (Number(texto) > new Date().getFullYear() + 1) return "El año no puede ser posterior al próximo año.";
  return null;
}

/** Los guiones y espacios del ISBN se ignoran: 978-0-307-47472-8 y 9780307474728 son lo mismo. */
const limpiarIsbn = (texto) => texto.replace(/[\s-]/g, "");

function validarIsbn(texto) {
  const isbn = limpiarIsbn(texto);
  if (!/^\d{1,13}$/.test(isbn)) return "Ingrese solo números (hasta 13 dígitos); los guiones y espacios se ignoran.";
  // El servidor guarda el ISBN como número entero: un cero inicial se perdería sin avisar.
  if (isbn.startsWith("0")) return "El ISBN no puede empezar por 0 (el servidor lo guarda como número). Use el ISBN-13.";
  return null;
}

const TIPOS_CAMPO = {
  texto: { modo: "text", validar: () => null, convertir: (texto) => texto },
  area: { modo: "text", validar: () => null, convertir: (texto) => texto },
  entero: { modo: "numeric", validar: validarEntero, convertir: (texto) => Number(texto) },
  anio: { modo: "numeric", validar: validarAnio, convertir: (texto) => texto },
  isbn: { modo: "numeric", validar: validarIsbn, convertir: (texto) => Number(limpiarIsbn(texto)) },
};

// ---------------------------------------------------------------- Tipos de audio

const TIPOS_AUDIO = {
  1: {
    icono: "🎵",
    rutaLista: "/musicas",
    claveLista: "musica",
    rutaRegistro: "/musica",
    claveRegistro: "objMusica",
    campos: [
      { clave: "titulo_cancion", etiqueta: "Título de la canción", tipo: "texto", obligatorio: true, ancho: "completo", placeholder: "Billie Jean" },
      { clave: "artista_principal", etiqueta: "Artista principal", tipo: "texto", placeholder: "Michael Jackson" },
      { clave: "album", etiqueta: "Álbum", tipo: "texto", placeholder: "Thriller" },
      { clave: "genero", etiqueta: "Género", tipo: "texto", placeholder: "Pop" },
      { clave: "sello_discografico", etiqueta: "Sello discográfico", tipo: "texto", placeholder: "Epic Records" },
      { clave: "anio_lanzamiento", etiqueta: "Año de lanzamiento", tipo: "anio", placeholder: "1982" },
    ],
  },

  2: {
    icono: "🌊",
    rutaLista: "/ruidosblancos",
    claveLista: "ruidosBlancos",
    rutaRegistro: "/ruidoblanco",
    claveRegistro: "objRuidoBlanco",
    campos: [
      {
        clave: "tipo_sonido", etiqueta: "Tipo de sonido", tipo: "texto", obligatorio: true, ancho: "completo",
        placeholder: "Ruido Blanco", sugerencias: ["Ruido Blanco", "Ruido Rosa", "Ruido Marrón"],
      },
      { clave: "fuente_audio", etiqueta: "Fuente del audio", tipo: "texto", placeholder: "Lluvia" },
      {
        clave: "uso_sugerido", etiqueta: "Uso sugerido", tipo: "texto",
        placeholder: "Dormir", sugerencias: ["Dormir", "Concentración", "Relajación"],
      },
      { clave: "proveedor", etiqueta: "Proveedor", tipo: "texto", placeholder: "SleepSounds Studio" },
      {
        clave: "duracion", etiqueta: "Duración (segundos)", tipo: "entero", min: 1, formato: formatearDuracion,
        placeholder: "60", ayuda: "Se completa con la duración del archivo; puede corregirla.",
        autocompletarDesdeArchivo: "duracion",
      },
      {
        clave: "frecuencia_dominante", etiqueta: "Frecuencia dominante", tipo: "texto",
        placeholder: "Graves", sugerencias: ["Graves", "Medios", "Agudos"],
      },
    ],
  },

  3: {
    icono: "🎙️",
    rutaLista: "/podcasts",
    claveLista: "podcasts",
    rutaRegistro: "/podcast",
    claveRegistro: "objPodcast",
    campos: [
      { clave: "nombre", etiqueta: "Nombre del podcast", tipo: "texto", placeholder: "Radio Ambulante" },
      { clave: "titulo_episodio", etiqueta: "Título del episodio", tipo: "texto", obligatorio: true, placeholder: "El Polizón" },
      { clave: "numero_temporada", etiqueta: "Número de temporada", tipo: "entero", min: 1, placeholder: "1" },
      {
        clave: "clasificacion_contenido", etiqueta: "Clasificación del contenido", tipo: "texto",
        placeholder: "Para toda la familia", sugerencias: ["Para toda la familia", "Explícito"],
      },
      { clave: "notas_show", etiqueta: "Notas del show", tipo: "area", ancho: "completo", placeholder: "Breve descripción del episodio" },
    ],
  },

  4: {
    icono: "📖",
    rutaLista: "/audiolibros",
    claveLista: "audiolibros",
    rutaRegistro: "/audiolibro",
    claveRegistro: "objAudiolibro",
    campos: [
      { clave: "titulo_libro", etiqueta: "Título del libro", tipo: "texto", obligatorio: true, ancho: "completo", placeholder: "El Principito" },
      { clave: "autor", etiqueta: "Autor", tipo: "texto", placeholder: "Antoine de Saint-Exupéry" },
      { clave: "narrador", etiqueta: "Narrador", tipo: "texto", placeholder: "José María Carrascal" },
      { clave: "editorial", etiqueta: "Editorial", tipo: "texto", placeholder: "Salamandra" },
      { clave: "isbn", etiqueta: "ISBN", tipo: "isbn", placeholder: "978-84-7888-000-0", ayuda: "Solo números; los guiones se ignoran." },
      { clave: "capitulo", etiqueta: "Capítulo", tipo: "entero", min: 1, placeholder: "1" },
    ],
  },
};
