# Stage 1: Build the Go application
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /app/aggregator ./cmd/aggregator

# Stage 2: Create a minimal runtime image
FROM golang:1.19

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/aggregator /app/aggregator

# Make the binary executable
RUN chmod +x ./aggregator

# Expose the port that the application listens on (if applicable)
EXPOSE 8012

# Command to run the application
CMD ["./aggregator/aggregator"]
