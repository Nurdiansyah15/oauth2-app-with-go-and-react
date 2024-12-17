# Menggunakan image Golang sebagai base image
FROM golang:1.23 AS builder

# Set work directory di dalam container
WORKDIR /app

# Salin file modul dan dependency
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Salin semua kode ke dalam container
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage kedua: Menggunakan alpine untuk image yang lebih kecil
FROM alpine:latest

# Install ca-certificates untuk mendukung HTTPS
RUN apk --no-cache add ca-certificates

# Set work directory di dalam container
WORKDIR /root/

# Salin file binary dari stage builder
COPY --from=builder /app/main .

COPY ./templates ./templates

COPY ./public/assets ./public/assets

EXPOSE 8080


# Jalankan aplikasi
CMD ["./main"]
