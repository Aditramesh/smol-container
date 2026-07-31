# Build stage
FROM golang:1.26.2 AS builder

WORKDIR /src 

COPY go.mod ./
RUN go mod download

RUN mkdir ubuntu-fs && \
    curl -L https://cdimage.ubuntu.com/ubuntu-base/releases/24.04/release/ubuntu-base-24.04.3-base-arm64.tar.gz \
    | tar -xz -C ubuntu-fs

COPY . ./

RUN mkdir ubuntu-fs/ROOT_FOR_CONTAINER

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o app .

# Runtime stage
FROM golang:1.26.2

COPY --from=builder /src/app /app
COPY --from=builder /src/ubuntu-fs /ubuntu-fs

ENTRYPOINT ["/app"]