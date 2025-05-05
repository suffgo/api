# -------- Build stage --------
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1) Copiamos sólo los módulos (para cachear deps)
COPY api/go.mod api/go.sum ./
RUN go mod download

# 2) Ahora copiamos TODO el código
COPY api/ ./

# 3) Construimos el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
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
