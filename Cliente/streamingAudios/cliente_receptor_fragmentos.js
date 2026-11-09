// cliente_receptor_fragmentos.js

const { peticionDTO } = require('../stubs_generados/servicios_pb.js');
const { AudioServicePromiseClient } = require('../stubs_generados/servicios_grpc_web_pb.js');

const ENVOY_HOST = 'http://localhost:8090';

const client = new AudioServicePromiseClient(ENVOY_HOST);

// --- Utilidades del elemento <audio> ---

/** Libera el Blob que tenga cargado el reproductor (cada Blob URL ocupa memoria hasta que se revoca). */
function liberarUrl(reproductor) {
    if (reproductor.dataset.urlObjeto) {
        URL.revokeObjectURL(reproductor.dataset.urlObjeto);
        delete reproductor.dataset.urlObjeto;
    }
}

/** Detiene el audio del reproductor y lo deja vacío. */
function liberarReproductor(reproductor) {
    if (!reproductor) return;
    reproductor.pause();
    liberarUrl(reproductor);
    reproductor.removeAttribute('src');
    reproductor.load();
}

/**
 * Pide una canción al servidor de streaming (gRPC-Web a través de Envoy), recibe sus fragmentos y,
 * cuando termina la transmisión, la carga en el reproductor y la reproduce.
 *
 * @param {string} titulo   Nombre del archivo a pedir (p. ej. "1.mp3").
 * @param {string} formato  Formato del audio (p. ej. "mp3").
 * @param {Object} [opciones]  Opcional; sin opciones se comporta con los valores por defecto.
 * @param {HTMLAudioElement} [opciones.reproductor]  <audio> donde reproducir; por omisión #audio-player.
 * @param {(datos: Uint8Array) => void} [opciones.onFragmento]  Se invoca con cada fragmento recibido.
 * @param {() => void} [opciones.onRecibido]  Se invoca cuando llegó todo el audio y ya está en el reproductor.
 * @param {(error: {codigo: number, mensaje: string}) => void} [opciones.onError]  Fallo en la transmisión.
 * @returns {{detener: () => void}}  Permite detener la recepción y la reproducción.
 */
function iniciar_streaming_cancion(titulo, formato, opciones = {}) {

    console.log("Iniciando llamada gRPC-Web a través de Envoy...");

    // El <audio> se busca al pedir la canción (en audios.html se crea al mostrar la ficha del audio).
    const reproductor = opciones.reproductor || document.getElementById('audio-player');

    // Estado propio de esta transmisión (así dos transmisiones no mezclan sus fragmentos).
    const audioChunks = [];
    let primerFragmento = true;
    let terminada = false; // ya hubo 'end', 'error' o se detuvo: ignorar cualquier evento posterior

    const request = new peticionDTO();
    request.setTitulo(titulo);
    request.setFormato(formato);

    console.log("Solicitud preparada:", request.toObject());
    const stream = client.enviarCancionMedianteStream(request, {});

    // Maneja la llegada de cada fragmento
    stream.on('data', (response) => {
        if (terminada) return;

        // Obtiene el array de bytes (Uint8Array)
        const data = response.getData_asU8();
        console.log(`Fragmento recibido desde el servidor. Tamaño: ${data.length} bytes`);

        // Almacena el fragmento en el array
        audioChunks.push(data);

        // Si es el primer fragmento, mostramos el reproductor y escribimos en el log
        if (primerFragmento) {
            primerFragmento = false;
            if (reproductor) {
                reproductor.classList.remove('hidden');
            }
            if (typeof window !== 'undefined' && window.__writeAudioLog) {
                window.__writeAudioLog('Primer fragmento recibido. Reproductor mostrado.', 'success');
            } else {
                console.log('Primer fragmento recibido. Reproductor mostrado.');
            }
        }

        if (opciones.onFragmento) opciones.onFragmento(data);
    });

    // Maneja el final del stream
    stream.on('end', () => {
        if (terminada) return;
        terminada = true;

        if (audioChunks.length === 0) {
            console.error('El servidor cerró la transmisión sin enviar datos de audio.');
            if (opciones.onError) opciones.onError({ codigo: 0, mensaje: 'El servidor no envió datos de audio.' });
            return;
        }

        console.log('✅ Transmisión de la canción finalizada. Preparando reproducción...');

        // Combina todos los fragmentos en un solo Blob
        const audioBlob = new Blob(audioChunks, { type: 'audio/mp3' });
        audioChunks.length = 0;

        // Crea una URL de objeto para el Blob
        const audioUrl = URL.createObjectURL(audioBlob);

        // Asigna la URL al reproductor y comienza la reproducción
        if (reproductor) {
            liberarUrl(reproductor); // el Blob de una reproducción anterior ya no hace falta
            reproductor.dataset.urlObjeto = audioUrl;
            reproductor.src = audioUrl;
            reproductor.play()
                .then(() => {
                    console.log('▶️ Reproducción de la canción iniciada con éxito.');
                })
                .catch(error => {
                    console.error('❌ Error al iniciar la reproducción automática. Intente dar clic en el reproductor.', error);
                });
        }

        if (opciones.onRecibido) opciones.onRecibido();
    });

    stream.on('error', (err) => {
        if (terminada) return;
        terminada = true;
        audioChunks.length = 0;

        console.error('ERROR en la comunicación gRPC-Web:', err);

        if (opciones.onError) opciones.onError({ codigo: err.code, mensaje: err.message });
    });

    return {
        // Corta la recepción (si sigue en curso) y detiene la reproducción.
        detener() {
            terminada = true;
            audioChunks.length = 0;
            stream.cancel();
            liberarReproductor(reproductor);
        },
    };
}

// Export implementation to global window for consumers (audios.html / funcion.js)

if (typeof window !== 'undefined') {

    window.iniciar_streaming_cancion = iniciar_streaming_cancion;

    // backward compatible name
    window.iniciarStreamGRPCImpl = iniciar_streaming_cancion;

    console.log('Exported iniciar_streaming_cancion to window');

}
