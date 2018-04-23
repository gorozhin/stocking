package main

import (
	"net"
	"strconv"

	
	"./util"
	"./stockingConnection"
)

type StockingServer struct{
	host string
	port int64
	listener *net.TCPListener
}

func (s *StockingServer)Run(){
	service := s.host+":"+strconv.FormatInt(s.port, 10)
	
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	util.CheckError(err)

	if (s.listener == nil){
		s.listener, err = net.ListenTCP("tcp", tcpAddr)
	}
	util.CheckError(err)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		var sc stockingConnection.StockingConnection = stockingConnection.StockingConnection{Conn : conn}
		go sc.HandleRequest()
	}
}

func main() {
	server := StockingServer{"0.0.0.0", 7777, nil}
	server.Run()
}


