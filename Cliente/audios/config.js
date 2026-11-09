// Configuración del catálogo de audios (metadatos).
const CONFIG_METADATOS = {
  // URL base del servicio de metadatos EXPUESTO POR EL API GATEWAY.
  // El gateway enruta /api/metadatosaudios/* -> /metadatosaudios/* del ServidorDeMetadatos-Audios.
  URL_BASE: "http://localhost:8080/api/metadatosaudios",

  // Si el servidor no responde en este tiempo, la consulta se considera fallida.
  TIEMPO_ESPERA_MS: 5000,
};

// Configuración de la reproducción por streaming (gRPC-Web, ver streamingAudios/).
const CONFIG_STREAMING = {
  FORMATO: "mp3",

  // Archivo que se pide al servidor de streaming para un audio del catálogo.
  // Igual que el cliente de consola: `${idAudio}.mp3`. Si cambia la convención de nombres
  // (p. ej. un archivo por tipo), es lo único que hay que tocar.
  archivoDeAudio: (idTipo, idAudio) => `cancion${idAudio}.${CONFIG_STREAMING.FORMATO}`,
};
