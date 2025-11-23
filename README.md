# Proyecto Sistemas Distribuidos: Plataforma de Micro-Reseñas

Integrantes: 

-Vicente Luongo ROL 202073637-5 

-Esteban ROL 

-Antonio ROL

Este sistema implementa una arquitectura distribuida que soporta **Consistencia Eventual** (entre DataNodes) y los modelos **Read Your Writes / Monotonic Reads** (para el Cliente), desarrollado en Go utilizando gRPC.

---

## 1. Requisitos Previos

Antes de comenzar, asegúrate de tener instalado lo siguiente en tu sistema (o en cada Máquina Virtual):

### Instalación de Go y Herramientas
1.  **Go (Golang):** Versión 1.19 o superior.
2.  **Make:** Para ejecutar los comandos de automatización.
3.  **Compilador de Protocol Buffers (Protoc):**

**En Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y golang-go make protobuf-compiler

Instalación de Plugins de Go para Protoc
Es necesario instalar los generadores de código para Go y gRPC. Ejecuta:

Bash

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
Nota: Asegúrate de que la ruta de instalación de Go esté en tu PATH. Agrega esto a tu ~/.bashrc si no te funcionan los comandos anteriores: export PATH=$PATH:$(go env GOPATH)/bin

2. Descarga y Configuración Inicial
Clonar/Descargar el repositorio y entrar en la carpeta.

Descargar dependencias del proyecto:

Bash

go mod tidy
3. Configuración de Red (CRÍTICO)
El sistema se configura a través del archivo Makefile. Debes elegir una de las dos modalidades de ejecución:

Opción A: Ejecución en una sola máquina (Localhost)
Si quieres probar todo el sistema en tu propio computador:

Abre el archivo Makefile.

Busca la sección de configuración de IPs al inicio.

Cambia TODAS las direcciones IP a 127.0.0.1.

Makefile

# Ejemplo para Localhost
VM1_IP=127.0.0.1
VM2_IP=127.0.0.1
VM3_IP=127.0.0.1
VM_COORD_IP=127.0.0.1
Opción B: Ejecución en Múltiples Máquinas Virtuales (Laboratorio)
Si vas a desplegar cada nodo en una VM distinta:

Obtén la IP de cada máquina (comando ip addr o ifconfig).

Edita el archivo Makefile en todas las máquinas para que tengan la misma configuración.

Makefile

# Ejemplo Distribuido (Reemplazar con IPs reales)
VM1_IP=192.168.1.10  # IP de la máquina que correrá el DataNode 1
VM2_IP=192.168.1.11  # IP de la máquina que correrá el DataNode 2
VM3_IP=192.168.1.12  # IP de la máquina que correrá el DataNode 3
VM_COORD_IP=192.168.1.13 # IP de la máquina que correrá el Coordinador
4. Compilación
Una vez configurado el Makefile, compila los binarios y genera el código gRPC. Este paso debe realizarse en todas las máquinas si estás en un entorno distribuido.

Bash

make build
Esto generará los ejecutables en la carpeta /bin y los archivos .pb.go en /proto.

5. Ejecución del Sistema
El orden de ejecución es importante. Debes levantar primero los DataNodes y luego el Coordinador.

Paso 1: Levantar los DataNodes
Abre 3 terminales distintas (o ve a cada VM correspondiente) y ejecuta:

Terminal/VM 1: make run-node-1

Terminal/VM 2: make run-node-2

Terminal/VM 3: make run-node-3

Paso 2: Levantar el Coordinador
Una vez que los nodos de datos estén corriendo, inicia el coordinador:

Terminal/VM 4: make run-coord

Paso 3: Ejecutar el Cliente
El cliente puede correr en cualquier terminal o máquina con acceso a la red del Coordinador:

Terminal Cliente: make run-client

6. Uso del Cliente
Al iniciar el cliente verás un menú interactivo:

Escribir Reseña (Opción 1):

Envía un texto al sistema.

Read Your Writes: El cliente guardará internamente la IP del nodo físico donde se guardó el dato.

Leer Reseña (Opción 2):

Solicita la lectura de la reseña.

Read Your Writes: El sistema intentará leer del mismo nodo donde escribiste anteriormente para asegurar consistencia inmediata.

Monotonic Reads: El cliente verificará que la versión del dato recibido no sea más antigua que la última que ha visto.

Solución de Problemas Comunes
Error protoc: command not found:

No has instalado el compilador. Revisa el paso 1 "Requisitos Previos".

Error protoc-gen-go: program not found:

Faltan los plugins de Go o tu PATH no está configurado correctamente. Ejecuta export PATH=$PATH:$(go env GOPATH)/bin e intenta compilar de nuevo.

Connection Refused / Context Deadline Exceeded:

Verifica que las IPs en el Makefile sean correctas.

Si usas varias VMs, asegúrate de que no haya un Firewall bloqueando los puertos (por defecto el sistema usa rangos 50050-50060).

Prueba hacer ping entre las máquinas para verificar conectividad básica.