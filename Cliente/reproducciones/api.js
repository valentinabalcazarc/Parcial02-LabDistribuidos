// Acceso al servicio REST de estadísticas de reproducciones (a través del API Gateway).

/**
 * Error al consultar el servicio de estadísticas.
 * `reintentable` es true cuando tiene sentido volver a intentar (red caída, tiempo agotado, 5xx).
 */
class ErrorEstadisticas extends Error {
  constructor(mensaje, { estado = 0, reintentable = false } = {}) {
    super(mensaje);
    this.name = "ErrorEstadisticas";
    this.estado = estado;
    this.reintentable = reintentable;
  }
}

const ApiEstadisticas = (() => {
  /** GET a `URL_BASE + ruta`; devuelve el JSON o lanza ErrorEstadisticas. */
  async function consultar(ruta) {
    const controlador = new AbortController();
    const temporizador = setTimeout(() => controlador.abort(), CONFIG_REPRODUCCIONES.TIEMPO_ESPERA_MS);

    let respuesta;
    let datos;
    try {
      respuesta = await fetch(CONFIG_REPRODUCCIONES.URL_BASE + ruta, { signal: controlador.signal });
      datos = await respuesta.json().catch(() => ({}));
    } catch (error) {
      if (error.name === "AbortError") {
        throw new ErrorEstadisticas(
          `El servidor no respondió en ${CONFIG_REPRODUCCIONES.TIEMPO_ESPERA_MS / 1000} s.`,
          { reintentable: true }
        );
      }
      throw new ErrorEstadisticas(
        "No fue posible conectar con el servidor. Verifique que el API Gateway y el servidor de estadísticas estén activos.",
        { reintentable: true }
      );
    } finally {
      clearTimeout(temporizador);
    }

    if (!respuesta.ok) {
      throw new ErrorEstadisticas(datos.mensaje || datos.error || `Error HTTP ${respuesta.status}`, {
        estado: respuesta.status,
        reintentable: respuesta.status >= 500,
      });
    }
    return datos;
  }

  return {
    /**
     * GET /cliente/:cliente -> { codigo, mensaje, datos: { cliente, total_reproducciones, audios: [...] } }
     * Obtiene los audios reproducidos por un cliente y la cantidad de veces.
     */
    async obtenerReproduccionesCliente(clienteId) {
      const datos = await consultar(`/cliente/${encodeURIComponent(clienteId)}`);
      return datos.datos || { cliente: clienteId, total_reproducciones: 0, audios: [] };
    },
  };
})();

