# rabbitmq-golang-client
A RabbitMQ producer and consumer written in Golang and deployed with Docker

## Run
`
 docker-compose up -d --build
`

## Test

Navigate to http://localhost:4444/hello and check the logs. You will see a message in the consumer's log: 
`
Received a message: hellos
`