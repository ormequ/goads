FROM golang:1.20-alpine AS builder
WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/api /app/api

ENV CONFIG_PATH ""
ENV GRPC_PORT 80
EXPOSE $GRPC_PORT

RUN chmod +x api

CMD ["./api"]
