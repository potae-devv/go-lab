# Start with golang base image
FROM golang:latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app with explicit GOOS and GOARCH for Linux
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main .

# Use PORT environment variable from Cloud Run
ENV PORT=8000

# Expose the port
EXPOSE ${PORT}

# Run the application
CMD ["./main"]