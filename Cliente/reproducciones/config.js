// Configuración de la página de reproducciones (estadísticas).
const CONFIG_REPRODUCCIONES = {
  // URL base del servicio de estadísticas EXPUESTO POR EL API GATEWAY.
  // El gateway enruta /api/estadisticas/* -> /estadisticas/* del ServidorDeEstadisticas.
  URL_BASE: "http://localhost:8080/api/estadisticas",

  // Si el servidor no responde en este tiempo, la consulta se considera fallida.
  TIEMPO_ESPERA_MS: 5000,
};

