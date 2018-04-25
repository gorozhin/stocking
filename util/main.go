package util

import (
	"fmt"
	"net"
	"os"
	"io"
	"strconv"
	"strings"
	
)

func NetCopy(input, output net.Conn) (err error) {
	buf := make([]byte, 8192)
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF && count > 0 {
				output.Write(buf[:count])
			}
			break
		}
		if count > 0 {
			output.Write(buf[:count])
		}
	}

	return 
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


func DispIp4(octets []byte) string {
	str := make([]string, len(octets))
	for i := 0; i < len(str); i++ {
		str[i] = strconv.FormatInt(int64(octets[i]), 10)
	}
	return strings.Join(str, ".")

}

func DispDomain(octets []byte) string {
	return string(octets[:len(octets)])
}

func DispIp6(octets []byte) string {

	str := make([]string, len(octets))
	for i := 0; i < len(str); i++ {
		pretender_str := strconv.FormatInt(int64(octets[i]), 16)
		pretender_str_len := len(pretender_str)

		add_zeros := (4 - pretender_str_len)
		zeros := strings.Repeat("0", add_zeros)

		joined_str := []string{zeros, pretender_str}

		str[i] = strings.Join(joined_str, "")
	}
	return strings.Join(str, ":")
}

func DispPort(portOcts []byte) string {
	return strconv.FormatInt(int64(portOcts[0])*256+int64(portOcts[1]), 10)
}
