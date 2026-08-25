# --- STAGE 1: Build the binary ---
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app statically linked for a scratch/alpine environment
# CGO_ENABLED=0 disables dynamic C libraries for maximum portability
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# --- STAGE 2: Create the minimal runtime image ---
FROM scratch

# Set working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp /myapp

# Copy SSL certificates (Crucial if your app makes HTTPS/API requests)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port your application listens on (e.g., 8080)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/myapp"]