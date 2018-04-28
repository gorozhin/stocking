package middlewareInterface

type MiddlewareInterface interface{
	ServerStarted()
	UnsupportedSocksVersion(byte)
	SuccessfullHandShake(string, byte)
	UnsuccessfullHandShake(string, error)
}


func ServerStarted(mw []MiddlewareInterface){
	for i := 0; i < len(mw); i++ {
		mw[i].ServerStarted()
	}
}

func UnsupportedSocksVersion(mw []MiddlewareInterface, ver byte){
	for i := 0; i < len(mw); i++ {
		mw[i].UnsupportedSocksVersion(ver)
	}
}

func SuccessfullHandShake(mw []MiddlewareInterface, remoteAddress string, auth_method byte){
	for i := 0; i < len(mw); i++ {
		mw[i].SuccessfullHandShake(remoteAddress, auth_method)
	}
}

func UnsuccessfullHandShake(mw []MiddlewareInterface, remoteAddress string, err error){
	for i := 0; i < len(mw); i++ {
		mw[i].UnsuccessfullHandShake(remoteAddress, err)
	}	
}
