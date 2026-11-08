// Configuración del cliente.
const CONFIG = {
  // URL base de la pasarela de pagos (proyecto servidor).
  URL_BASE: "http://localhost:8080/api",

  // Si el servidor no responde en este tiempo, el intento se considera fallido.
  TIEMPO_ESPERA_MS: 3000,

  // Esperas antes de cada reintento (4 reintentos con backoff exponencial).
  // Reintento 1 a los 2 s, reintento 2 a los 4 s, reintento 3 a los 8 s, reintento 4 a los 16 s.
  ESPERAS_REINTENTO_MS: [2000, 4000, 8000, 16000],
};
