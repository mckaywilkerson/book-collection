# Start official go image
FROM golang:1.24.2

# Create app directory to hold source code
WORKDIR /app

# Copies go.mod and go.sum into /app
COPY go.mod go.sum ./

# install go dependencies
RUN go mod download

# Copy rest of file
COPY . .

# Builds app with optional configuration
RUN go build -o app ./cmd/bookapp

# Tell docker which port your container listens on
EXPOSE 8081

# Run binary
CMD ["./app"]