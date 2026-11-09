// Lógica de la interfaz del administrador: registro de audios (metadatos + archivo) según el tipo.
//
// Flujo (el del administrador de consola, pero con los metadatos propios de cada tipo):
//   tipos de audio  ->  formulario del tipo (metadatos + archivo .mp3)  ->  confirmación
//
// La vista actual vive en el hash de la URL, así el botón "atrás" del navegador
// y recargar la página funcionan:
//   #/                      tipos de audio
//   #/tipos/1/nuevo         formulario para registrar un audio del tipo 1

const $ = (id) => document.getElementById(id);

const estado = {
  tipos: null, // caché de los tipos de audio (se piden una sola vez)
  navegacion: 0, // contador para descartar respuestas de navegaciones ya superadas
  envio: null, // { cancelar() } mientras se registra un audio; null si no hay ninguno en curso
  liberar: null, // libera los recursos de la vista actual (URL de la vista previa del audio)
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
const hrefFormulario = (idTipo) => `#/tipos/${idTipo}/nuevo`;

const CLASES_ENTRADA =
  "block w-full rounded-xl bg-[#f5f5f7] border border-gray-200 px-4 py-2.5 text-gray-800 placeholder:text-gray-400 " +
  "focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent transition disabled:opacity-60 disabled:cursor-not-allowed";
const CLASES_BOTON_PRIMARIO =
  "rounded-full bg-accent hover:bg-accent-hover active:brightness-90 px-6 py-2.5 font-semibold text-white transition-all shadow-sm " +
  "disabled:opacity-60 disabled:cursor-not-allowed";
const CLASES_BOTON_SECUNDARIO =
  "inline-flex items-center rounded-full border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors";
const CLASES_TARJETA = "bg-white rounded-2xl border border-gray-200 shadow-sm p-6 space-y-4";
const CLASES_TITULO_TARJETA = "text-sm font-bold uppercase tracking-wide text-gray-500";

function textoTamanio(bytes) {
  return bytes < 1024 * 1024 ? `${Math.max(1, Math.round(bytes / 1024))} KB` : `${(bytes / 1024 / 1024).toFixed(1)} MB`;
}

function textoReloj(segundos) {
  const minutos = Math.floor(segundos / 60);
  return `${minutos}:${String(segundos % 60).padStart(2, "0")} min`;
}

/** Muestra "—" cuando el servidor no trae el dato. Un número 0 también es "vacío": el servidor no distingue ausente de 0. */
function textoValor(valor, campo) {
  const numerico = campo.tipo === "entero" || campo.tipo === "isbn";
  if (valor === null || valor === undefined || valor === "" || (numerico && valor === 0)) return "—";
  return String(campo.formato ? campo.formato(valor) : valor);
}

// ---------------------------------------------------------------- Estado del servidor (chip del encabezado)

const CHIPS_SERVIDOR = {
  consultando: ["Consultando servidor…", "bg-yellow-100 text-yellow-700 border-yellow-200"],
  ok: ["Servidor conectado", "bg-green-100 text-green-700 border-green-200"],
  error: ["Servidor no disponible", "bg-red-100 text-red-700 border-red-200"],
};

function estadoServidor(clave) {
  const [texto, clases] = CHIPS_SERVIDOR[clave];
  const chip = $("estado-servidor");
  chip.textContent = texto;
  chip.className = `whitespace-nowrap text-xs font-medium px-3 py-1 rounded-full border ${clases}`;
}

// ---------------------------------------------------------------- Rutas

function rutaActual() {
  const partes = location.hash.replace(/^#\/?/, "").split("/").filter(Boolean);
  if (partes.length === 3 && partes[0] === "tipos" && /^\d+$/.test(partes[1]) && partes[2] === "nuevo") {
    return { vista: "formulario", idTipo: Number(partes[1]) };
  }
  return { vista: "tipos" };
}

async function obtenerTipos() {
  if (!estado.tipos) estado.tipos = await ApiAdmin.listarTipos();
  return estado.tipos;
}

async function enrutar() {
  const numero = ++estado.navegacion;
  const vigente = () => numero === estado.navegacion;
  const ruta = rutaActual();

  estado.envio?.cancelar(); // al cambiar de vista no debe seguir una subida cuyo resultado ya no se mostraría
  if (!estado.tipos) estadoServidor("consultando");
  pintarCargando();
  try {
    const tipos = await obtenerTipos();
    if (!vigente()) return;
    estadoServidor("ok");

    if (ruta.vista === "tipos") {
      pintarTipos(tipos);
      return;
    }

    const volver = { href: hrefTipos(), texto: "Volver a los tipos" };
    const tipo = tipos.find((t) => t.id === ruta.idTipo);
    if (!tipo) {
      pintarError("Tipo de audio no encontrado.", { volver });
      return;
    }
    const definicion = TIPOS_AUDIO[tipo.id];
    if (!definicion) {
      pintarError(`Tipo de audio #${tipo.id} no soportado.`, { volver });
      return;
    }
    pintarFormulario(tipo, definicion);
  } catch (error) {
    if (!vigente()) return;
    estadoServidor("error");
    pintarError(error.message, { reintentar: error.reintentable !== false });
  }
}

// ---------------------------------------------------------------- Vistas

/** Reemplaza el contenido de la página y lleva el foco al título (útil para lectores de pantalla). */
function mostrar(tituloPagina, nodos) {
  estado.liberar?.();
  estado.liberar = null;
  document.title = `${tituloPagina} · Administración de Audios`;
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
        volver ? el("a", { href: volver.href, clase: CLASES_BOTON_SECUNDARIO, texto: `← ${volver.texto}` }) : null,
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
          el("a", { href: hrefFormulario(tipo.id), "data-tipo": tipo.id, clase: "tarjeta bg-white rounded-2xl border border-gray-200 shadow-sm p-6 flex flex-col gap-3 hover:shadow-md transition-shadow" }, [
            el("span", { clase: "text-3xl", "aria-hidden": "true", texto: TIPOS_AUDIO[tipo.id]?.icono || "🎧" }),
            el("h3", { clase: "text-lg font-bold text-gray-900", texto: tipo.nombre }),
            el("span", { clase: "text-sm font-semibold text-accent", texto: "Registrar audio →" }),
          ])
        )
      );

  mostrar("Tipos de audio", [
    encabezado({ titulo: "Registrar un audio", subtitulo: "Seleccione el tipo de audio. Cada tipo tiene sus propios metadatos." }),
    tarjetas,
  ]);
}

