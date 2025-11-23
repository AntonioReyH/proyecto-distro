package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	//"time"

	pb "proyecto-sd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Coordinator struct {
	pb.UnimplementedDistributedServiceServer
	Datanodes []string
}

func (c *Coordinator) CreateReview(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	// Elegir un datanode aleatorio (Load Balancing)
	target := c.Datanodes[rand.Intn(len(c.Datanodes))]
	
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewDistributedServiceClient(conn)
	resp, err := client.WriteData(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// El coordinador retorna el ID del nodo que escribió para habilitar RYW
	return resp, nil
}

func (c *Coordinator) GetReview(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	var target string
	
	// Lógica Read Your Writes:
	// Si el cliente dice "mi última escritura fue en X", intentamos leer de X.
	if req.PreferredNode != "" {
		target = req.PreferredNode // Usar la IP del datanode preferido
		log.Printf("RYW: Redirigiendo lectura a nodo preferido %s", target)
	} else {
		target = c.Datanodes[rand.Intn(len(c.Datanodes))]
		log.Printf("Lectura balanceada a %s", target)
	}

	// Intentar leer
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Fallback simple si falla la conexión
		return nil, err
	}
	defer conn.Close()

	client := pb.NewDistributedServiceClient(conn)
	return client.ReadData(ctx, req)
}

func main() {
	nodesEnv := os.Getenv("DATANODES") // "ip1:port,ip2:port..."
	nodes := strings.Split(nodesEnv, ",")

	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDistributedServiceServer(s, &Coordinator{Datanodes: nodes})
	
	log.Println("Coordinador escuchando en :50050")
	s.Serve(lis)
}