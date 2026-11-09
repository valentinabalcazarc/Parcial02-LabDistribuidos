// reproducciones/historial.js
// Módulo auxiliar que registra el identificador del cliente para las estadísticas de reproducción.
// El cliente de streaming (cliente_receptor_fragmentos.js) envía este mismo ID como metadata
// gRPC 'x-client-id', y la página de reproducciones lo usa para consultar las estadísticas
// del cliente automáticamente, sin necesidad de ingresarlo manualmente.

const HISTORIAL_REPRODUCCIONES = (() => {
  const CLAVE_STORAGE = "reproducciones_cliente_id";

  /**
   * Obtiene el ID del cliente almacenado. Si no existe, genera uno nuevo y lo persiste.
   * Este ID es el mismo que envía el cliente de streaming al servidor via gRPC metadata.
   */
  function obtenerClienteId() {
    let id = localStorage.getItem(CLAVE_STORAGE);
    if (!id) {
      id = 'cliente-web-' + Math.random().toString(36).substring(2, 10);
      localStorage.setItem(CLAVE_STORAGE, id);
    }
    return id;
  }

  /** Guarda un ID de cliente personalizado. */
  function establecerClienteId(id) {
    localStorage.setItem(CLAVE_STORAGE, id);
  }

  return { obtenerClienteId, establecerClienteId };
})();

