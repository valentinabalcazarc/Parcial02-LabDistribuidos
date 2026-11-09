// Lógica de la interfaz de la página "Mis Reproducciones".
// Muestra tarjetas con los audios que el cliente ha reproducido y cuántas veces.

const $r = (id) => document.getElementById(id);

/** Crea un elemento DOM. Los textos se asignan con textContent (sin inyección de HTML). */
function elR(etiqueta, atributos = {}, hijos = []) {
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

/** Iconos por tipo de audio. */
const ICONOS_TIPO = {
  musica: "🎵",
  Musica: "🎵",
  Música: "🎵",
  podcast: "🎙️",
  Podcast: "🎙️",
  audiolibro: "📖",
  Audiolibro: "📖",
  ruido_blanco: "🌊",
  "Ruido Blanco": "🌊",
  RuidoBlanco: "🌊",
};

function iconoTipo(tipoAudio) {
  return ICONOS_TIPO[tipoAudio] || "🎧";
}

/** Formatea la fecha de última reproducción. */
function formatearFechaReproduccion(fechaStr) {
  if (!fechaStr) return "—";
  try {
    const fecha = new Date(fechaStr);
    if (isNaN(fecha.getTime())) return fechaStr;
    return fecha.toLocaleString("es-CO", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return fechaStr;
  }
}

/** Crea una tarjeta para un audio reproducido. */
function crearTarjetaAudio(audio) {
  const icono = iconoTipo(audio.tipo_audio);
  const tipo = audio.tipo_audio || "Audio";

  return elR("article", {
    clase: "tarjeta-reproduccion bg-white rounded-2xl border border-gray-200 shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow",
  }, [
    // Encabezado: ícono + tipo
    elR("div", { clase: "flex items-center justify-between" }, [
      elR("span", { clase: "text-3xl", "aria-hidden": "true", texto: icono }),
      elR("span", {
        clase: "rounded-full bg-accent-light text-accent text-xs font-semibold px-3 py-1",
        texto: tipo,
      }),
    ]),
    // Título del audio
    elR("h3", {
      clase: "text-base font-bold text-gray-900 truncate",
      title: audio.titulo,
      texto: audio.titulo,
    }),
    // Contador de reproducciones
    elR("div", { clase: "flex items-center gap-2" }, [
      elR("span", { clase: "contador-reproducciones text-2xl font-bold text-accent", texto: String(audio.cantidad_veces) }),
      elR("span", {
        clase: "text-sm text-gray-500",
        texto: audio.cantidad_veces === 1 ? "reproducción" : "reproducciones",
      }),
    ]),
    // Última reproducción
    elR("p", {
      clase: "text-xs text-gray-400 mt-auto",
      texto: `Última vez: ${formatearFechaReproduccion(audio.ultima_reproduccion)}`,
    }),
  ]);
}

/** Muestra un estado de error en el grid. */
function mostrarErrorReproducciones(mensaje, reintentable) {
  const contenedor = $r("lista-reproducciones");
  contenedor.replaceChildren(
    elR("div", {
      role: "alert",
      clase: "col-span-full rounded-2xl border border-red-200 bg-red-50 p-6 space-y-4",
    }, [
      elR("h3", { clase: "text-lg font-bold text-red-700", texto: "No fue posible cargar las reproducciones" }),
      elR("p", { clase: "text-sm text-red-600", texto: mensaje }),
      reintentable
        ? elR("button", {
            type: "button",
            clase: "rounded-full bg-accent hover:bg-accent-hover active:brightness-90 px-6 py-2 text-sm font-semibold text-white transition-all shadow-sm",
            texto: "↻ Reintentar",
            onclick: cargarReproducciones,
          })
        : null,
    ])
  );
}

/** Muestra estado vacío. */
function mostrarVacio() {
  const contenedor = $r("lista-reproducciones");
  contenedor.replaceChildren(
    elR("div", { clase: "col-span-full text-center py-16 space-y-4" }, [
      elR("span", { clase: "text-6xl block", "aria-hidden": "true", texto: "🎧" }),
      elR("h3", { clase: "text-lg font-semibold text-gray-700", texto: "Sin reproducciones aún" }),
      elR("p", { clase: "text-sm text-gray-500 max-w-sm mx-auto", texto: "Cuando reproduzcas un audio desde el catálogo, aparecerá aquí con su contador de reproducciones." }),
      elR("a", {
        href: "audios.html",
        clase: "inline-block rounded-full bg-accent hover:bg-accent-hover px-6 py-2 text-sm font-semibold text-white transition-all shadow-sm",
        texto: "Ir al catálogo de audios →",
      }),
    ])
  );
}

/** Obtiene el ID de cliente persistente (el mismo que usa el streaming gRPC). */
function obtenerClienteId() {
  const CLAVE_STORAGE = "reproducciones_cliente_id";
  let id = localStorage.getItem(CLAVE_STORAGE);
  if (!id) {
    id = 'cliente-web-' + Math.random().toString(36).substring(2, 10);
    localStorage.setItem(CLAVE_STORAGE, id);
  }
  return id;
}

// ---------------------------------------------------------------- Resolución de títulos desde metadatos

/**
 * Extrae el id numérico del audio a partir del nombre de archivo almacenado.
 * Ej: "cancion1.mp3" → 1, "cancion2.mp3" → 2
 */
function extraerIdDeArchivo(nombreArchivo) {
  const match = nombreArchivo.match(/cancion(\d+)/i);
  return match ? Number(match[1]) : null;
}

/** Caché de títulos: { idAudio -> { titulo, tipo_audio } } */
let cacheTitulos = null;

/**
 * Consulta todos los tipos de audio y sus audios para construir un mapa
 * idAudio → título real de los metadatos.
 * Si falla la consulta de metadatos, devuelve un mapa vacío (se usará el nombre de archivo).
 */
async function cargarMapaTitulos() {
  if (cacheTitulos) return cacheTitulos;

  cacheTitulos = {};
  try {
    if (typeof ApiMetadatos === "undefined" || typeof TIPOS_AUDIO === "undefined") {
      return cacheTitulos;
    }

    const tipos = await ApiMetadatos.listarTipos();

    // Cargar los audios de cada tipo en paralelo
    const promesas = tipos.map(async (tipo) => {
      const definicion = TIPOS_AUDIO[tipo.id];
      if (!definicion) return;

      try {
        const audios = await ApiMetadatos.listarAudios(tipo.id);
        for (const audio of audios) {
          const titulo = definicion.titulo(audio);
          const subtitulo = definicion.subtitulo ? definicion.subtitulo(audio) : "";
          cacheTitulos[audio.id] = {
            titulo: titulo || `Audio #${audio.id}`,
            subtitulo: subtitulo || "",
            tipo_audio: tipo.nombre,
          };
        }
      } catch (e) {
        console.warn(`No se pudieron cargar audios del tipo ${tipo.nombre}:`, e.message);
      }
    });

    await Promise.all(promesas);
  } catch (e) {
    console.warn("No se pudieron cargar los metadatos para resolver títulos:", e.message);
  }

  return cacheTitulos;
}

/**
 * Enriquece los audios del resumen de reproducciones con los títulos reales de metadatos.
 * Reemplaza el nombre de archivo (ej: "cancion1.mp3") por el título del catálogo.
 */
async function resolverTitulos(audios) {
  const mapa = await cargarMapaTitulos();

  for (const audio of audios) {
    const idAudio = extraerIdDeArchivo(audio.titulo);
    if (idAudio !== null && mapa[idAudio]) {
      const meta = mapa[idAudio];
      audio.titulo = meta.titulo;
      // Si el tipo de audio no viene de las estadísticas, usar el de metadatos
      if (!audio.tipo_audio && meta.tipo_audio) {
        audio.tipo_audio = meta.tipo_audio;
      }
    }
  }

  return audios;
}

/** Carga y muestra las reproducciones del cliente. */
async function cargarReproducciones() {
  const contenedor = $r("lista-reproducciones");
  const totalEl = $r("total-reproducciones");
  const clienteId = obtenerClienteId();

  if (!clienteId) {
    contenedor.replaceChildren(
      elR("p", { clase: "texto-tenue col-span-full", texto: "Ingrese un identificador de cliente para consultar sus reproducciones." })
    );
    totalEl.textContent = "";
    return;
  }

  contenedor.replaceChildren(elR("p", { clase: "texto-tenue col-span-full", texto: "Cargando…" }));
  totalEl.textContent = "";

  try {
    const resumen = await ApiEstadisticas.obtenerReproduccionesCliente(clienteId);

    if (!resumen.audios || resumen.audios.length === 0) {
      mostrarVacio();
      totalEl.textContent = "No hay reproducciones registradas.";
      return;
    }

    // Resolver los nombres de archivo a títulos reales de los metadatos
    await resolverTitulos(resumen.audios);

    // Actualizar el total
    const totalAudios = resumen.audios.length;
    const totalReproducciones = resumen.total_reproducciones;
    totalEl.textContent = `${totalAudios} audio${totalAudios !== 1 ? "s" : ""} · ${totalReproducciones} reproducción${totalReproducciones !== 1 ? "es" : ""} en total`;

    // Renderizar tarjetas
    contenedor.replaceChildren(
      ...resumen.audios.map((audio) => crearTarjetaAudio(audio))
    );
  } catch (error) {
    mostrarErrorReproducciones(error.message, error.reintentable !== false);
  }
}

/** Limpia el historial visual (solo en el cliente, no borra datos del servidor). */
function limpiarHistorial() {
  const contenedor = $r("lista-reproducciones");
  const totalEl = $r("total-reproducciones");
  contenedor.replaceChildren(
    elR("div", { clase: "col-span-full text-center py-16 space-y-4" }, [
      elR("span", { clase: "text-6xl block", "aria-hidden": "true", texto: "🗑️" }),
      elR("h3", { clase: "text-lg font-semibold text-gray-700", texto: "Historial limpiado" }),
      elR("p", { clase: "text-sm text-gray-500", texto: "Presiona \"Recargar\" para volver a cargar las reproducciones del servidor." }),
    ])
  );
  totalEl.textContent = "";
}

// ---------------------------------------------------------------- Inicio

document.addEventListener("DOMContentLoaded", () => {
  $r("btn-recargar").addEventListener("click", cargarReproducciones);
  $r("btn-limpiar").addEventListener("click", limpiarHistorial);

  cargarReproducciones();
});
