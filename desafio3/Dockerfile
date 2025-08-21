# Build stage
FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./cmd/ordersystem

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/cmd/ordersystem/.env .
EXPOSE 8000
EXPOSE 50051
CMD ["./main"]