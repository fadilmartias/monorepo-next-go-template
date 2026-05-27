# Stage 1: Build tahap kompilasi
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk caching module
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy seluruh source code
COPY backend/ .

# Build aplikasi Go
# Build Web Server
RUN CGO_ENABLED=0 GOOS=linux go build -o server .
# Build CLI
RUN CGO_ENABLED=0 GOOS=linux go build -o cli ./cmd/main.go

# Stage 2: Hasil akhir yang sangat ringan
FROM alpine:latest

WORKDIR /app

# PERBAIKAN: Gunakan ./ (titik slash) agar dicopy ke DALAM folder /app
COPY --from=builder /app/server ./server
COPY --from=builder /app/cli ./cli
COPY --from=builder /app/public ./public

# PERBAIKAN: Tidak perlu copy storage dari builder. 
# Cukup buat struktur foldernya saja di sini agar tidak error saat start.
# Nanti isinya akan diisi oleh volume dari docker-compose.
RUN mkdir -p ./storage/logs ./storage/uploads/tmp ./storage/uploads/tinymce

EXPOSE 9001
CMD ["./server"]