// ---------------------------------------------------------------- Formulario: campos de metadatos

/**
 * Un campo del formulario de metadatos. `leer()` devuelve { valor } (ya convertido al tipo del DTO),
 * { vacio: true } si es opcional y no se escribió nada, o { error } con el mensaje de validación.
 */
function crearCampo(campo) {
  const clase = TIPOS_CAMPO[campo.tipo];
  const id = `campo-${campo.clave}`;
  const idAyuda = `${id}-ayuda`;
  const idError = `${id}-error`;

  const atributos = {
    id,
    name: campo.clave,
    autocomplete: "off",
    clase: CLASES_ENTRADA,
    "aria-describedby": campo.ayuda ? `${idAyuda} ${idError}` : idError,
  };
  if (campo.obligatorio) atributos.required = "";
  if (campo.placeholder) atributos.placeholder = campo.placeholder;

  const entrada =
    campo.tipo === "area"
      ? el("textarea", { ...atributos, rows: "3" })
      : el("input", { ...atributos, type: "text", inputmode: clase.modo, ...(campo.sugerencias ? { list: `${id}-lista` } : {}) });
  const lista = campo.sugerencias
    ? el("datalist", { id: `${id}-lista` }, campo.sugerencias.map((valor) => el("option", { value: valor })))
    : null;
  const mensajeError = el("p", { id: idError, role: "alert", clase: "hidden text-xs text-red-600" });

  const control = {
    campo,
    entrada,
    editadoAMano: false, // si la persona escribió en el campo, no se sobrescribe con datos del archivo
    nodo: el("div", { clase: `space-y-1 ${campo.ancho === "completo" ? "sm:col-span-2" : ""}` }, [
      el("label", { for: id, clase: "text-sm font-medium text-gray-600", texto: campo.etiqueta }, [
        campo.obligatorio ? el("span", { clase: "text-accent ml-0.5", "aria-hidden": "true", texto: "*" }) : null,
      ]),
      entrada,
      lista,
      campo.ayuda ? el("p", { id: idAyuda, clase: "text-xs text-gray-400", texto: campo.ayuda }) : null,
      mensajeError,
    ]),

    mostrarError(mensaje) {
      entrada.setAttribute("aria-invalid", "true");
      mensajeError.textContent = mensaje;
      mensajeError.classList.remove("hidden");
    },
    limpiarError() {
      entrada.removeAttribute("aria-invalid");
      mensajeError.classList.add("hidden");
    },
    leer() {
      const texto = entrada.value.trim();
      if (texto === "") return campo.obligatorio ? { error: "Este campo es obligatorio." } : { vacio: true };
      const error = clase.validar(texto, campo);
      return error ? { error } : { valor: clase.convertir(texto) };
    },
  };

  entrada.addEventListener("input", () => {
    control.editadoAMano = true;
    control.limpiarError();
  });
  return control;
}

