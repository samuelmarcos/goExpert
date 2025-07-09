# Build stage
FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./cmd/ordersystem/main.go

# Run stage
FROM scratch
COPY --from=builder /app/main .
EXPOSE 8000
CMD ["./main"] 