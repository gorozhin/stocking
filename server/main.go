package server

import(
	"net"
	//"fmt"
	"strconv"
	"../util"
	"../connection"
	"../authInterface"
)

type StockingServer struct{
	Host []byte
	Port uint16
	Listener *net.TCPListener
	Auth authInterface.AuthInterface
}

func (s *StockingServer)GetAddr() []byte {
	return s.Host
}

func (s *StockingServer)GetPort() uint16 {
	return s.Port
}

func (s *StockingServer) GetAuth() authInterface.AuthInterface{
	return s.Auth
}

func (s *StockingServer)Run(){
	service := util.DispIp4(s.Host)+":"+strconv.FormatInt(int64(s.Port), 10)
	
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	util.CheckError(err)

	if (s.Listener == nil){
		s.Listener, err = net.ListenTCP("tcp", tcpAddr)
	}
	util.CheckError(err)

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			continue
		}
		var sc connection.Connection = connection.Connection{Conn : conn, Server: s}
		go sc.HandleRequest()
	}
}
