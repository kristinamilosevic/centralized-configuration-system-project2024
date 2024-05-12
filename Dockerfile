# Stage 1: Build stage
FROM golang:1.17 AS builder

# Postavljanje radnog direktorijuma u kontejneru
WORKDIR /app

# Kopiranje Go modula i dependencija
COPY go.mod go.sum ./
RUN go mod download

# Kopiranje izvornog koda aplikacije
COPY . .

# Stage 2: Production stage
FROM golang:1.17 AS final

# Postavljanje radnog direktorijuma za izvršavanje aplikacije
WORKDIR /app

# Kopiranje izvornog koda aplikacije iz prethodnog stage-a
COPY --from=builder /app .

# Port na kojem će aplikacija slušati zahtjeve
EXPOSE 8000

# Komanda za pokretanje aplikacije
CMD ["go", "run", "main.go"]