// ---------------------------------------------------------------- Formulario: archivo de audio

/**
 * Zona para elegir (o arrastrar) el archivo de audio, con vista previa para escucharlo antes de registrarlo.
 * `alConocerDuracion(segundos)` se llama cuando el navegador lee la duración del archivo elegido.
 */
function crearZonaArchivo({ alConocerDuracion }) {
  const formato = CONFIG_ADMIN.FORMATO_AUDIO;
  let archivo = null;
  let urlVistaPrevia = null;

  const entrada = el("input", {
    id: "campo-archivo",
    name: "archivo",
    type: "file",
    accept: `.${formato},audio/mpeg`,
    clase: "sr-only",
    "aria-describedby": "archivo-ayuda archivo-aviso archivo-error",
  });
  const titulo = el("span", { clase: "text-sm font-semibold text-gray-700 break-words max-w-full", texto: `Seleccione o arrastre un archivo .${formato}` });
  const detalle = el("span", { clase: "text-xs text-gray-400" });
  const zona = el("label", { for: "campo-archivo", clase: "zona-archivo relative flex flex-col items-center justify-center gap-1 rounded-2xl border-2 border-dashed border-gray-300 bg-[#f5f5f7] px-6 py-8 text-center cursor-pointer transition-colors" }, [
    el("span", { clase: "text-3xl", "aria-hidden": "true", texto: "🎧" }),
    titulo,
    detalle,
    entrada,
  ]);
  const vistaPrevia = el("audio", { controls: "", preload: "metadata", clase: "hidden w-full", "aria-label": "Vista previa del audio elegido" });
  const aviso = el("p", { id: "archivo-aviso", clase: "hidden text-xs text-amber-600" });
  const mensajeError = el("p", { id: "archivo-error", role: "alert", clase: "hidden text-xs text-red-600" });

  function mostrarError(mensaje) {
    mensajeError.textContent = mensaje;
    mensajeError.classList.remove("hidden");
    zona.classList.add("invalida");
  }

  function limpiarError() {
    mensajeError.classList.add("hidden");
    zona.classList.remove("invalida");
  }

  /** null si el archivo es válido; si no, el mensaje de error. */
  function validar() {
    if (!archivo) return "Seleccione el archivo de audio.";
    if (!archivo.name.toLowerCase().endsWith(`.${formato}`)) {
      return `El archivo debe ser .${formato}: es el único formato que entrega el servidor de streaming.`;
    }
    if (archivo.size === 0) return "El archivo está vacío.";
    return null;
  }

  function liberar() {
    if (urlVistaPrevia) URL.revokeObjectURL(urlVistaPrevia);
    urlVistaPrevia = null;
    vistaPrevia.removeAttribute("src");
    vistaPrevia.load();
  }

  function elegir(nuevo) {
    liberar();
    limpiarError();
    aviso.classList.add("hidden");
    archivo = nuevo || null;

    if (!archivo) {
      titulo.textContent = `Seleccione o arrastre un archivo .${formato}`;
      detalle.textContent = "";
      vistaPrevia.classList.add("hidden");
      return;
    }
    titulo.textContent = archivo.name;
    detalle.textContent = `${textoTamanio(archivo.size)} · haga clic o arrastre otro archivo para cambiarlo`;
    urlVistaPrevia = URL.createObjectURL(archivo);
    vistaPrevia.src = urlVistaPrevia;
    vistaPrevia.classList.remove("hidden");

    const mensaje = validar();
    if (mensaje) mostrarError(mensaje);
  }

  vistaPrevia.addEventListener("loadedmetadata", () => {
    if (!archivo || !Number.isFinite(vistaPrevia.duration)) return;
    const segundos = Math.round(vistaPrevia.duration);
    detalle.textContent = `${textoTamanio(archivo.size)} · ${textoReloj(segundos)} · haga clic o arrastre otro archivo para cambiarlo`;
    alConocerDuracion(segundos);
  });
  vistaPrevia.addEventListener("error", () => {
    if (!vistaPrevia.hasAttribute("src") || validar()) return; // si ya hay un error más claro (vacío, extensión), no se duplica
    aviso.textContent = "El navegador no pudo leer este archivo como audio. Verifique que sea un MP3 válido.";
    aviso.classList.remove("hidden");
  });

  entrada.addEventListener("change", () => elegir(entrada.files[0]));
  zona.addEventListener("dragover", (evento) => {
    evento.preventDefault();
    if (!entrada.disabled) zona.classList.add("arrastrando");
  });
  zona.addEventListener("dragleave", () => zona.classList.remove("arrastrando"));
  zona.addEventListener("drop", (evento) => {
    evento.preventDefault();
    zona.classList.remove("arrastrando");
    if (entrada.disabled) return;
    const soltado = evento.dataTransfer.files[0];
    if (!soltado) return;
    try {
      entrada.files = evento.dataTransfer.files; // mantiene el <input> sincronizado con lo que se ve
    } catch {
      /* navegadores que no permiten asignarlo: el archivo igual queda elegido en `archivo` */
    }
    elegir(soltado);
  });

  return {
    nodo: el("div", { clase: "space-y-2" }, [
      zona,
      el("p", {
        id: "archivo-ayuda",
        clase: "text-xs text-gray-400",
        texto: `Se guardará en el servidor como ${CONFIG_ADMIN.archivoDeAudio(0, "{ID}")}, el nombre con el que el cliente lo reproduce. El ID se asigna al registrar.`,
      }),
      vistaPrevia,
      aviso,
      mensajeError,
    ]),
    entrada,
    obtener: () => archivo,
    validar,
    mostrarError,
    liberar,
  };
}

