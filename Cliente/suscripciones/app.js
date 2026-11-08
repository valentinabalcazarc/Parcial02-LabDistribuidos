// Lógica de la interfaz del cliente.

const estado = {
  suscripciones: [],
  ocupado: false,
  inicioPago: 0,
};

const $ = (id) => document.getElementById(id);

const formatoMoneda = new Intl.NumberFormat("es-CO", { style: "currency", currency: "COP", maximumFractionDigits: 0 });
const formatearMonto = (valor) => formatoMoneda.format(valor);
const formatearFecha = (iso) => (iso ? new Date(iso).toLocaleString("es-CO") : "—");
const nombreSuscripcion = (id) => estado.suscripciones.find((s) => s.id === id)?.nombre || id;

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

// ---------------------------------------------------------------- Estado del servidor

async function cargarEstadoServidor() {
  const indicador = $("estado-servidor");
  try {
    await Api.estado();
    indicador.className = "text-xs font-medium px-3 py-1 rounded-full bg-green-50 text-green-700 border border-green-200";
    indicador.textContent = "Servidor activo";
  } catch {
    indicador.className = "text-xs font-medium px-3 py-1 rounded-full bg-red-50 text-red-600 border border-red-200";
    indicador.textContent = "Servidor no disponible";
  }
}

// ---------------------------------------------------------------- Suscripciones

async function cargarSuscripciones() {
  const contenedor = $("lista-suscripciones");
  try {
    estado.suscripciones = await Api.listarSuscripciones();
  } catch (error) {
    contenedor.replaceChildren(el("p", { clase: "texto-tenue col-span-full", texto: `No fue posible listar las suscripciones: ${error.message}` }));
    return;
  }

  contenedor.replaceChildren(
    ...estado.suscripciones.map((s) => {
      // Build feature items from the description
      const features = [];
      features.push(`${s.descripcion}`);
      features.push(`${s.cuentas} cuenta${s.cuentas > 1 ? "s" : ""} simultánea${s.cuentas > 1 ? "s" : ""}`);
      features.push("Acceso a todo el catálogo de música");
      features.push("Sin anuncios");

      const featureList = el("ul", { clase: "space-y-3 mt-6 mb-8 flex-1" },
        features.map(f => {
          const li = document.createElement("li");
          li.className = "flex items-start gap-2 text-sm text-gray-600 leading-relaxed";
          const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
          svg.setAttribute("class", "w-4 h-4 mt-0.5 shrink-0 text-accent");
          svg.setAttribute("viewBox", "0 0 20 20");
          svg.setAttribute("fill", "currentColor");
          const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
          path.setAttribute("fill-rule", "evenodd");
          path.setAttribute("d", "M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z");
          path.setAttribute("clip-rule", "evenodd");
          svg.appendChild(path);
          li.appendChild(svg);
          li.appendChild(document.createTextNode(f));
          return li;
        })
      );

      return el("article", { clase: "tarjeta bg-white rounded-2xl border border-gray-200 shadow-sm p-6 flex flex-col hover:shadow-md transition-shadow", "data-id": s.id }, [
        el("h3", { clase: "text-xl font-bold text-accent", texto: s.nombre }),
        el("div", { clase: "mt-2" }, [
          el("span", { clase: "text-2xl font-bold text-gray-900", texto: `${formatearMonto(s.precio)}` }),
          el("span", { clase: "text-sm text-gray-500", texto: " al mes." }),
        ]),
        el("p", { clase: "text-sm font-semibold text-gray-900 mt-1", texto: "Primer mes gratis para nuevas suscripciones." }),
        featureList,
        el("button", {
          clase: "rounded-full bg-accent hover:bg-accent-hover active:brightness-90 px-6 py-2.5 text-sm font-semibold text-white transition-all shadow-sm w-fit",
          type: "button",
          texto: "Elegir plan",
          onclick: () => seleccionarSuscripcion(s.id, true),
        }),
      ]);
    })
  );

  const selector = $("suscripcion");
  const actual = selector.value;
  selector.replaceChildren(...estado.suscripciones.map((s) => el("option", { value: s.id, texto: s.nombre })));
  seleccionarSuscripcion(actual || estado.suscripciones[0]?.id, false);
}

