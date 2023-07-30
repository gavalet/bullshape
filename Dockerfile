FROM golang:1.18.1 AS builder
RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the webservice
RUN CGO_ENABLED=0 GOOS=linux go build -o myapi cmd/bullshape/main.go

FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage to the final stage
COPY --from=builder /app/myapi .
COPY --from=builder /app/*conf .
CMD ["./myapi"]