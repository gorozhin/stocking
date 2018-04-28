package main

import (
	"./server"
	"./authRedis"
	"./middlewareInterface"
	"./statusMiddleware"
	
	"github.com/go-redis/redis"
)

func main() {

	/*
auth.Container(map[string]string{
			"john" : "doe",
			"login" : "password",
		})}
*/


	status := &statusMiddleware.Status{}

	
	server := server.StockingServer{
		[]byte{188,166,180,165},
		7777,
		nil,
		authRedis.Container{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		})},
		[](middlewareInterface.MiddlewareInterface){status}}
	server.Run()
}


