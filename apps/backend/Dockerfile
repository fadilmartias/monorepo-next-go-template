# Stage 1: Build tahap kompilasi
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk caching module
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy seluruh source code

COPY backend/ .

# Build aplikasi Go
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Hasil akhir yang sangat ringan
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/public /public
COPY --from=builder /app/storage /storage

COPY widget/ ./widget

EXPOSE 3090
CMD ["./main"]