FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o /main ./cmd/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /main .

# Expose port 8080 for the container
EXPOSE 8080

CMD ["/app/main"]