function seleccionarSuscripcion(id, desplazar) {
  const suscripcion = estado.suscripciones.find((s) => s.id === id);
  if (!suscripcion) return;
  $("suscripcion").value = id;
  $("monto").value = formatearMonto(suscripcion.precio);
  document.querySelectorAll(".tarjeta").forEach((t) => {
    t.classList.toggle("seleccionada", t.dataset.id === id);
    // Also scale up selected card slightly for emphasis
    t.classList.toggle("scale-[1.02]", t.dataset.id === id);
  });
  if (desplazar) $("form-pago").scrollIntoView({ behavior: "smooth", block: "center" });
}

// ---------------------------------------------------------------- Bitácora

function limpiarBitacora() {
  $("bitacora").replaceChildren();
  $("resultado").className = "resultado oculto";
}

function registrar({ tipo, mensaje }) {
  const segundos = ((Date.now() - estado.inicioPago) / 1000).toFixed(1);
  const hora = new Date().toLocaleTimeString("es-CO");
  const item = el("li", { clase: tipo }, [el("span", { clase: "hora", texto: `${hora} (+${segundos} s)` }), mensaje]);
  $("bitacora").append(item);
  item.scrollIntoView({ block: "nearest" });
}

function mostrarResultado(tipo, titulo, filas = []) {
  const caja = $("resultado");
  caja.className = `resultado ${tipo}`;
  const lista = el("dl", {}, filas.flatMap(([etiqueta, valor]) => [el("dt", { texto: etiqueta }), el("dd", { texto: valor })]));
  caja.replaceChildren(el("h3", { texto: titulo }), lista);
}

function setOcupado(ocupado) {
  estado.ocupado = ocupado;
  $("btn-pagar").disabled = ocupado;
  document.querySelectorAll("[data-reenviar]").forEach((b) => (b.disabled = ocupado));
}

// ---------------------------------------------------------------- Pago

async function procesarPago(evento) {
  evento.preventDefault();
  if (estado.ocupado) return;

  const suscripcion = estado.suscripciones.find((s) => s.id === $("suscripcion").value);
  if (!suscripcion) return;

  const datosPago = {
    cedula: $("cedula").value.trim(),
    suscripcion_id: suscripcion.id,
    monto: suscripcion.precio,
    nombre_titular: $("nombre").value.trim(),
    numero_tarjeta: $("tarjeta").value.replace(/\s+/g, ""),
  };

  limpiarBitacora();
  estado.inicioPago = Date.now();
  setOcupado(true);
  try {
    registrar({ tipo: "info", mensaje: `Solicitando token de idempotencia (POST /tokens) para ${suscripcion.nombre}…` });
    const token = await Api.generarToken(datosPago);
    registrar({ tipo: "ok", mensaje: `Token de idempotencia obtenido: ${token.token_idempotencia}` });
    await ejecutarPago(token.token_idempotencia, token.monto);
  } catch (error) {
    registrar({ tipo: "error", mensaje: error.message });
    mostrarResultado("error", "No fue posible obtener el token de idempotencia", [["Detalle", error.message]]);
  } finally {
    setOcupado(false);
  }
}

async function reenviarConToken(token, monto) {
  if (estado.ocupado) return;
  limpiarBitacora();
  estado.inicioPago = Date.now();
  setOcupado(true);
  registrar({ tipo: "info", mensaje: `Reenviando pago con un token ya existente: ${token}` });
  try {
    await ejecutarPago(token, monto);
  } finally {
    setOcupado(false);
  }
}

