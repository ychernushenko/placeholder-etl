# Use the latest Go version required by your project
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o placeholder-etl ./cmd/main.go

# Run the application
CMD ["./placeholder-etl"]
