package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"proyecto-sd/common"
	pb "proyecto-sd/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataNode struct {
	pb.UnimplementedDistributedServiceServer
	ID        string
	Port      string
	Peers     []string // IPs de otros datanodes
	Storage   map[string]*pb.Review
	mu        sync.Mutex
}

func (s *DataNode) WriteData(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Crear nueva reseña o actualizar
	reviewID := "review-1" // Simplificación: una sola reseña para demo
	
	// Inicializar reloj si no existe
	currentReview, exists := s.Storage[reviewID]
	var newClock *pb.VectorClock
	
	if exists {
		newClock = common.MergeClocks(currentReview.Clock, nil)
	} else {
		newClock = &pb.VectorClock{Versions: make(map[string]int64)}
	}

	// Incrementar reloj lógico propio
	newClock.Versions[s.ID]++

	newReview := &pb.Review{
		Id:        reviewID,
		Content:   req.Content,
		Clock:     newClock,
		Timestamp: time.Now().UnixNano(),
	}

	s.Storage[reviewID] = newReview
	log.Printf("Escritura local: %s, Clock: %v", req.Content, newClock.Versions)

	// Trigger asíncrono de replicación (Eventual Consistency)
	go s.broadcastReplication(newReview)

	return &pb.WriteResponse{ReviewId: reviewID, WrittenAtNode: s.ID}, nil
}

func (s *DataNode) ReadData(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, exists := s.Storage[req.ReviewId]
	if !exists {
		return nil, fmt.Errorf("review not found")
	}

	// Chequeo Monotonic Reads: Si el dato local es más viejo que lo que el cliente ya vio
	if !common.IsAfter(review.Clock, req.MinClock) {
		// En un sistema real esperaríamos o rechazaríamos.
		// Aquí retornamos lo que hay pero logueamos la advertencia.
		log.Printf("Warning: Violación potencial de Monotonic Reads. Local clock antiguo.")
	}

	return &pb.ReadResponse{Review: review, SourceNode: s.ID}, nil
}

func (s *DataNode) Replicate(ctx context.Context, req *pb.ReplicateRequest) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	incoming := req.Review
	current, exists := s.Storage[incoming.Id]

	log.Printf("Recibiendo replicación de %s. Clock entrante: %v", req.SenderNodeId, incoming.Clock.Versions)

	if !exists {
		s.Storage[incoming.Id] = incoming
		log.Println("Datos nuevos aceptados.")
		return &pb.Empty{}, nil
	}

	// Resolución de conflictos y Convergencia
	isIncomingNewer := common.IsAfter(incoming.Clock, current.Clock)
	isCurrentNewer := common.IsAfter(current.Clock, incoming.Clock)

	if isIncomingNewer {
		s.Storage[incoming.Id] = incoming
		log.Println("Actualización aceptada (Clock posterior).")
	} else if !isCurrentNewer {
		// Conflicto (concurrentes). Usamos Timestamp (Last Writer Wins) para converger.
		if incoming.Timestamp > current.Timestamp {
			// Merge clocks para capturar la causalidad de ambas ramas
			mergedClock := common.MergeClocks(current.Clock, incoming.Clock)
			incoming.Clock = mergedClock
			s.Storage[incoming.Id] = incoming
			log.Println("Conflicto resuelto: Gana entrante (LWW). Clock mergeado.")
		} else {
			// Actualizar solo el reloj del local para reflejar que conoce la otra rama
			current.Clock = common.MergeClocks(current.Clock, incoming.Clock)
			log.Println("Conflicto resuelto: Gana local (LWW). Clock mergeado.")
		}
	} else {
		log.Println("Replicación ignorada (Datos viejos).")
	}

	return &pb.Empty{}, nil
}

func (s *DataNode) broadcastReplication(review *pb.Review) {
	// Simular latencia aleatoria para demostrar consistencia eventual
	time.Sleep(2 * time.Second) 
	
	for _, peer := range s.Peers {
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Error conectando a peer %s: %v", peer, err)
			continue
		}
		client := pb.NewDistributedServiceClient(conn)
		_, err = client.Replicate(context.Background(), &pb.ReplicateRequest{
			Review:       review,
			SenderNodeId: s.ID,
		})
		if err != nil {
			log.Printf("Error replicando a %s: %v", peer, err)
		}
		conn.Close()
	}
}

func main() {
	id := os.Getenv("NODE_ID")
	port := os.Getenv("PORT")
	peersEnv := os.Getenv("PEERS") // Separados por coma: "ip:port,ip:port"
	
	peers := []string{}
	if peersEnv != "" {
		peers = strings.Split(peersEnv, ",")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	node := &DataNode{
		ID:      id,
		Port:    port,
		Peers:   peers,
		Storage: make(map[string]*pb.Review),
	}

	pb.RegisterDistributedServiceServer(grpcServer, node)
	log.Printf("DataNode %s escuchando en puerto %s", id, port)
	grpcServer.Serve(lis)
}