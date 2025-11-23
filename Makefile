# ==========================================
# CONFIGURACIÓN DE RED (VARIABLES DE ENTORNO)
# ==========================================

# Edita estas IPs con las direcciones reales de tus 4 Máquinas Virtuales
VM1_IP=192.168.1.10
VM2_IP=192.168.1.11
VM3_IP=192.168.1.12
VM_COORD_IP=192.168.1.20

# Puertos (Pueden ser el mismo si están en máquinas distintas)
PORT_DN=50051
PORT_COORD=50050

# Listas de Peers (Calculadas automáticamente en base a las IPs de arriba)
# Cada nodo necesita saber la dirección de los otros dos
PEERS_NODE_1=$(VM2_IP):$(PORT_DN),$(VM3_IP):$(PORT_DN)
PEERS_NODE_2=$(VM1_IP):$(PORT_DN),$(VM3_IP):$(PORT_DN)
PEERS_NODE_3=$(VM1_IP):$(PORT_DN),$(VM2_IP):$(PORT_DN)

# Lista completa de nodos para el Coordinador
ALL_DATANODES=$(VM1_IP):$(PORT_DN),$(VM2_IP):$(PORT_DN),$(VM3_IP):$(PORT_DN)

# Dirección del Coordinador para el Cliente
COORD_ADDR=$(VM_COORD_IP):$(PORT_COORD)

# ==========================================
# COMANDOS DE COMPILACIÓN
# ==========================================

gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/service.proto

build: gen
	go build -o bin/datanode datanode/main.go
	go build -o bin/coordinator coordinator/main.go
	go build -o bin/client client/main.go

clean:
	rm -rf bin/ proto/*.pb.go

# ==========================================
# EJECUCIÓN POR MÁQUINA (Usar en cada VM)
# ==========================================

# --- Ejecutar en la Máquina Virtual 1 ---
run-node-1: build
	@echo "Iniciando DataNode 1 en $(VM1_IP)..."
	NODE_ID="node1" \
	PORT="$(PORT_DN)" \
	PEERS="$(PEERS_NODE_1)" \
	./bin/datanode

# --- Ejecutar en la Máquina Virtual 2 ---
run-node-2: build
	@echo "Iniciando DataNode 2 en $(VM2_IP)..."
	NODE_ID="node2" \
	PORT="$(PORT_DN)" \
	PEERS="$(PEERS_NODE_2)" \
	./bin/datanode

# --- Ejecutar en la Máquina Virtual 3 ---
run-node-3: build
	@echo "Iniciando DataNode 3 en $(VM3_IP)..."
	NODE_ID="node3" \
	PORT="$(PORT_DN)" \
	PEERS="$(PEERS_NODE_3)" \
	./bin/datanode

# --- Ejecutar en la Máquina Virtual del Coordinador ---
run-coord: build
	@echo "Iniciando Coordinador en $(VM_COORD_IP)..."
	DATANODES="$(ALL_DATANODES)" \
	./bin/coordinator

# --- Ejecutar Cliente (En cualquier máquina) ---
run-client: build
	@echo "Conectando cliente al coordinador en $(COORD_ADDR)..."
	COORDINATOR_ADDR="$(COORD_ADDR)" \
	./bin/client