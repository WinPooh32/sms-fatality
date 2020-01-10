module worker

go 1.13

require (
	common v0.0.0
	github.com/lib/pq v1.3.0
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
)

replace common => ../common
