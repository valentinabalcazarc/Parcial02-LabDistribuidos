// Módulo de reacciones en tiempo real mediante WebSocket (STOMP + SockJS).
//
// Flujo:
//   1. Cuando el usuario reproduce un audio, se conecta al broker de reacciones
//      y se suscribe al topic /brokerDeReacciones/audio/{idAudio}.
//      (Esa suscripción le dice al servidor que este cliente está escuchando ese audio.)
//   2. El usuario puede enviar reacciones (👍, ❤️, 😂, 🔥, 👏) al destino /app/enviarReaccion.
//   3. El servidor reenvía la reacción a todos los clientes suscritos al mismo topic (callback).
//   4. Cada reacción recibida se muestra como una animación flotante sobre el reproductor.

const Reacciones = (() => {
  // ─── Configuración ───────────────────────────────────────────
  const GATEWAY_URL = "http://localhost:8080";
  const WS_PATH = "/api/reacciones/ws";

  // Los emojis que el usuario puede enviar
  const EMOJIS = ["👍", "❤️", "😂", "🔥", "👏"];

  // ─── Estado ──────────────────────────────────────────────────
  let stompClient = null;
  let suscripcionActual = null;
  let idAudioActual = null;
  let conectado = false;

  // ─── Conexión STOMP sobre SockJS ─────────────────────────────

  /**
   * Conecta al broker de reacciones y, una vez conectado, se suscribe al topic
   * del audio indicado. Devuelve una Promise que se resuelve cuando la conexión
   * y la suscripción están listas.
   *
   * @param {string} idAudio  Identificador del audio que el cliente está reproduciendo.
   * @param {(reaccion: {idAudio: string, tipoReaccion: string}) => void} onReaccion
   *        Callback que se invoca cada vez que llega una reacción por el WebSocket.
   */
  function conectar(idAudio, onReaccion) {
    // Si ya estamos conectados al mismo audio, no reconectar
    if (conectado && idAudioActual === idAudio) return Promise.resolve();

    // Si hay conexión previa (otro audio), desconectar primero
    desconectar();

    idAudioActual = idAudio;

    return new Promise((resolve, reject) => {
      const socket = new SockJS(GATEWAY_URL + WS_PATH);
      stompClient = Stomp.over(socket);

      // Silenciar los logs de debug de STOMP en la consola (se pueden habilitar para depuración)
      stompClient.debug = null;

      stompClient.connect(
        {}, // headers (vacío; el gateway agrega X-Gateway-Token)
        function onConectado() {
          conectado = true;
          console.log(`[Reacciones] Conectado al broker. Suscribiendo al audio "${idAudio}"…`);

          // Suscripción al topic: el servidor registra que esta sesión "reproduce" este audio
          const topic = "/brokerDeReacciones/audio/" + idAudio;
          suscripcionActual = stompClient.subscribe(topic, function (mensaje) {
            try {
              const reaccion = JSON.parse(mensaje.body);
              if (onReaccion) onReaccion(reaccion);
            } catch (e) {
              console.warn("[Reacciones] Mensaje no válido:", mensaje.body);
            }
          });

          console.log(`[Reacciones] Suscrito a ${topic}`);
          resolve();
        },
        function onError(error) {
          conectado = false;
          console.error("[Reacciones] Error de conexión WebSocket:", error);
          reject(error);
        }
      );
    });
  }

  /**
   * Envía una reacción al servidor. Solo funciona si hay conexión activa.
   *
   * @param {string} tipoReaccion  Emoji o texto de la reacción (p. ej. "👍").
   */
  function enviar(tipoReaccion) {
    if (!stompClient || !conectado || !idAudioActual) {
      console.warn("[Reacciones] No conectado; no se puede enviar la reacción.");
      return;
    }

    const payload = JSON.stringify({
      idAudio: idAudioActual,
      tipoReaccion: tipoReaccion,
    });

    stompClient.send("/app/enviarReaccion", {}, payload);
    console.log(`[Reacciones] Enviada reacción "${tipoReaccion}" para audio "${idAudioActual}"`);
  }

  /**
   * Desconecta del broker y limpia el estado. Se llama al cambiar de audio o de vista.
   */
  function desconectar() {
    if (suscripcionActual) {
      try { suscripcionActual.unsubscribe(); } catch (_) { /* ya desconectado */ }
      suscripcionActual = null;
    }
    if (stompClient) {
      try { stompClient.disconnect(); } catch (_) { /* ignorar */ }
      stompClient = null;
    }
    conectado = false;
    idAudioActual = null;
  }

  // ─── Componente visual ───────────────────────────────────────

  /**
   * Crea la sección de reacciones que se integra en la ficha del audio.
   * Contiene:
   *   - Estado de la conexión WebSocket
   *   - Botones de reacción (habilitados solo cuando hay audio reproduciéndose)
   *   - Zona donde aparecen las reacciones flotantes recibidas
   *
   * @param {string} idAudio  Identificador del audio (para suscripción al topic).
   * @returns {HTMLElement}   Sección lista para insertar en el DOM.
   */
  function crearSeccion(idAudio) {
    // Contenedor principal (usa la utilidad `el` definida en app.js)
    const seccion = document.createElement("section");
    seccion.className = "bg-white rounded-2xl border border-gray-200 shadow-sm p-6 space-y-4";

    // Título
    const titulo = document.createElement("h3");
    titulo.className = "text-sm font-bold uppercase tracking-wide text-gray-500";
    titulo.textContent = "Reacciones en vivo";
    seccion.appendChild(titulo);

    // Estado de conexión
    const estadoConexion = document.createElement("p");
    estadoConexion.className = "text-sm text-gray-400";
    estadoConexion.textContent = "⏳ Esperando reproducción para habilitar reacciones…";
    seccion.appendChild(estadoConexion);

    // Contenedor de botones de reacciones
    const botonesContenedor = document.createElement("div");
    botonesContenedor.className = "flex flex-wrap gap-2";

    EMOJIS.forEach((emoji) => {
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className =
        "reaccion-btn opacity-50 cursor-not-allowed text-2xl px-3 py-2 rounded-xl border border-gray-200 bg-[#f5f5f7] hover:bg-gray-100 transition-all transform hover:scale-110 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100";
      btn.textContent = emoji;
      btn.disabled = true;
      btn.title = `Reaccionar con ${emoji}`;
      btn.addEventListener("click", () => {
        if (btn.disabled) return;
        enviar(emoji);
        // Pequeña animación de "pulso" al hacer click
        btn.classList.add("scale-125");
        setTimeout(() => btn.classList.remove("scale-125"), 150);
      });
      botonesContenedor.appendChild(btn);
    });

    seccion.appendChild(botonesContenedor);

    // Zona de reacciones flotantes (las reacciones recibidas se muestran aquí)
    const zonaReacciones = document.createElement("div");
    zonaReacciones.className =
      "relative overflow-hidden rounded-xl bg-gradient-to-b from-[#f5f5f7] to-white border border-gray-100 h-24 flex items-end justify-center";
    zonaReacciones.id = "zona-reacciones";

    const placeholder = document.createElement("p");
    placeholder.className = "text-xs text-gray-300 pb-2 reacciones-placeholder";
    placeholder.textContent = "Las reacciones de todos los oyentes aparecerán aquí";
    zonaReacciones.appendChild(placeholder);

    seccion.appendChild(zonaReacciones);

    // ─── Funciones internas del componente ────────────────

    /** Habilita los botones y conecta al WebSocket */
    function activar() {
      estadoConexion.textContent = "🔌 Conectando al servidor de reacciones…";
      estadoConexion.className = "text-sm text-yellow-600";

      conectar(idAudio, (reaccion) => {
        mostrarReaccionFlotante(zonaReacciones, reaccion.tipoReaccion);
      })
        .then(() => {
          estadoConexion.textContent = "🟢 Conectado · Reacciones habilitadas";
          estadoConexion.className = "text-sm text-green-600 font-medium";
          botonesContenedor.querySelectorAll(".reaccion-btn").forEach((btn) => {
            btn.disabled = false;
            btn.classList.remove("opacity-50", "cursor-not-allowed");
            btn.classList.add("cursor-pointer");
          });
        })
        .catch(() => {
          estadoConexion.textContent = "🔴 No se pudo conectar al servidor de reacciones";
          estadoConexion.className = "text-sm text-red-500";
        });
    }

    /** Deshabilita los botones y desconecta del WebSocket */
    function desactivar() {
      desconectar();
      estadoConexion.textContent = "⏳ Esperando reproducción para habilitar reacciones…";
      estadoConexion.className = "text-sm text-gray-400";
      botonesContenedor.querySelectorAll(".reaccion-btn").forEach((btn) => {
        btn.disabled = true;
        btn.classList.add("opacity-50", "cursor-not-allowed");
        btn.classList.remove("cursor-pointer");
      });
    }

    // Exponer métodos para que app.js pueda activar/desactivar las reacciones
    seccion._reacciones = { activar, desactivar };

    return seccion;
  }

  /**
   * Muestra un emoji flotante animado en la zona de reacciones.
   */
  function mostrarReaccionFlotante(zona, emoji) {
    // Quitar el placeholder si existe
    const ph = zona.querySelector(".reacciones-placeholder");
    if (ph) ph.remove();

    const span = document.createElement("span");
    span.textContent = emoji;
    span.className = "reaccion-flotante";
    // Posición horizontal aleatoria
    span.style.left = Math.random() * 80 + 10 + "%";
    zona.appendChild(span);

    // Retirar del DOM después de la animación
    span.addEventListener("animationend", () => span.remove());
  }

  // ─── API pública ─────────────────────────────────────────────
  return {
    conectar,
    enviar,
    desconectar,
    crearSeccion,
    get conectado() { return conectado; },
  };
})();