// ---------------------------------------------------------------- Formulario de registro

function pintarFormulario(tipo, definicion) {
  const controles = definicion.campos.map(crearCampo);
  const zonaArchivo = crearZonaArchivo({
    // Para campos como la duración del ruido blanco: se propone el dato del archivo mientras la persona no lo haya escrito.
    alConocerDuracion(segundos) {
      const control = controles.find((c) => c.campo.autocompletarDesdeArchivo === "duracion");
      if (!control || control.editadoAMano) return;
      control.entrada.value = String(segundos);
      control.limpiarError();
    },
  });

  const estadoTexto = el("p", { clase: "text-sm text-gray-600" });
  const barra = el("div", { clase: "h-full bg-accent transition-all", style: "width: 0%" });
  const pista = el("div", { clase: "hidden h-2 rounded-full bg-gray-200 overflow-hidden", "aria-hidden": "true" }, [barra]);
  const panelEnvio = el("div", { role: "status", clase: "hidden space-y-2" }, [estadoTexto, pista]);
  const alerta = el("div", { role: "alert", clase: "hidden rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-600" });
  const botonRegistrar = el("button", { type: "submit", clase: CLASES_BOTON_PRIMARIO, texto: "Registrar audio" });
  const campos = el("fieldset", { clase: "space-y-6 min-w-0" }, [
    el("section", { clase: CLASES_TARJETA }, [
      el("h3", { clase: CLASES_TITULO_TARJETA, texto: `Metadatos · ${tipo.nombre}` }),
      el("div", { clase: "grid gap-4 sm:grid-cols-2" }, controles.map((c) => c.nodo)),
      el("p", { clase: "text-xs text-gray-400" }, [
        el("span", { clase: "text-accent", "aria-hidden": "true", texto: "*" }),
        " Obligatorio. Los campos numéricos que se dejen vacíos se guardan como 0.",
      ]),
    ]),
    el("section", { clase: CLASES_TARJETA }, [
      el("h3", { clase: CLASES_TITULO_TARJETA }, ["Archivo de audio ", el("span", { clase: "text-accent", "aria-hidden": "true", texto: "*" })]),
      zonaArchivo.nodo,
    ]),
  ]);

  function bloquear(bloqueado) {
    campos.disabled = bloqueado;
    botonRegistrar.disabled = bloqueado;
  }

  function mostrarEstado(texto) {
    estadoTexto.textContent = texto;
    panelEnvio.classList.remove("hidden");
  }

  function mostrarProgreso(fraccion) {
    const porcentaje = Math.round(fraccion * 100);
    pista.classList.remove("hidden");
    barra.style.width = `${porcentaje}%`;
    estadoTexto.textContent = `Subiendo el audio… ${porcentaje} %`;
  }

  /** Valida todo el formulario. Devuelve los metadatos listos para el servidor, o null (y lleva el foco al primer error). */
  function validarFormulario() {
    const valores = {};
    let primerError = null;

    for (const control of controles) {
      const resultado = control.leer();
      if (resultado.error) {
        control.mostrarError(resultado.error);
        primerError ??= control.entrada;
      } else {
        control.limpiarError();
        if (!resultado.vacio) valores[control.campo.clave] = resultado.valor;
      }
    }

    const errorArchivo = zonaArchivo.validar();
    if (errorArchivo) {
      zonaArchivo.mostrarError(errorArchivo);
      primerError ??= zonaArchivo.entrada;
    }

    if (primerError) {
      primerError.focus();
      return null;
    }
    return valores;
  }

  async function enviar(evento) {
    evento.preventDefault();
    alerta.classList.add("hidden");
    const valores = validarFormulario();
    if (!valores) return;

    const controlador = new AbortController();
    const envio = { cancelar: () => controlador.abort() };
    estado.envio = envio;
    bloquear(true);
    pista.classList.add("hidden");
    barra.style.width = "0%";

    try {
      const resultado = await ApiAdmin.registrarAudio(tipo.id, valores, zonaArchivo.obtener(), {
        senal: controlador.signal,
        alEstado: mostrarEstado,
        alProgresar: mostrarProgreso,
      });
      if (controlador.signal.aborted) return;
      estadoServidor("ok");
      pintarConfirmacion(tipo, definicion, resultado);
    } catch (error) {
      if (controlador.signal.aborted) return; // la persona salió de la vista: ya no hay nada que mostrar
      if (error.estado === 0 && error.reintentable) estadoServidor("error"); // sin respuesta del servidor
      panelEnvio.classList.add("hidden");
      alerta.textContent = error.message;
      alerta.classList.remove("hidden");
      bloquear(false);
      alerta.scrollIntoView({ block: "nearest" });
    } finally {
      if (estado.envio === envio) estado.envio = null;
    }
  }

  mostrar(`Registrar ${tipo.nombre}`, [
    el("div", { clase: "max-w-3xl space-y-6" }, [
      encabezado({
        titulo: `${definicion.icono} Registrar ${tipo.nombre}`,
        subtitulo: `Complete los metadatos del audio y adjunte el archivo .${CONFIG_ADMIN.FORMATO_AUDIO}.`,
        volver: { href: hrefTipos(), texto: "Volver a los tipos" },
      }),
      el("form", { novalidate: "", autocomplete: "off", clase: "space-y-6", onsubmit: enviar }, [
        campos,
        panelEnvio,
        alerta,
        el("div", { clase: "flex flex-wrap items-center gap-3" }, [
          botonRegistrar,
          el("a", { href: hrefTipos(), clase: CLASES_BOTON_SECUNDARIO, texto: "Cancelar" }),
        ]),
      ]),
    ]),
  ]);
  estado.liberar = zonaArchivo.liberar;
}

