package serverInterface

import (
	"../authInterface"
	"../middlewareInterface"
)

type ServerInterface interface {
	GetAddr() []byte
	GetPort() uint16
	GetAuth() authInterface.AuthInterface
	GetMiddleware() []middlewareInterface.MiddlewareInterface
}