async function ejecutarPago(token, monto) {
  try {
    const { datos, intento } = await Api.realizarPagoConReintentos(token, monto, registrar);
    const enQueIntento = intento === 0 ? "Primer intento" : `Reintento ${intento}`;

    if (datos.transaccion_realizada) {
      mostrarResultado("ok", "✔ Pago realizado exitosamente", [
        ["Obtenido en", enQueIntento],
        ["Monto de la transacción", formatearMonto(datos.monto_transaccion)],
        ["Transacción realizada", "Sí"],
        ["Realizada previamente", "No"],
        ["ID transacción", datos.id_transaccion],
        ["Token", datos.token_idempotencia],
        ["Fecha", formatearFecha(datos.fecha_transaccion)],
      ]);
    } else if (datos.previamente_realizada) {
      mostrarResultado("aviso", "⚠ Este token ya había sido utilizado: el pago no se ejecutó de nuevo", [
        ["Obtenido en", enQueIntento],
        ["Monto enviado", formatearMonto(datos.monto_transaccion)],
        ["Transacción realizada", "No"],
        ["Realizada previamente", "Sí"],
        ["Monto de la transacción previa", formatearMonto(datos.monto_transaccion_previa)],
        ["ID transacción previa", datos.id_transaccion],
        ["Fecha transacción previa", formatearFecha(datos.fecha_transaccion)],
        ["Token", datos.token_idempotencia],
      ]);
    }
    registrar({ tipo: "ok", mensaje: datos.mensaje });
  } catch (error) {
    if (error instanceof ErrorSinRespuesta) {
      registrar({ tipo: "error", mensaje: error.message });
      mostrarResultado("error", "✖ No fue posible acceder al servidor para realizar el pago", [
        ["Token utilizado", token],
        ["Intentos", `1 intento inicial + ${CONFIG.ESPERAS_REINTENTO_MS.length} reintentos sin respuesta`],
        ["Sugerencia", "Puede reenviar el pago con el mismo token: si el servidor alcanzó a procesarlo, no se cobrará dos veces."],
      ]);
    } else {
      registrar({ tipo: "error", mensaje: error.message });
      mostrarResultado("error", "✖ El servidor rechazó el pago", [["Detalle", error.message]]);
    }
  }
}

// ---------------------------------------------------------------- Consulta de transacciones

async function consultarTransacciones(evento) {
  evento?.preventDefault();
  const contenedor = $("resultado-consulta");
  const cedula = $("cedula-consulta").value.trim();
  contenedor.replaceChildren(el("p", { clase: "texto-tenue", texto: "Consultando…" }));

  let datos;
  try {
    datos = await Api.consultarTransacciones(cedula);
  } catch (error) {
    contenedor.replaceChildren(el("p", { clase: "texto-tenue", texto: `Error: ${error.message}` }));
    return;
  }

  const filasTransacciones = datos.transacciones.map((t, i) =>
    el("tr", {}, [
      el("td", { texto: String(i + 1) }),
      el("td", { texto: formatearFecha(t.fecha) }),
      el("td", { texto: nombreSuscripcion(t.suscripcion_id) }),
      el("td", { texto: formatearMonto(t.monto) }),
      el("td", {}, el("code", { texto: t.token_idempotencia })),
      el("td", {},
        el("button", {
          clase: "rounded-full bg-[#FA2D48] hover:bg-[#E0223C] px-3 py-1 text-xs font-semibold text-white transition-all", type: "button", "data-reenviar": "", texto: "Reenviar con este token",
          onclick: () => reenviarConToken(t.token_idempotencia, t.monto),
        })
      ),
    ])
  );

  contenedor.replaceChildren(
    el("p", { clase: "resumen" }, [`Cédula ${datos.cedula} · transacciones realizadas: `, el("span", { texto: String(datos.total_transacciones) })]),
    datos.transacciones.length === 0
      ? el("p", { clase: "texto-tenue", texto: "No hay transacciones registradas." })
      : el("div", { clase: "tabla-contenedor" },
          el("table", {}, [
            el("thead", {}, el("tr", {}, ["#", "Fecha", "Suscripción", "Monto", "Token de idempotencia", ""].map((t) => el("th", { texto: t })))),
            el("tbody", {}, filasTransacciones),
          ])
        )
  );
}

// ---------------------------------------------------------------- Inicio

document.addEventListener("DOMContentLoaded", () => {
  $("form-pago").addEventListener("submit", procesarPago);
  $("form-consulta").addEventListener("submit", consultarTransacciones);
  $("suscripcion").addEventListener("change", (e) => seleccionarSuscripcion(e.target.value, false));
  $("btn-recargar-suscripciones").addEventListener("click", () => {
    cargarEstadoServidor();
    cargarSuscripciones();
  });

  cargarEstadoServidor();
  cargarSuscripciones();
});
