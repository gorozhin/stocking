package stockingConnection

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"errors"
	"../util"
)

type StockingConnection struct{
	Conn net.Conn
}

func (sc *StockingConnection)HandShake() error{
	request := make([]byte, 258)
	
	read_len, err := sc.Conn.Read(request)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if request[0] != 5 {
		return errors.New("Unsupported SOCKS version")
	}

	auth_methods := request[2:read_len]

	err = errors.New("Unsupported auth method")
	for _, elt := range auth_methods {
		switch elt {
		case 0:
			sc.Conn.Write([]byte{5, 0})
			err = nil
			break;
			
		}
	}

	if err != nil {
		sc.Conn.Write([]byte{5, 0xff})
	}
	return err
	
}

func (sc *StockingConnection)DispatchRequest() (addr,  port string, err error){
	request := make([]byte, 258)
	
 	read_len, err := sc.Conn.Read(request)
	if err != nil {
		fmt.Println("Connection error: ",err)
		return "", "", errors.New("Connection crashed")
	}
	
	fmt.Println(request[:read_len])
	
	destAddr := strconv.FormatInt(int64(request[4]), 10)+"."+
		strconv.FormatInt(int64(request[5]), 10)+"."+
		strconv.FormatInt(int64(request[6]), 10)+"."+
		strconv.FormatInt(int64(request[7]), 10)
	destPort  := strconv.FormatInt(int64(request[8])*256+int64(request[9]), 10)
	
	fmt.Println("daddr: ", destAddr,":", destPort)	
	// sc.Conn.Write([]byte{5, 0, 0, 1, <IPHERE>, <IPHERE>, <IPHERE>, <IPHERE>, 0, 80})

	return destAddr,  destPort, nil
}

func (sc *StockingConnection)HandleRequest() {
	sc.Conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	defer sc.Conn.Close()

	
	err := sc.HandShake()	
	if err != nil {
		fmt.Println("Something went wrong 1")
		return
	}

	
	destAddr, destPort, err := sc.DispatchRequest()
	if err != nil {
		fmt.Println("Something went wrong 2")
		return
	}

	
	proxied, err := net.Dial("tcp", destAddr+":"+destPort)

	if err != nil {
		fmt.Println("Something went wrong 3")
		return
	}
	defer proxied.Close()
	
	
	go util.NetCopy(sc.Conn, proxied)
	util.NetCopy(proxied, sc.Conn)
}
