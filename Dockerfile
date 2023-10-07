# Build stage
FROM golang:1.21.2-alpine3.18 AS builder
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/main .
# COPY .env .

EXPOSE 8080
CMD [ "/app/main" ]
