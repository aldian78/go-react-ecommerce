FROM golang:1.23.0-alpine


WORKDIR /app

# Copy semua file dari root proyek (bukan hanya dari docker-compose/)
COPY .. .

# Get dependency (jika diperlukan)
RUN go mod download

# Build REST services
RUN go build -o rest ./cmd/rest/main.go

# Jalankan binary
CMD ["./rest"]
