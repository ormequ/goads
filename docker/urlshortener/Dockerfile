FROM golang:1.20-alpine AS builder
WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o urlshortener cmd/urlshortener/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/urlshortener /build/urlshortener

ENV CONFIG_PATH ""
ENV GRPC_PORT 8000
EXPOSE $GRPC_PORT

RUN chmod +x urlshortener

CMD ["./urlshortener"]
