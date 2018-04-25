package serverInterface

import (
	"../authInterface"
)

type ServerInterface interface {
	GetAddr() []byte
	GetPort() uint16
	GetAuth() authInterface.AuthInterface
}
