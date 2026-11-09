// Configuración del administrador web.
const CONFIG_ADMIN = {
  // URL base del servicio de metadatos EXPUESTO POR EL API GATEWAY (la misma que usa el cliente web).
  // El gateway enruta /api/metadatosaudios/* -> /metadatosaudios/* del ServidorDeMetadatos-Audios.
  URL_BASE: "http://localhost:8080/api/metadatosaudios",

  // Si el servidor no responde en este tiempo, la consulta o el registro de metadatos se considera fallido.
  // La subida del archivo de audio no tiene este límite: depende del tamaño del archivo y de la red.
  TIEMPO_ESPERA_MS: 5000,

  // El servidor de streaming solo entrega este formato (igual que CONFIG_STREAMING.FORMATO del cliente).
  FORMATO_AUDIO: "mp3",

  // Nombre con el que se guarda el archivo en el servidor para un audio del catálogo.
  // DEBE coincidir con CONFIG_STREAMING.archivoDeAudio del cliente (audios/config.js): es el nombre que
  // el cliente le pide al servidor de streaming para reproducir el audio. Si cambia la convención,
  // hay que cambiarla en los dos lugares.
  archivoDeAudio: (idTipo, idAudio) => `cancion${idAudio}.${CONFIG_ADMIN.FORMATO_AUDIO}`,
};
