FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o semantic-search-api ./cmd/main.go

FROM scratch

COPY --from=builder /app/semantic-search-api /semantic-search-api

ENTRYPOINT ["/semantic-search"]