FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o /app/main ./cmd/

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

COPY .env .env

COPY --from=builder /app/main .

EXPOSE 10800

CMD ["./main"]
