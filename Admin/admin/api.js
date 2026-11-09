
class ErrorAdmin extends Error {
  constructor(mensaje, { estado = 0, reintentable = false } = {}) {
    super(mensaje);
    this.name = "ErrorAdmin";
    this.estado = estado;
    this.reintentable = reintentable;
  }
}

const ApiAdmin = (() => {
  const url = (ruta) => CONFIG_ADMIN.URL_BASE + ruta;

  const MENSAJE_SIN_CONEXION =
    "No fue posible conectar con el servidor. Verifique que el API Gateway y el servidor de metadatos estén activos.";

  /** El servidor de metadatos responde { codigo, mensaje }; el gateway { error }. */
  function errorHttp(estado, datos) {
    return new ErrorAdmin(datos?.mensaje || datos?.error || `Error HTTP ${estado}`, {
      estado,
      reintentable: estado >= 500,
    });
  }

  /** `senal` es una AbortSignal opcional para cancelar desde afuera (p. ej. el usuario salió de la vista). */
  async function solicitar(ruta, { metodo = "GET", json = null, senal = null } = {}) {
    const controlador = new AbortController();
    let agotado = false;
    const temporizador = setTimeout(() => {
      agotado = true;
      controlador.abort();
    }, CONFIG_ADMIN.TIEMPO_ESPERA_MS);
    const cancelar = () => controlador.abort();
    senal?.addEventListener("abort", cancelar, { once: true });
    if (senal?.aborted) cancelar();

    let respuesta;
    let datos;
    try {
      respuesta = await fetch(url(ruta), {
        method: metodo,
        headers: json ? { "Content-Type": "application/json" } : undefined,
        body: json ? JSON.stringify(json) : undefined,
        signal: controlador.signal,
      });
      datos = await respuesta.json().catch(() => ({}));
    } catch (error) {
      if (senal?.aborted) throw error; // cancelado desde afuera: no es un fallo del servidor
      if (agotado) {
        throw new ErrorAdmin(`El servidor no respondió en ${CONFIG_ADMIN.TIEMPO_ESPERA_MS / 1000} s.`, { reintentable: true });
      }
      throw new ErrorAdmin(MENSAJE_SIN_CONEXION, { reintentable: true });
    } finally {
      clearTimeout(temporizador);
      senal?.removeEventListener("abort", cancelar);
    }

    if (!respuesta.ok) throw errorHttp(respuesta.status, datos);
    return datos;
  }

  function definicion(idTipo) {
    const tipo = TIPOS_AUDIO[idTipo];
    if (!tipo) throw new ErrorAdmin(`Tipo de audio #${idTipo} no soportado.`);
    return tipo;
  }

  /** Antepone el contexto del paso que falló al mensaje del error (los errores ajenos, como AbortError, pasan sin tocar). */
  async function paso(contexto, accion, consejo = () => "") {
    try {
      return await accion();
    } catch (error) {
      if (!(error instanceof ErrorAdmin)) throw error;
      throw new ErrorAdmin([contexto, error.message, consejo(error)].filter(Boolean).join(" "), {
        estado: error.estado,
        reintentable: error.reintentable,
      });
    }
  }

  /** Mayor ID de todo el catálogo + 1. Los cuatro tipos comparten la numeración porque el archivo no lleva el tipo. */
  async function siguienteId(senal) {
    const listas = await Promise.all(
      Object.values(TIPOS_AUDIO).map(async (tipo) => {
        const datos = await solicitar(tipo.rutaLista, { senal });
        return datos[tipo.claveLista] || []; // el servidor devuelve `null` (no []) cuando no hay ninguno
      })
    );
    const mayor = listas.flat().reduce((maximo, audio) => Math.max(maximo, Number(audio.id) || 0), 0);
    return mayor + 1;
  }

  /**
   * Sube el archivo de audio. Usa XMLHttpRequest (y no fetch) porque es la única forma de informar
   * el avance de la subida. Campos que espera Audio_controller: "archivo" y "nombre_archivo".
   */
  function subirAudio(nombreArchivo, archivo, { senal = null, alProgresar = () => { } } = {}) {
    return new Promise((resolver, rechazar) => {
      if (senal?.aborted) return rechazar(new DOMException("Subida cancelada", "AbortError"));

      const cuerpo = new FormData();
      cuerpo.append("archivo", archivo, archivo.name);
      cuerpo.append("nombre_archivo", nombreArchivo);

      const xhr = new XMLHttpRequest();
      const cancelar = () => xhr.abort();
      const terminar = () => senal?.removeEventListener("abort", cancelar);

      xhr.open("POST", url("/audio"));
      xhr.upload.onprogress = (evento) => {
        if (evento.lengthComputable) alProgresar(evento.loaded / evento.total);
      };
      xhr.onload = () => {
        terminar();
        let datos = {};
        try {
          datos = JSON.parse(xhr.responseText);
        } catch {
          /* respuesta sin JSON: se informa solo el código HTTP */
        }
        if (xhr.status >= 200 && xhr.status < 300) resolver(datos);
        else rechazar(errorHttp(xhr.status, datos));
      };
      xhr.onerror = () => {
        terminar();
        rechazar(new ErrorAdmin(MENSAJE_SIN_CONEXION, { reintentable: true }));
      };
      xhr.onabort = () => {
        terminar();
        rechazar(new DOMException("Subida cancelada", "AbortError"));
      };

      senal?.addEventListener("abort", cancelar, { once: true });
      xhr.send(cuerpo);
    });
  }

  return {
    /** GET /tipos-audio -> [{ id, nombre }] */
    async listarTipos() {
      const datos = await solicitar("/tipos-audio");
      return datos.tipos || [];
    },

    /**
     * Registra un audio: reserva un ID, sube el archivo y registra los metadatos.
     * @param idTipo   id del tipo de audio.
     * @param valores  metadatos ya convertidos a los tipos del DTO ({ titulo_cancion: "...", ... }).
     * @param archivo  File con el audio.
     * @param opciones { senal, alEstado(texto), alProgresar(fraccion 0..1) }
     * @returns { audio, nombreArchivo } con el objeto que devolvió el servidor (incluye el id asignado).
     */
    async registrarAudio(idTipo, valores, archivo, { senal = null, alEstado = () => { }, alProgresar = () => { } } = {}) {
      const tipo = definicion(idTipo);

      alEstado("Reservando un ID en el catálogo…");
      const id = await paso("No fue posible consultar el catálogo para asignar el ID del audio.", () => siguienteId(senal));

      const nombreArchivo = CONFIG_ADMIN.archivoDeAudio(idTipo, id);
      alEstado("Subiendo el audio…");
      await paso("No se pudo subir el audio.", () => subirAudio(nombreArchivo, archivo, { senal, alProgresar }));

      alEstado("Registrando los metadatos…");
      const datos = await paso(
        "El audio se subió, pero no se pudieron registrar los metadatos.",
        () => solicitar(tipo.rutaRegistro, { metodo: "POST", json: { id, ...valores }, senal }),
        (error) => {
          if (error.estado === 0) return "Verifique en el catálogo si el registro quedó guardado antes de reintentar.";
          if (error.reintentable) return "Puede volver a intentarlo: el audio ya está en el servidor y se reemplazará.";
          return "";
        }
      );

      const audio = datos[tipo.claveRegistro];
      if (!audio) throw new ErrorAdmin("La respuesta del servidor no contiene el audio registrado.");
      return { audio, nombreArchivo };
    },
  };
})();
