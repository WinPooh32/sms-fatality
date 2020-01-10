FROM golang:1.13-alpine as builder
WORKDIR /go/src/
COPY . .
WORKDIR /go/src/acceptor
RUN GOOS=linux go build -o acceptor .

FROM alpine:3
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/acceptor/acceptor .
ENTRYPOINT ["./acceptor", "-c=./acceptor-config/acceptor.ini"]
