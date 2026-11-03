FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && \
    for cmd in cmd/*/; do \
        if [ -d "$cmd" ]; then \
            binary=$(basename "$cmd"); \
            echo "Building $binary..."; \
            CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "bin/$binary" "./$cmd"; \
        fi \
    done

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/* ./

RUN mkdir -p /root/uploads

EXPOSE 8080

CMD ["./api"]
