module github.com/Dpbm/jobsServer

go 1.23.5

require (
	github.com/Dpbm/shared v0.0.1
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.10.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)

replace github.com/Dpbm/shared v0.0.1 => ../shared
