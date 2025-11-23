package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "proyecto-sd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	coordinatorAddr := os.Getenv("COORDINATOR_ADDR")
	conn, err := grpc.Dial(coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDistributedServiceClient(conn)

	// Estado del Cliente
	var lastWrittenNode string
	var maxSeenClock *pb.VectorClock

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Escribir Reseña (RYW)")
		fmt.Println("2. Leer Reseña (Monotonic + RYW)")
		fmt.Print("Seleccione opción: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "1" {
			fmt.Print("Ingrese contenido: ")
			content, _ := reader.ReadString('\n')
			content = strings.TrimSpace(content)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			resp, err := c.CreateReview(ctx, &pb.WriteRequest{Content: content, ClientId: "client-1"})
			cancel()
			
			if err != nil {
				log.Printf("Error escribiendo: %v", err)
			} else {
				// Guardamos el nodo donde escribimos para pedir RYW en la próxima lectura
				lastWrittenNode = resp.WrittenAtNode // Esta debería ser la dirección IP:Port del datanode
				// Nota: En una impl real, el coordinador debería devolver la IP pública del datanode
				// Aquí asumimos que resp.WrittenAtNode es suficiente para identificarlo o mapearlo
				fmt.Printf("Escrito exitosamente en nodo: %s\n", lastWrittenNode)
			}

		} else if text == "2" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			
			// Solicitud con Preferencia (RYW) y Restricción (Monotonic)
			// Nota: Para que RYW funcione en VMs distintas, lastWrittenNode debe ser una dirección alcanzable por el Coordinador
			req := &pb.ReadRequest{
				ReviewId:      "review-1",
				PreferredNode: lastWrittenNode, // Enviamos donde escribimos la última vez
				MinClock:      maxSeenClock,    // Enviamos lo último que vimos
			}
			
			resp, err := c.GetReview(ctx, req)
			cancel()

			if err != nil {
				log.Printf("Error leyendo: %v", err)
			} else {
				fmt.Printf("Leído: '%s' desde nodo %s\n", resp.Review.Content, resp.SourceNode)
				fmt.Printf("Reloj: %v\n", resp.Review.Clock.Versions)
				
				// Actualizamos nuestro reloj Monotónico
				// En Go simple sin librería: si maxSeenClock es nil o resp.Clock es mayor, actualizar.
				maxSeenClock = resp.Review.Clock
			}
		}
	}
}