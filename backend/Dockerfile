FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o /out/terraroute-api ./cmd/api

FROM alpine:3.21

RUN addgroup -S terraroute && adduser -S terraroute -G terraroute

WORKDIR /app
COPY --from=builder /out/terraroute-api /app/terraroute-api

USER terraroute
EXPOSE 8080

ENTRYPOINT ["/app/terraroute-api"]
