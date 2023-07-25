FROM golang:1.20-alpine AS builder
WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o auth cmd/auth/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/auth /build/auth
COPY --from=builder /build/config/auth.env /build/config/config.env

COPY cert cert

ENV CONFIG_PATH /config/config.env

EXPOSE 8080
EXPOSE 8081

RUN chmod +x auth

CMD ["./auth"]