// ---------------------------------------------------------------- Confirmación

function pintarConfirmacion(tipo, definicion, { audio, nombreArchivo }) {
  const etiqueta = (texto) => el("span", { clase: "rounded-full bg-accent-light text-accent text-xs font-semibold px-3 py-1", texto });

  const filas = [
    ["ID asignado", String(audio.id)],
    ["Archivo de audio", nombreArchivo],
    ...definicion.campos.map((campo) => [campo.etiqueta, textoValor(audio[campo.clave], campo)]),
  ].flatMap(([nombre, valor]) => [el("dt", { texto: nombre }), el("dd", { texto: valor })]);

  mostrar("Audio registrado", [
    el("div", { clase: "max-w-3xl space-y-6" }, [
      encabezado({
        titulo: "✅ Audio registrado",
        subtitulo: "El archivo y sus metadatos quedaron guardados en el servidor.",
      }),
      el("section", { clase: CLASES_TARJETA }, [
        el("div", { clase: "flex flex-wrap items-center gap-2" }, [
          el("h3", { clase: `${CLASES_TITULO_TARJETA} mr-2`, texto: "Audio registrado" }),
          etiqueta(`ID ${audio.id}`),
          etiqueta(audio.tipo?.nombre || tipo.nombre),
        ]),
        el("dl", { clase: "ficha" }, filas),
      ]),
      el("div", { clase: "flex flex-wrap items-center gap-3" }, [
        el("button", {
          type: "button",
          clase: CLASES_BOTON_PRIMARIO,
          texto: `Registrar otro audio de ${tipo.nombre}`,
          onclick: () => pintarFormulario(tipo, definicion),
        }),
        el("a", { href: hrefTipos(), clase: CLASES_BOTON_SECUNDARIO, texto: "Elegir otro tipo" }),
      ]),
    ]),
  ]);
}

// ---------------------------------------------------------------- Inicio

// No perder una subida a medias por cerrar la pestaña o recargar sin querer.
window.addEventListener("beforeunload", (evento) => {
  if (!estado.envio) return;
  evento.preventDefault();
  evento.returnValue = "";
});

document.addEventListener("DOMContentLoaded", () => {
  window.addEventListener("hashchange", enrutar);
  enrutar();
});
