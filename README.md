Proyecto Sistemas Distribuidos: Plataforma de Micro-Reseñas
Integrantes: Vicente Luongo ROL 202073637-5, Esteban ROL, Antonio ROL

Este proyecto implementa un sistema distribuido que soporta Consistencia Eventual (entre DataNodes) y Read Your Writes / Monotonic Reads (para el Cliente), utilizando Go y gRPC.

1. Requisitos Previos (Instalación de Herramientas)
Antes de compilar, asegúrate de que cada Máquina Virtual (VM) tenga instaladas las herramientas necesarias.

Ejecuta los siguientes comandos en la terminal de todas las VMs (asumiendo Ubuntu/Debian):

Bash

# 1. Actualizar repositorios
sudo apt-get update

# 2. Instalar Go (Golang) y Make
sudo apt-get install -y golang-go make

# 3. Instalar el compilador de Protocol Buffers
sudo apt-get install -y protobuf-compiler

# 4. Instalar plugins de Go para Protobuf y gRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 5. Asegurar que los binarios de Go estén en el PATH
export PATH="$PATH:$(go env GOPATH)/bin"
2. Inicialización del Proyecto
Si acabas de copiar y pegar el código en la carpeta, es probable que necesites inicializar el módulo de Go y descargar las dependencias.

Ejecuta esto en la carpeta raíz del proyecto (proyecto-sd/):

Bash

# Inicializar dependencias
go mod tidy
3. Configuración de Red (IMPORTANTE)
Para que las máquinas se comuniquen, debes editar el archivo Makefile con las direcciones IP reales de tus VMs.

Abre el archivo Makefile:

Bash

nano Makefile
Busca la sección CONFIGURACIÓN DE RED al inicio del archivo.

Modifica las siguientes variables con las IPs correspondientes a tu laboratorio:

Makefile

VM1_IP=192.168.X.X  <-- IP donde correrá el DataNode 1
VM2_IP=192.168.X.X  <-- IP donde correrá el DataNode 2
VM3_IP=192.168.X.X  <-- IP donde correrá el DataNode 3
VM_COORD_IP=192.168.X.X <-- IP donde correrá el Coordinador
Guarda los cambios (Ctrl+O, Enter) y sal (Ctrl+X).

Nota: Debes realizar este cambio en el Makefile de todas las máquinas virtuales para que tengan la misma configuración.

4. Compilación y Generación de Código
Antes de ejecutar los nodos, necesitas generar el código gRPC y compilar los binarios.

Ejecuta el siguiente comando en todas las máquinas:

Bash

# Genera los archivos .pb.go y compila los binarios en la carpeta /bin
make build
5. Ejecución del Sistema
El sistema consta de 4 roles distintos. Ejecuta cada comando en su VM correspondiente.

En la Máquina Virtual 1 (DataNode 1)
Bash

make run-node-1
En la Máquina Virtual 2 (DataNode 2)
Bash

make run-node-2
En la Máquina Virtual 3 (DataNode 3)
Bash

make run-node-3
En la Máquina Virtual 4 (Coordinador)
Nota: Asegúrate de que los 3 DataNodes ya estén corriendo antes de iniciar el coordinador.

Bash

make run-coord
6. Ejecución del Cliente
Puedes correr el cliente en cualquiera de las máquinas (o en tu máquina local si tiene acceso a la red de las VMs).

Bash

make run-client
Sigue las instrucciones en pantalla para:

Escribir (Opción 1): Envía una reseña. El sistema te confirmará en qué nodo físico se guardó.

Leer (Opción 2): Solicita la reseña. El sistema intentará leer del mismo nodo donde escribiste (Read Your Writes) y verificará que el reloj lógico no sea antiguo (Monotonic Reads).

Solución de Problemas Comunes
Error "command not found: protoc": No instalaste el compilador. Revisa el paso 1.

Error de conexión (Connection Refused/Timeout):

Verifica que las IPs en el Makefile sean correctas.

Verifica que los firewalls de las VMs permitan tráfico en los puertos 50050 y 50051.

Error "protoc-gen-go: program not found": Asegúrate de haber ejecutado el export PATH del paso 1.