// Acceso a los servicios REST de metadatos de audios (a través del API Gateway).

/**
 * Error al consultar el servicio de metadatos.
 * `reintentable` es true cuando tiene sentido volver a intentar (red caída, tiempo agotado, 5xx)
 * y false cuando el servidor respondió con un error de negocio (4xx, p. ej. audio no encontrado).
 */
class ErrorMetadatos extends Error {
  constructor(mensaje, { estado = 0, reintentable = false } = {}) {
    super(mensaje);
    this.name = "ErrorMetadatos";
    this.estado = estado;
    this.reintentable = reintentable;
  }
}

const ApiMetadatos = (() => {
  /** GET a `URL_BASE + ruta`; devuelve el JSON o lanza ErrorMetadatos. */
  async function consultar(ruta) {
    const controlador = new AbortController();
    const temporizador = setTimeout(() => controlador.abort(), CONFIG_METADATOS.TIEMPO_ESPERA_MS);

    let respuesta;
    let datos;
    try {
      respuesta = await fetch(CONFIG_METADATOS.URL_BASE + ruta, { signal: controlador.signal });
      datos = await respuesta.json().catch(() => ({}));
    } catch (error) {
      if (error.name === "AbortError") {
        throw new ErrorMetadatos(
          `El servidor no respondió en ${CONFIG_METADATOS.TIEMPO_ESPERA_MS / 1000} s.`,
          { reintentable: true }
        );
      }
      throw new ErrorMetadatos(
        "No fue posible conectar con el servidor. Verifique que el API Gateway y el servidor de metadatos estén activos.",
        { reintentable: true }
      );
    } finally {
      clearTimeout(temporizador);
    }

    if (!respuesta.ok) {
      // El servidor de metadatos responde { codigo, mensaje }; el gateway { error }.
      throw new ErrorMetadatos(datos.mensaje || datos.error || `Error HTTP ${respuesta.status}`, {
        estado: respuesta.status,
        reintentable: respuesta.status >= 500,
      });
    }
    return datos;
  }

  function definicion(idTipo) {
    const tipo = TIPOS_AUDIO[idTipo];
    if (!tipo) throw new ErrorMetadatos(`Tipo de audio #${idTipo} no soportado.`);
    return tipo;
  }

  return {
    /** GET /tipos-audio -> [{ id, nombre }] */
    async listarTipos() {
      const datos = await consultar("/tipos-audio");
      return datos.tipos || [];
    },

    /** Lista de audios de un tipo. El servidor devuelve `null` (no []) cuando no hay ninguno. */
    async listarAudios(idTipo) {
      const tipo = definicion(idTipo);
      const datos = await consultar(tipo.rutaLista);
      return datos[tipo.claveLista] || [];
    },

    /** Metadatos completos de un audio. */
    async consultarAudio(idTipo, idAudio) {
      const tipo = definicion(idTipo);
      const datos = await consultar(tipo.rutaDetalle(idAudio));
      const audio = datos[tipo.claveDetalle];
      if (!audio) throw new ErrorMetadatos("La respuesta del servidor no contiene los metadatos del audio.");
      return audio;
    },
  };
})();
