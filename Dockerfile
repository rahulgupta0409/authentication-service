# Use the official Golang image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

RUN apt-get update
RUN apt install -y protobuf-compiler

# Copy the Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy your Go service source code into the container
COPY . .

RUN /app/generate-pb.sh


# Expose the gRPC service port
EXPOSE 3001

# Command to run the binary
CMD ["go", "run", "/app/main.go"]
