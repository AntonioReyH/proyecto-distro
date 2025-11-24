## Proyecto Sistemas Distribuidos: Plataforma de Micro-Reseñas

# Integrantes:
- Vicente Luongo ROL 202073637-5

- Esteban Carrasco ROL 201773546-5

- Antonio Rey ROL 202173633-6 


Este sistema implementa una arquitectura distribuida en Go utilizando gRPC. Soporta Consistencia Eventual entre nodos de datos y garantiza los modelos Read Your Writes y Monotonic Reads para el cliente.

# Requisitos Previos
Asegúrate de tener instalado lo siguiente en tu entorno (o en cada máquina virtual):

Go (Golang): Versión 1.19 o superior.

Make: Para la automatización de tareas.

Compilador Protocol Buffers (protoc):

 

# Ubuntu/Debian
```
sudo apt update && sudo apt install -y golang-go make protobuf-compiler

```
Plugins de Go para gRPC: 

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
Asegúrate de que $GOPATH/bin esté en tu $PATH.

# 1. Configuración e Instalación
Paso 1.1: Descargar dependencias
En la raíz del proyecto, ejecuta:

 

```
go mod tidy
```
# Paso 1.2: Configurar IPs en el Makefile
El archivo Makefile controla la configuración de red. Debes editarlo según tu entorno de ejecución.

# Opción A: Modo Local (Una sola máquina) Si vas a probar todo en tu computador personal:

Abre el Makefile.

Busca la variable LOCAL_HOST (línea ~25) y asegúrate de que sea tu IP local o localhost:

Makefile
```
LOCAL_HOST=127.0.0.1
```
Nota: El archivo original tiene 172.22.87.38, cámbialo si estás probando en tu propia máquina.


# Opción B: Modo Distribuido (4 Máquinas Virtuales) Si tienes 4 máquinas (3 DataNodes + 1 Coordinador):

Abre el Makefile.

Edita la sección "1. MODO DISTRIBUIDO" con las IPs reales de tus VMs:

Makefile
```
VM1_IP=192.168.1.10  # IP para DataNode 1
VM2_IP=192.168.1.11  # IP para DataNode 2
VM3_IP=192.168.1.12  # IP para DataNode 3
VM_COORD_IP=192.168.1.20 # IP para el Coordinador
```
# 2. Compilación
Este paso genera los binarios en la carpeta bin/ y el código de gRPC en proto/. Ejecuta este comando (en todas las máquinas si es distribuido):

 
```
make build
```
# 3. Ejecución
Selecciona la guía según el modo que configuraste. El orden es importante: primero los DataNodes, luego el Coordinador, y al final el Cliente.

# Guía para MODO LOCAL (1 Máquina)
Necesitarás abrir 5 terminales diferentes en la raíz del proyecto.

Terminal 1 (DataNode 1):

 
```
make run-local-node-1
```
Terminal 2 (DataNode 2):

 
```
make run-local-node-2
```
Terminal 3 (DataNode 3):

 

```
make run-local-node-3
```
Terminal 4 (Coordinador):

 
```
make run-local-coord
```
Terminal 5 (Cliente):


```
make run-local-client
```
# Guía para MODO DISTRIBUIDO (4 VMs)
Ejecuta el comando correspondiente en cada máquina virtual según su rol.

En dist078 (DataNode 1):

 
```
make run-dist-node-1
En dist079 (DataNode 2):
```

 
```
make run-dist-node-2
```
En dist080 (DataNode 3):

 
```
make run-dist-node-3
```
En dist102 (Coordinador): Asegúrate de que los DataNodes ya estén corriendo.

 
```
make run-dist-coord
```
Cliente (Cualquier máquina con acceso al Coordinador),
En nuestro caso lo ejecutamos en la dist078:

 
```
make run-dist-client
```
4. Uso del Cliente
Una vez iniciado el cliente, verás el siguiente menú:

Escribir Reseña (Opción 1):

Envía una nueva reseña al sistema.

El cliente guarda internamente la IP del nodo que procesó la escritura ("Read Your Writes").

Leer Reseña (Opción 2):

Solicita leer la reseña.

Consistencia Read Your Writes: El coordinador intentará redirigir tu lectura al mismo nodo donde escribiste por última vez.

Consistencia Monotonic Reads: El cliente verificará si la versión del dato (Reloj Vectorial) es posterior o igual a la última que vio. Si el dato es antiguo, mostrará una advertencia en los logs del servidor.

Limpieza
Para borrar los binarios y archivos generados por protobuf:

 
```
make clean
```