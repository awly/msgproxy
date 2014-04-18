package msgproxy

import (
	"log"
	"net"
	"sync"
)

var waiting = struct {
	names map[string]net.Conn
	sync.Mutex
}{names: make(map[string]net.Conn)}

func ListenAndServe(laddr string) error {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}
	defer l.Close()
	return Serve(l)
}

func Serve(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go serve(c)
	}

	return nil
}

func serve(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1)
	_, err := c.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	buf = make([]byte, buf[0])
	_, err = c.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	name := string(buf)
	log.Println(name, "connected")

	waiting.Lock()
	defer waiting.Unlock()
	if other, ok := waiting.names[name]; ok {
		delete(waiting.names, name)
		join(c, other)
	} else {
		waiting.names[name] = c
	}
}

func join(a, b net.Conn) {}
