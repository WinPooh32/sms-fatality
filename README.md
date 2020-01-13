# sms-fatality
[![Go Report Card](https://goreportcard.com/badge/github.com/WinPooh32/sms-fatality)](https://goreportcard.com/report/github.com/WinPooh32/sms-fatality)

Durable message delivery service:

```
http POST SMS{phone, text string} --> Acceptor --> RabbitMQ --> Worker(N) --> PostgreSQL
                                   ^            ^                          ^
                                   |            |                          |
~ 7000 rps limit -------------------            |                          |
                                                |                          |
~ 350 rps limit ---------------------------------                          |
                                                                           |
~ 3500 rps limit -----------------------------------------------------------
```
\* All services run on Ubuntu 19.10 Linux 5.3.0-26, 4/8 core i5-u8250, 8gb ram, R3M120G8 ssd
<br>\*\* **RabbitMQ** is slow due to confirmation mode enabled on durable queue, message rate can be increased by multiply mq channels or multiply acceptor instances (e.g. using **nginx** balancing)

# Deploy
`docker-compose up --build`

# Stack
 * Docker
 * RabbitMQ
 * PostgreSQL
 * Golang: `net/http`, `streadway/amqp`, `lib/pq`

# Tools
 * [docker-compose](https://docs.docker.com/compose/) -  a tool for defining and running multi-container Docker applications
 * [goose](https://github.com/pressly/goose) - a database migration tool
 * [sqlc](https://github.com/kyleconroy/sqlc) - generates fully-type safe idiomatic Go code from SQL