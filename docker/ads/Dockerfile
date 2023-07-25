FROM golang:1.20-alpine AS builder
WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o ads cmd/ads/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/ads /build/ads
COPY --from=builder /build/config/ads.env /build/config/config.env

ENV CONFIG_PATH /config/config.env

EXPOSE 18080
EXPOSE 18081

RUN chmod +x ads

CMD ["./ads"]
