FROM golang:latest AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# 1. Copiar solo los archivos necesarios para descargar dependencias (evita rebuilds innecesarios)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copiar el resto del código y compilar
COPY . .

ARG DB_HOST
ARG DB_USER

RUN go build -o main .

# Exponer el puerto en el que la aplicación escucha
EXPOSE 5555

# Comando para ejecutar la aplicación
CMD ["./main"]
