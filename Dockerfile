FROM golang:alpine AS builder
WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ads cmd/main/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/ads /build/ads
COPY /config /config

ENV CONFIG_PATH /config/config.env

EXPOSE 18080
EXPOSE 18081

CMD ["./ads"]
