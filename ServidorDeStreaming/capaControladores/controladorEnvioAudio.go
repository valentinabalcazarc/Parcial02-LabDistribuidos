package capacontrollers

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	capafachadaservices "servidor.local/grpc-servidor/capaFachadaServices"
	"servidor.local/grpc-servidor/capalogger"
	pb "servidor.local/grpc-servidor/serviciosCancion"
)

type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
	logger *capalogger.Logger
}

func NewControladorServidor(logger *capalogger.Logger) *ControladorServidor {
	return &ControladorServidor{
		logger: logger,
	}
}

// Implementación del procedimiento remoto
func (thisC *ControladorServidor) EnviarCancionMedianteStream(
	req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {

	direcionCliente := ObtenerDireccionCliente(stream.Context())

	// Invocación a operación asíncrona de log
	go thisC.logger.AlmacenarSolicitud(req.Titulo, direcionCliente)

	titulo := req.GetTitulo()
	formato := req.GetFormato()
	if formato != "" && !strings.HasSuffix(titulo, "."+formato) {
		titulo = titulo + "." + formato
	}

	return capafachadaservices.StreamAudioFile(
		titulo,
		direcionCliente,
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}

/**
 * @brief Obtiene la identidad o dirección del cliente.
 *
 * Si el cliente envía 'x-client-id' en los metadatos gRPC, utiliza ese valor.
 * De lo contrario, toma la dirección IP del socket o 'cliente_anonimo'.
 */
func ObtenerDireccionCliente(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-client-id"); len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
		if vals := md.Get("cliente"); len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return "cliente_anonimo"
}
