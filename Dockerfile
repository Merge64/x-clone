# Use the official Golang image as the builder
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy only the go.mod and go.sum files first to cache the dependency download
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o x-clone ./server/main/main.go

# Use a minimal scratch image for the final stage
FROM scratch

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/x-clone .

# Command to run the application
CMD ["./x-clone"]