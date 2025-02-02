FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["go", "run", "main.go"]
