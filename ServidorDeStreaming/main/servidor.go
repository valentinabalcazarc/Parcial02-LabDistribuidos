package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	capacontroladores "servidor.local/grpc-servidor/capaControladores"
	"servidor.local/grpc-servidor/capalogger"
	pb "servidor.local/grpc-servidor/serviciosCancion"

	discovery "servidor.local/grpc-servidor/discovery"
)

const nombreServicio = "ms-servidorDeStreaming"

func main() {
	puerto := 50051
	host := env("SERVICE_HOST", "127.0.0.1")
	consulAddr := env("CONSUL_ADDR", "127.0.0.1:8500")

	// 1. Iniciar listener TCP en el puerto gRPC 50051
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(puerto))
	if err != nil {
		log.Fatalf("Error escuchando en el puerto %d: %v", puerto, err)
	}

	grpcServer := grpc.NewServer()

	// 2. Crear logger como singleton
	objLogger := capalogger.CrearUnicaInstanciaDelLogger()
	defer objLogger.CerrarArchivo()

	// 3. Registrar el controlador del servicio de streaming
	objControlador := capacontroladores.NewControladorServidor(objLogger)
	pb.RegisterAudioServiceServer(grpcServer, objControlador)

	// 4. Iniciar el servidor gRPC en una goroutine para no bloquear el registro en Consul
	go func() {
		fmt.Printf("Servidor gRPC escuchando en :%d...\n", puerto)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error al iniciar servidor gRPC: %v", err)
		}
	}()

	// 5. Registrar en el directorio de descubrimiento Consul
	registro, err := discovery.Registrar(consulAddr, nombreServicio, host, puerto)
	if err != nil {
		log.Printf("⚠️  [ADVERTENCIA] No fue posible registrarse en Consul (%s): %v", consulAddr, err)
		log.Println("   (El servidor gRPC continuará funcionando normalmente)")
	}

	// 6. Esperar señal de apagado (CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\nApagando servidor gRPC...")
	if registro != nil {
		registro.Desregistrar()
	}
	grpcServer.GracefulStop()
	fmt.Println("Servidor gRPC detenido correctamente.")
}

func env(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}
