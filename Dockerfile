# ==========================================
# ETAPA 1: Compilación (Builder)
# ==========================================
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# ==========================================
# ETAPA 2: Imagen Final de Producción (Runtime)
# ==========================================
FROM alpine:latest

WORKDIR /app


RUN apk --no-cache add ca-certificates


COPY --from=builder /app/server .

COPY --from=builder /app/dist ./dist

EXPOSE 8000

# Comando por defecto para iniciar el servidor
CMD ["./server"]