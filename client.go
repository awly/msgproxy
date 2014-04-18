package msgproxy

import "net"

func Dial(addr string, name string) (net.Conn, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	namemsg := append([]byte{byte(len(name))}, []byte(name)...)
	_, err = con.Write(namemsg)
	if err != nil {
		con.Close()
		return nil, err
	}
	return con, nil
}
