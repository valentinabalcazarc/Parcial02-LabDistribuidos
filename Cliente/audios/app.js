// Lógica de la interfaz del catálogo de audios: consulta de metadatos y reproducción por streaming.
//
// Flujo (el mismo del cliente de consola):
//   tipos de audio  ->  lista de audios del tipo  ->  ficha técnica del audio  ->  reproducir (streaming)
//
// La vista actual vive en el hash de la URL, así el botón "atrás" del navegador
// y recargar la página funcionan:
//   #/                      tipos de audio
//   #/tipos/1               audios del tipo 1
//   #/tipos/1/audios/3      ficha técnica del audio 3 del tipo 1, con su reproductor

const $ = (id) => document.getElementById(id);

const estado = {
  tipos: null, // caché de los tipos de audio (se piden una sola vez)
  navegacion: 0, // contador para descartar respuestas de navegaciones ya superadas
};

/** Crea un elemento DOM. Los textos se asignan con textContent (sin inyección de HTML). */
function el(etiqueta, atributos = {}, hijos = []) {
  const nodo = document.createElement(etiqueta);
  for (const [clave, valor] of Object.entries(atributos)) {
    if (clave === "texto") nodo.textContent = valor;
    else if (clave === "clase") nodo.className = valor;
    else if (clave.startsWith("on")) nodo.addEventListener(clave.slice(2), valor);
    else nodo.setAttribute(clave, valor);
  }
  for (const hijo of [].concat(hijos)) {
    if (hijo != null) nodo.append(hijo);
  }
  return nodo;
}

const hrefTipos = () => "#/";
const hrefLista = (idTipo) => `#/tipos/${idTipo}`;
const hrefDetalle = (idTipo, idAudio) => `#/tipos/${idTipo}/audios/${idAudio}`;

/** Muestra "—" cuando el servidor no trae el dato. Los números se muestran tal cual (p. ej. el ISBN). */
const textoValor = (valor) => (valor === null || valor === undefined || valor === "" ? "—" : String(valor));

/** Nombre a mostrar de un audio; si el servidor no lo trae, "Audio #id" (igual que el cliente de consola). */
const nombreAudio = (definicion, audio) => definicion.titulo(audio) || `Audio #${audio.id}`;

// ---------------------------------------------------------------- Rutas

