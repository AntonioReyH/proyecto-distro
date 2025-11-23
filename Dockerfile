FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Instalar compilador de protobuf si no se generan localmente, 
# pero mejor usar el makefile local y copiar los .go generados.
# Asumiremos que ejecutas 'make gen' antes de docker build o que copiamos todo.
RUN go build -o /bin/datanode datanode/main.go
RUN go build -o /bin/coordinator coordinator/main.go
RUN go build -o /bin/client client/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /bin/datanode .
COPY --from=builder /bin/coordinator .
COPY --from=builder /bin/client .

# Exponemos puertos típicos
EXPOSE 50051 50052 50053 50050