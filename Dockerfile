# Build stage
FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o app .

# Runtime stage
FROM golang:1.26.2

WORKDIR /app

COPY --from=builder /app/app .

ENTRYPOINT ["./app"]