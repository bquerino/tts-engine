# Etapa de build
FROM golang:1.23.2 AS builder

WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download dependencies and build the binary
RUN go mod tidy
RUN go build -o /app/bin/tts-engine ./cmd/server/main.go

# Executing stage
FROM alpine

WORKDIR /app

# install ca-certificates in order to access HTTPS endpoints
RUN apk add --no-cache ca-certificates libc6-compat

# copy the binary from the builder stage
COPY --from=builder /app/bin/tts-engine /app/tts-engine

# port that the container should expose
EXPOSE 8080

# init the application
CMD ["/app/tts-engine"]
