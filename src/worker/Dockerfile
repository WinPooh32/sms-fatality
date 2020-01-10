FROM golang:1.13-alpine as builder
WORKDIR /go/src/
COPY . .
WORKDIR /go/src/worker
RUN GOOS=linux go build -o worker .

FROM alpine:3
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/worker/worker .
ENTRYPOINT ["./worker", "-c=./worker-config/worker.ini"]
