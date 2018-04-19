# Dank Stocks Inc
Dank Stocks Inc is a toy day trading system, built as part of a higher year distributed systems class at school. It's essentially a REST API, written in `Golang` accompanied by `postgres`,`redis`, `rabbitmq`, and `docker`. The picture below was the final architecture for the system by the end of the term. It's not perfect, but we scaled it to about ~15k transactions per second, for a 1.2 million transaction test file.

![Architecture](https://github.com/therafatm/Dank-Stocks-Inc/raw/master/arch-actual.jpg "Architecture!")

The core business logic lives of transaction server lives [here](https://github.com/therafatm/transaction-service). 

## Things we could add
* Add local user cache
* Shard user database
* Convert REST structure to a completely asynchronous system. Expose an API, which issues requests as messages on rabbit, which gets picked up by a worker. Use websockets to send data back.

## Running Dev
* `git submodules init` 
* `git submodules update` 
* `cd ./src`
* `docker-compose build`
* `docker-compose up` 

## Running Tests
* `cd ./src`
*  docker-compose exec test go test .
*  docker-compose exec test go run workload_generator/generator.go