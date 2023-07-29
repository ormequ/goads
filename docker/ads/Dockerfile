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

ENV CONFIG_PATH ""
ENV GRPC_PORT 8000
EXPOSE $GRPC_PORT

RUN chmod +x ads

CMD ["./ads"]