function rutaActual() {
  const partes = location.hash.replace(/^#\/?/, "").split("/").filter(Boolean);
  const entero = (texto) => (/^\d+$/.test(texto ?? "") ? Number(texto) : null);

  if (partes[0] === "tipos" && entero(partes[1]) !== null) {
    if (partes.length === 2) return { vista: "lista", idTipo: entero(partes[1]) };
    if (partes.length === 4 && partes[2] === "audios" && entero(partes[3]) !== null) {
      return { vista: "detalle", idTipo: entero(partes[1]), idAudio: entero(partes[3]) };
    }
  }
  return { vista: "tipos" };
}

async function obtenerTipos() {
  if (!estado.tipos) estado.tipos = await ApiMetadatos.listarTipos();
  return estado.tipos;
}

async function enrutar() {
  const numero = ++estado.navegacion;
  const vigente = () => numero === estado.navegacion;
  const ruta = rutaActual();

  detenerReproduccion(); // al cambiar de vista (o reintentar) no debe seguir sonando el audio anterior
  limpiarReacciones();   // desconectar reacciones al cambiar de vista
  pintarCargando();
  try {
    const tipos = await obtenerTipos();
    if (!vigente()) return;

    if (ruta.vista === "tipos") {
      pintarTipos(tipos);
      return;
    }

    const tipo = tipos.find((t) => t.id === ruta.idTipo);
    if (!tipo) {
      pintarError("Tipo de audio no encontrado.", { volver: { href: hrefTipos(), texto: "Volver a los tipos" } });
      return;
    }
    const definicion = TIPOS_AUDIO[tipo.id];
    if (!definicion) {
      pintarError(`Tipo de audio #${tipo.id} no soportado.`, { volver: { href: hrefTipos(), texto: "Volver a los tipos" } });
      return;
    }

    if (ruta.vista === "lista") {
      const audios = await ApiMetadatos.listarAudios(tipo.id);
      if (!vigente()) return;
      pintarLista(tipo, definicion, audios);
      return;
    }

    const audio = await ApiMetadatos.consultarAudio(tipo.id, ruta.idAudio);
    if (!vigente()) return;
    pintarDetalle(tipo, definicion, audio);
  } catch (error) {
    if (!vigente()) return;
    const volver =
      ruta.vista === "detalle"
        ? { href: hrefLista(ruta.idTipo), texto: "Volver a la lista" }
        : ruta.vista === "lista"
          ? { href: hrefTipos(), texto: "Volver a los tipos" }
          : null;
    pintarError(error.message, { volver, reintentar: error.reintentable !== false });
  }
}

// ---------------------------------------------------------------- Vistas

/** Reemplaza el contenido de la página y lleva el foco al título (útil para lectores de pantalla). */
function mostrar(tituloPagina, nodos) {
  document.title = `${tituloPagina} · Catálogo de Audios`;
  $("contenido").replaceChildren(...nodos);
  $("titulo-vista")?.focus({ preventScroll: true });
}

function encabezado({ titulo, subtitulo, volver }) {
  return el("div", { clase: "space-y-3" }, [
    volver
      ? el("a", {
        href: volver.href,
        clase: "inline-flex items-center rounded-full border border-gray-300 px-4 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors",
        texto: `← ${volver.texto}`,
      })
      : null,
    el("h2", { id: "titulo-vista", tabindex: "-1", clase: "text-2xl font-bold text-gray-900 outline-none", texto: titulo }),
    subtitulo ? el("p", { clase: "text-sm text-gray-500", texto: subtitulo }) : null,
  ]);
}

function pintarCargando() {
  $("contenido").replaceChildren(el("p", { clase: "texto-tenue", texto: "Cargando…" }));
}

function pintarError(mensaje, { volver = null, reintentar = false } = {}) {
  mostrar("Error", [
    el("div", { role: "alert", clase: "rounded-2xl border border-red-200 bg-red-50 p-6 space-y-4" }, [
      el("h2", { id: "titulo-vista", tabindex: "-1", clase: "text-lg font-bold text-red-700 outline-none", texto: "No fue posible completar la consulta" }),
      el("p", { clase: "text-sm text-red-600", texto: mensaje }),
      el("div", { clase: "flex flex-wrap gap-3" }, [
        reintentar
          ? el("button", {
            type: "button",
            clase: "rounded-full bg-accent hover:bg-accent-hover active:brightness-90 px-6 py-2 text-sm font-semibold text-white transition-all shadow-sm",
            texto: "↻ Reintentar",
            onclick: enrutar,
          })
          : null,
        volver
          ? el("a", {
            href: volver.href,
            clase: "inline-flex items-center rounded-full border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors",
            texto: `← ${volver.texto}`,
          })
          : null,
      ]),
    ]),
  ]);
}

function pintarTipos(tipos) {
  const tarjetas =
    tipos.length === 0
      ? el("p", { clase: "texto-tenue", texto: "No hay tipos de audio disponibles." })
      : el(
        "div",
        { clase: "grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4" },
        tipos.map((tipo) =>
          el("a", { href: hrefLista(tipo.id), "data-tipo": tipo.id, clase: "tarjeta bg-white rounded-2xl border border-gray-200 shadow-sm p-6 flex flex-col gap-3 hover:shadow-md transition-shadow" }, [
            el("span", { clase: "text-3xl", "aria-hidden": "true", texto: TIPOS_AUDIO[tipo.id]?.icono || "🎧" }),
            el("h3", { clase: "text-lg font-bold text-gray-900", texto: tipo.nombre }),
            el("span", { clase: "text-sm font-semibold text-accent", texto: "Ver audios →" }),
          ])
        )
      );

  mostrar("Tipos de audio", [
    encabezado({ titulo: "Tipos de audio", subtitulo: "Seleccione un tipo para ver los audios disponibles." }),
    tarjetas,
  ]);
}

function pintarLista(tipo, definicion, audios) {
  const contenido =
    audios.length === 0
      ? el("p", { clase: "texto-tenue", texto: "No hay audios registrados en este tipo." })
      : el(
        "div",
        { clase: "bg-white rounded-2xl border border-gray-200 shadow-sm divide-y divide-gray-100 overflow-hidden" },
        audios.map((audio) => {
          const subtitulo = definicion.subtitulo(audio);
          return el("a", { href: hrefDetalle(tipo.id, audio.id), "data-audio": audio.id, clase: "fila-audio flex items-center gap-4 px-5 py-4 transition-colors" }, [
            el("span", { clase: "shrink-0 rounded-full bg-accent-light text-accent text-xs font-semibold px-3 py-1", texto: `ID ${audio.id}` }),
            el("div", { clase: "min-w-0 flex-1" }, [
              el("p", { clase: "font-semibold text-gray-900 truncate", texto: nombreAudio(definicion, audio) }),
              subtitulo ? el("p", { clase: "text-sm text-gray-500 truncate", texto: subtitulo }) : null,
            ]),
            el("span", { clase: "shrink-0 text-sm font-medium text-accent", texto: "Ver detalles →" }),
          ]);
        })
      );

  mostrar(tipo.nombre, [
    encabezado({
      titulo: `${definicion.icono} ${tipo.nombre}`,
      subtitulo: audios.length === 1 ? "1 audio disponible. Seleccione uno para ver sus detalles." : `${audios.length} audios disponibles. Seleccione uno para ver sus detalles.`,
      volver: { href: hrefTipos(), texto: "Volver a los tipos" },
    }),
    contenido,
  ]);
}

// ---------------------------------------------------------------- Streaming y reacciones (estado global)

/** Transmisión en curso ({ detener() } de streamingAudios/); null si no hay ninguna. Solo hay una a la vez. */
let transmisionActiva = null;

/** Sección de reacciones del audio que se está reproduciendo. */
let seccionReaccionesActual = null;

/** Detiene la recepción y la reproducción en curso, si las hay. */
function detenerReproduccion() {
  if (!transmisionActiva) return;
  transmisionActiva.detener();
  transmisionActiva = null;
  // Al detener la reproducción se desactivan las reacciones (ya no se reproduce)
  if (seccionReaccionesActual && seccionReaccionesActual._reacciones) {
    seccionReaccionesActual._reacciones.desactivar();
  }
}

/** Al cambiar de vista, desconectar también las reacciones. */
function limpiarReacciones() {
  if (typeof Reacciones !== "undefined") {
    Reacciones.desconectar();
  }
  seccionReaccionesActual = null;
}

const GRPC_NO_DISPONIBLE = 14;

/** grpc-web informa "http status code: 0" cuando no hubo respuesta alguna (gateway caído); 14 cuando el gateway responde 503. */
function mensajeErrorStreaming(error) {
  const sinConexion = error.codigo === GRPC_NO_DISPONIBLE || /http status code: 0/.test(error.mensaje || "");
  if (sinConexion || !error.mensaje) {
    return "No fue posible conectar con el servidor de streaming. Verifique que el API Gateway y el servidor estén activos.";
  }
  return error.mensaje; // p. ej. el servidor no encontró el archivo de audio
}

function textoTamanio(bytes) {
  return bytes < 1024 * 1024 ? `${Math.round(bytes / 1024)} KB` : `${(bytes / 1024 / 1024).toFixed(1)} MB`;
}

/**
 * Sección "Reproducción" de la ficha: botón para reproducir el audio mediante streaming, estado de la
 * transmisión y botón para detenerla (en la consola: ENTER). Pide el mismo archivo que el cliente de
 * consola, `${idAudio}.mp3` (ver CONFIG_STREAMING).
 */
function crearReproductor(tipo, definicion, audio) {
  const archivo = CONFIG_STREAMING.archivoDeAudio(tipo.id, audio.id);
  const nombre = nombreAudio(definicion, audio);

  const reproductor = el("audio", { id: "audio-player", controls: "", clase: "hidden w-full" });
  const estadoTexto = el("p", { role: "status", clase: "text-sm text-gray-500", texto: "Listo para reproducir." });
  // El avance por fragmentos cambia muy seguido: va aparte y oculto a lectores de pantalla.
  const progreso = el("p", { "aria-hidden": "true", clase: "text-xs text-gray-400 tabular-nums" });
  const botonReproducir = el("button", {
    type: "button",
    clase: "rounded-full bg-accent hover:bg-accent-hover active:brightness-90 px-6 py-2 text-sm font-semibold text-white transition-all shadow-sm",
    texto: "▶️ Reproducir mediante streaming",
    onclick: reproducir,
  });
  const botonDetener = el("button", {
    type: "button",
    clase: "hidden rounded-full border border-gray-300 bg-white px-6 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors",
    texto: "⏹️ Detener",
    onclick: detener,
  });

  /** `activo`: hay una transmisión o reproducción en curso (se ofrece "Detener" en lugar de "Reproducir"). */
  function actualizarEstado(texto, { activo = false, error = false, textoBoton = null } = {}) {
    estadoTexto.textContent = texto;
    estadoTexto.className = `text-sm ${error ? "text-red-600" : "text-gray-500"}`;
    botonReproducir.classList.toggle("hidden", activo);
    botonDetener.classList.toggle("hidden", !activo);
    if (textoBoton) botonReproducir.textContent = textoBoton;
  }

  function fallar(mensaje) {
    reproductor.classList.add("hidden");
    actualizarEstado(`❌ ${mensaje}`, { error: true, textoBoton: "↻ Reintentar" });
  }

  function reproducir() {
    detenerReproduccion(); // una sola transmisión a la vez
    progreso.textContent = "";
    actualizarEstado("📡 Conectando con el servidor de streaming…", { activo: true });

    let fragmentos = 0;
    let bytes = 0;
    const transmision =
      typeof pedirCancion === "function"
        ? pedirCancion(archivo, CONFIG_STREAMING.FORMATO, {
          reproductor,
          onFragmento: (datos) => {
            fragmentos += 1;
            bytes += datos.length;
            if (fragmentos === 1) actualizarEstado("📥 Recibiendo audio del servidor…", { activo: true });
            progreso.textContent = `Fragmento #${fragmentos} recibido (${datos.length} bytes) · ${textoTamanio(bytes)} en total`;
          },
          onRecibido: () => {
            progreso.textContent = `${fragmentos} fragmentos · ${textoTamanio(bytes)}`;
            actualizarEstado("📥 Audio recibido completo desde el servidor.", { activo: true });
          },
          onError: (error) => {
            transmisionActiva = null;
            fallar(mensajeErrorStreaming(error));
          },
        })
        : null;

    if (!transmision) {
      fallar("No se cargó el cliente de streaming (falta bundle.js o streamingAudios/funciones.js).");
      return;
    }
    transmisionActiva = transmision;
  }

  function detener() {
    detenerReproduccion();
    reproductor.classList.add("hidden");
    progreso.textContent = "";
    actualizarEstado("⏹️ Reproducción detenida por el usuario.", { textoBoton: "▶️ Reproducir de nuevo" });
  }

  // Eventos del propio <audio>: empieza a sonar, termina, o el navegador no pudo decodificarlo.
  reproductor.onplaying = () => {
    actualizarEstado(`▶️ Reproduciendo ${nombre}…`, { activo: true });

    // Activar reacciones cuando el audio empieza a sonar
    if (seccionReaccionesActual && seccionReaccionesActual._reacciones) {
      seccionReaccionesActual._reacciones.activar();
    }
  };
  reproductor.onended = () => {
    actualizarEstado("✅ Audio finalizado.", { textoBoton: "▶️ Reproducir de nuevo" });
    // Desactivar reacciones cuando el audio termina
    if (seccionReaccionesActual && seccionReaccionesActual._reacciones) {
      seccionReaccionesActual._reacciones.desactivar();
    }
  };
  reproductor.onerror = () => {
    if (reproductor.hasAttribute("src")) fallar("El navegador no pudo reproducir el audio recibido.");
  };

  return el("section", { clase: "bg-white rounded-2xl border border-gray-200 shadow-sm p-6 space-y-4" }, [
    el("h3", { clase: "text-sm font-bold uppercase tracking-wide text-gray-500", texto: "Reproducción" }),
    reproductor,
    estadoTexto,
    progreso,
    el("div", { clase: "flex flex-wrap gap-3" }, [botonReproducir, botonDetener]),
  ]);
}

// ---------------------------------------------------------------- Vistas de detalle

function pintarDetalle(tipo, definicion, audio) {
  const etiqueta = (texto) => el("span", { clase: "rounded-full bg-accent-light text-accent text-xs font-semibold px-3 py-1", texto });

  const filas = definicion.campos.flatMap(([nombre, clave, formato]) => [
    el("dt", { texto: nombre }),
    el("dd", { texto: textoValor(formato ? formato(audio[clave]) : audio[clave]) }),
  ]);

  // Crear la sección de reacciones (se activa al reproducir)
  const idAudioReacciones = String(audio.id);
  const seccionReacciones = (typeof Reacciones !== "undefined")
    ? Reacciones.crearSeccion(idAudioReacciones)
    : null;
  seccionReaccionesActual = seccionReacciones;

  mostrar(nombreAudio(definicion, audio), [
    encabezado({
      titulo: nombreAudio(definicion, audio),
      subtitulo: definicion.subtitulo(audio),
      volver: { href: hrefLista(tipo.id), texto: "Volver a la lista" },
    }),
    el("section", { clase: "bg-white rounded-2xl border border-gray-200 shadow-sm p-6 space-y-4" }, [
      el("div", { clase: "flex flex-wrap items-center gap-2" }, [
        el("h3", { clase: "text-sm font-bold uppercase tracking-wide text-gray-500 mr-2", texto: "Ficha técnica" }),
        etiqueta(`ID ${audio.id}`),
        etiqueta(audio.tipo?.nombre || tipo.nombre),
      ]),
      el("dl", { clase: "ficha" }, filas),
    ]),
    crearReproductor(tipo, definicion, audio),
    seccionReacciones,
  ].filter(Boolean));
}

// ---------------------------------------------------------------- Inicio

document.addEventListener("DOMContentLoaded", () => {
  window.addEventListener("hashchange", enrutar);
  enrutar();
});
