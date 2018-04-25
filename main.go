package main

import (
	"./server"
	"./auth"
)

func main() {
	
	server := server.StockingServer{
		7777,
		nil,
		auth.Container(map[string]string{
			"john" : "doe",
			"login" : "password",
		})}
	server.Run()
}


