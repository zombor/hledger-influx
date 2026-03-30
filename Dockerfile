FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o hledger-influx ./cmd/hledger-influx

FROM scratch
COPY --from=builder /app/hledger-influx /hledger-influx
ENTRYPOINT ["/hledger-influx"]
