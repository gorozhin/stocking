package connection

import (
	"../serverInterface"
	"../util"
	"errors"
	"fmt"
	"net"
	"time"
	"../middlewareInterface"
)

const (
	socks5ProtocolVersion = 0x05
	reservedByte          = 0x00

	ip4AddressType    = 0x01
	domainAddressType = 0x03
	ip6AddressType    = 0x04

	successResponseCode               = 0x00
	internalServerErrorResponseCode   = 0x01
	notAllowedByRulesetResponseCode   = 0x02
	networkUnreachableResponseCode    = 0x03
	hostUnreachhableResponseCode      = 0x04
	connectionRefusedResponseCode     = 0x05
	expiredTTLResponseCode            = 0x06
	commandNotSupportedResponseCode   = 0x07
	unsupportedAdressTypeResponseCode = 0x08
)

type Connection struct {
	Conn   net.Conn
	Server serverInterface.ServerInterface
}

type connectionResponse struct {
	version      byte
	responseCode byte

	addressType byte
	address     []byte
	port        uint16
}

func (pc *connectionResponse) formate() []byte {
	res := []byte{pc.version, pc.responseCode, reservedByte, pc.addressType}
	res = append(res, pc.address...)
	res = append(res, byte(pc.port>>8))
	res = append(res, byte(pc.port&0xFF))

	return res
}

func (sc *Connection) HandShake() error {
	request := make([]byte, 258)
	read_len, err := sc.Conn.Read(request)

	if err != nil {

		fmt.Println(err)
		return err
	}
	
	if request[0] != socks5ProtocolVersion {
		middlewareInterface.UnsupportedSocksVersion(sc.Server.GetMiddleware(), request[0])
		return errors.New("Unsupported SOCKS version")
	}

	auth_methods := request[2:read_len]
	
	err = errors.New("Unsupported auth method")
	var auth_method byte = 0xff
	for _, elt := range auth_methods {
		switch elt {
		case 2:
			auth_method = 0x02
			err = nil
			break
	
		case 0:
			auth_method = 0x0
			err = nil
			break
		}
	}
	
	
	sc.Conn.Write([]byte{socks5ProtocolVersion, auth_method})
	if err != nil {
		middlewareInterface.UnsuccessfullHandShake(sc.Server.GetMiddleware(), sc.Conn.RemoteAddr().String(), err)

		return err
	}
	middlewareInterface.SuccessfullHandShake(sc.Server.GetMiddleware(), sc.Conn.RemoteAddr().String(), auth_method)

	if auth_method == 0x02{
		auth_request := make([]byte, 515)
		auth_request_read_len, err := sc.Conn.Read(auth_request)
		if auth_request_read_len == 1{}
		if err != nil {
			return errors.New("Auth connectivity error")
		}

		login_len := auth_request[1]
		login := string(auth_request[2:2+login_len])

		password_len := auth_request[2+login_len]
		password := string(auth_request[3+login_len:3+login_len+password_len])
		
		if sc.Server.GetAuth().Valid(login,password) {
			sc.Conn.Write([]byte{0x01, 0x00})
		} else {
			sc.Conn.Write([]byte{0x01, 0x1})
			return errors.New("Invalid password")
		}
		
		
	}
	
	return err
}

func (sc *Connection) DispatchRequest() (err error) {
	request := make([]byte, 263)

	read_len, err := sc.Conn.Read(request)
	if err != nil {
		//fmt.Println("Connection error: ", err)
		//fmt.Println(err.Error())
		return errors.New("Connection crashed")
	}


	destAddr, destPort, err := getAddressAndPort(request, read_len)

	proxied, err := net.Dial("tcp", destAddr+":"+destPort)

	if err != nil {
		cr := connectionResponse{socks5ProtocolVersion,
			internalServerErrorResponseCode,
			ip4AddressType,
			sc.Server.GetAddr(),
			sc.Server.GetPort()}

		sc.Conn.Write(cr.formate())

		//fmt.Println("Something went wrong 3")
		//fmt.Println(err.Error())

		return err
	}
	defer proxied.Close()


	cr := connectionResponse{socks5ProtocolVersion,
		successResponseCode,
		ip4AddressType,
		sc.Server.GetAddr(),
		sc.Server.GetPort()}
	sc.Conn.Write(cr.formate())

	go util.NetCopy(sc.Conn, proxied)
	util.NetCopy(proxied, sc.Conn)

	return nil
}


func (sc *Connection) HandleRequest() {
	sc.Conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	defer sc.Conn.Close()

	err := sc.HandShake()
	if err != nil {
		//fmt.Println("Something went wrong with handshake", err)
		return
	}

	err = sc.DispatchRequest()
	if err != nil {
		//fmt.Println("Something went wrong with dispatch", err)
		//fmt.Println(err.Error())
		return 
	}

}

func getAddressAndPort(request []byte, read_len int) (destAddr, destPort string, err error){
	err = nil

	switch request[3] {
	case ip4AddressType:
		destAddr = util.DispIp4(request[4:8])
	case domainAddressType:
		len := request[4]
		destAddr = util.DispDomain(request[5:5+len])
	case ip6AddressType:
		destAddr = util.DispIp6(request[4:20])
	default:
		return "", "", errors.New("Unknown adress type")
	}

	destPort = util.DispPort(request[read_len-2:])

	return destAddr, destPort, err
}

