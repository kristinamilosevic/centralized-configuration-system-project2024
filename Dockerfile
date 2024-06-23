# Stage 1: Build stage
FROM golang:1.19 AS builder

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

# # Instalacija potrebnih alata
# RUN apt-get update && apt-get install -y curl unzip

# # Instalacija Consul agenta
# RUN apt-get update && apt-get install -y curl
# RUN curl -L -o consul.zip https://releases.hashicorp.com/consul/1.11.0/consul_1.11.0_linux_amd64.zip \
#     && unzip consul.zip \
#     && rm consul.zip \
#     && mv consul /usr/local/bin/consul

# Postavljanje radnog direktorijuma za izvršavanje aplikacije
WORKDIR /app

# Kopiranje izvrsnog fajla iz prethodnog stage-a
COPY --from=builder /app/main .

# Kopiranje konfiguracionih fajlova za Consul (opciono)
COPY consul-config.json /etc/consul/config.json

# Portovi za Consul
EXPOSE 8500
EXPOSE 8600

# Port na kojem će aplikacija slušati zahtjeve
EXPOSE 8000

# Komanda za pokretanje aplikacije i Consul agenta
CMD ["sh", "-c", "./main & consul agent -dev -config-dir=/etc/consul"]
