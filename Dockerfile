FROM golang:1.14 AS builder

ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w"

FROM gcr.io/distroless/base

COPY --from=builder /go/src/app/orders-api /go/src/app/configs/prod.yml /

EXPOSE 8080

ENTRYPOINT ["/orders-api"]