FROM golang:1.13-alpine as builder
# install git
RUN apk --no-cache add ca-certificates git less openssh
# install goose
RUN go get -v -u github.com/pressly/goose
WORKDIR /go/bin
RUN go build -tags='postgres no_mysql no_sqlite3 no_redshift no_tidb no_mssql' -i -o goose github.com/pressly/goose/cmd/goose


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /schema
WORKDIR /schema
COPY ./sql/messages/schema/ .
WORKDIR /schema
COPY --from=builder /go/bin/goose .

ENV START_SLEEP 15
ENV POSTGRES ""
ENTRYPOINT sleep $START_SLEEP && ./goose postgres "$POSTGRES" up
