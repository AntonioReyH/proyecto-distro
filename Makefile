# ==========================================
# 1. MODO DISTRIBUIDO (Multiples VMs)
# ==========================================
# Configura estas IPs si usas 4 máquinas distintas
VM1_IP=10.35.168.88			#dist078: Datanode 1
VM2_IP=10.35.168.89 		#dist079: Datanode 2
VM3_IP=10.35.168.90			#dist080: Datanode 3
VM_COORD_IP=10.35.168.112	#dist102: Coordinador

# En distribuido, todos pueden usar el mismo puerto pq son distintas IPs
PORT_DIST=50051
PORT_COORD_DIST=50050

PEERS_DIST_1=$(VM2_IP):$(PORT_DIST),$(VM3_IP):$(PORT_DIST)
PEERS_DIST_2=$(VM1_IP):$(PORT_DIST),$(VM3_IP):$(PORT_DIST)
PEERS_DIST_3=$(VM1_IP):$(PORT_DIST),$(VM2_IP):$(PORT_DIST)
ALL_DATANODES_DIST=$(VM1_IP):$(PORT_DIST),$(VM2_IP):$(PORT_DIST),$(VM3_IP):$(PORT_DIST)
COORD_ADDR_DIST=$(VM_COORD_IP):$(PORT_COORD_DIST)

# ==========================================
# 2. MODO LOCAL (Una sola máquina, N terminales)
# ==========================================
LOCAL_HOST=172.192.168.1
# Puertos diferentes para evitar choques en localhost
PORT_L_1=50051
PORT_L_2=50052
PORT_L_3=50053
PORT_L_COORD=50050

PEERS_LOCAL_1=$(LOCAL_HOST):$(PORT_L_2),$(LOCAL_HOST):$(PORT_L_3)
PEERS_LOCAL_2=$(LOCAL_HOST):$(PORT_L_1),$(LOCAL_HOST):$(PORT_L_3)
PEERS_LOCAL_3=$(LOCAL_HOST):$(PORT_L_1),$(LOCAL_HOST):$(PORT_L_2)
ALL_DATANODES_LOCAL=$(LOCAL_HOST):$(PORT_L_1),$(LOCAL_HOST):$(PORT_L_2),$(LOCAL_HOST):$(PORT_L_3)
COORD_ADDR_LOCAL=$(LOCAL_HOST):$(PORT_L_COORD)


# ==========================================
# COMANDOS DE COMPILACIÓN (Comunes)
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
# EJECUCIÓN: MODO DISTRIBUIDO (4 VMs)
# ==========================================
run-dist-node-1: build
	@echo "--- MODO DISTRIBUIDO: DataNode 1 ---"
	NODE_ID="node1" PORT="$(PORT_DIST)" PEERS="$(PEERS_DIST_1)" ./bin/datanode

run-dist-node-2: build
	@echo "--- MODO DISTRIBUIDO: DataNode 2 ---"
	NODE_ID="node2" PORT="$(PORT_DIST)" PEERS="$(PEERS_DIST_2)" ./bin/datanode

run-dist-node-3: build
	@echo "--- MODO DISTRIBUIDO: DataNode 3 ---"
	NODE_ID="node3" PORT="$(PORT_DIST)" PEERS="$(PEERS_DIST_3)" ./bin/datanode

run-dist-coord: build
	@echo "--- MODO DISTRIBUIDO: Coordinador ---"
	DATANODES="$(ALL_DATANODES_DIST)" ./bin/coordinator

run-dist-client: build
	@echo "--- MODO DISTRIBUIDO: Cliente ---"
	COORDINATOR_ADDR="$(COORD_ADDR_DIST)" ./bin/client


# ==========================================
# EJECUCIÓN: MODO LOCAL (1 VM, Multi-Terminal)
# ==========================================
run-local-node-1: build
	@echo "--- MODO LOCAL: DataNode 1 (Puerto $(PORT_L_1)) ---"
	# ANTES: NODE_ID="node1" ...
	# AHORA:
	NODE_ID="$(LOCAL_HOST):$(PORT_L_1)" PORT="$(PORT_L_1)" PEERS="$(PEERS_LOCAL_1)" ./bin/datanode

run-local-node-2: build
	@echo "--- MODO LOCAL: DataNode 2 (Puerto $(PORT_L_2)) ---"
	# ANTES: NODE_ID="node2" ...
	# AHORA:
	NODE_ID="$(LOCAL_HOST):$(PORT_L_2)" PORT="$(PORT_L_2)" PEERS="$(PEERS_LOCAL_2)" ./bin/datanode

run-local-node-3: build
	@echo "--- MODO LOCAL: DataNode 3 (Puerto $(PORT_L_3)) ---"
	# ANTES: NODE_ID="node3" ...
	# AHORA:
	NODE_ID="$(LOCAL_HOST):$(PORT_L_3)" PORT="$(PORT_L_3)" PEERS="$(PEERS_LOCAL_3)" ./bin/datanode

run-local-coord: build
	@echo "--- MODO LOCAL: Coordinador (Puerto $(PORT_L_COORD)) ---"
	DATANODES="$(ALL_DATANODES_LOCAL)" ./bin/coordinator

run-local-client: build
	@echo "--- MODO LOCAL: Cliente ---"
	COORDINATOR_ADDR="$(COORD_ADDR_LOCAL)" ./bin/client