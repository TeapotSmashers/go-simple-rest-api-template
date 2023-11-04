# Use the official Go image as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files into the container
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Build the Go application
RUN go build -o my-go-api ./cmd/todos-backend

# Expose the port the application will run on
EXPOSE 3000

# Run the application
CMD ["./my-go-api"]
