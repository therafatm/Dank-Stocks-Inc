# seng468
seng 468 Spring 2018

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
