# Stage 1: Build stage
FROM golang:1.17 AS builder

# Postavljanje radnog direktorijuma u kontejneru
WORKDIR /app

# Kopiranje Go modula i dependencija
COPY go.mod go.sum ./
RUN go mod download

# Kopiranje izvornog koda aplikacije
COPY . .

# Buildovanje izvrsnog fajla
RUN go build -o main .

# Stage 2: Production stage
FROM debian:buster-slim AS final

# Postavljanje radnog direktorijuma za izvršavanje aplikacije
WORKDIR /app

# Kopiranje izvrsnog fajla iz prethodnog stage-a
COPY --from=builder /app/main .

# Port na kojem će aplikacija slušati zahtjeve
EXPOSE 8000

# Komanda za pokretanje aplikacije
CMD ["./main"]
