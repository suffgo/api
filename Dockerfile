# -------- Build stage --------
FROM golang:1.24.2 AS builder
WORKDIR /app

# Solo copiamos los módulos para cachear capas de dependencias\ nCOPY api/go.mod api/go.sum ./
RUN go mod download

# Copiamos el código
COPY api/api/ ./

# Construimos un binario estático optimizado para Linux amd64
RUN CGO_ENABLED=0 \ 
    GOOS=linux \ 
    GOARCH=amd64 \ 
    go build -ldflags="-s -w" -o suffgo-api ./cmd/app

# -------- Final stage --------
FROM alpine:3.21 AS runner

# Certificados para conexiones HTTPS
RUN apk add --no-cache ca-certificates
WORKDIR /root/

# Copiamos el binario desde builder
COPY --from=builder /app/suffgo-api .

# Puerto que la aplicación expone (se controla vía $PORT en producción)
EXPOSE 10000

# Arrancamos la API
ENTRYPOINT ["./suffgo-api"]
