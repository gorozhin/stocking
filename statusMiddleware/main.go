package statusMiddleware

import (
	"fmt"
)

type Status struct {}


func (l *Status)ServerStarted(){
	fmt.Println("Server just started")
}

func (l *Status)UnsupportedSocksVersion(ver byte){
	fmt.Println("Unexpected SOCKS version:", ver)
}

func (l *Status)SuccessfullHandShake(remoteAddress string, auth_method byte){
	fmt.Println("Successfull handshake with: ", remoteAddress, "requested auth method:", auth_method)
}

func (l *Status)UnsuccessfullHandShake(remoteAddress string, err error){
	fmt.Println("Unsuccessfull handshake with: ", remoteAddress, "error: ", err.Error())

}
