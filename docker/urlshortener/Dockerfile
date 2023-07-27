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
COPY --from=builder /build/config/urlshortener.env /build/config/config.env

COPY cert cert

ENV CONFIG_PATH /config/config.env

EXPOSE 8888

RUN chmod +x urlshortener

CMD ["./urlshortener"]
