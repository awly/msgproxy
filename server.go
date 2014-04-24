package msgproxy

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const timeout = time.Minute

var waiting = struct {
	names map[string]conn
	sync.Mutex
	sync.Once
}{names: make(map[string]conn)}

type conn struct {
	net.Conn
	t time.Time
}

// clean up waiting map periodically
func gcWaiting() {
	for {
		waiting.Lock()
		for k, v := range waiting.names {
			log.Println(v.RemoteAddr(), v.t)
			if time.Since(v.t) > timeout {
				log.Println(v.RemoteAddr(), "timeout")
				v.Close()
				delete(waiting.names, k)
			}
		}
		waiting.Unlock()
	}
}

func ListenAndServe(laddr string) error {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}
	defer l.Close()
	return Serve(l)
}

func Serve(l net.Listener) error {
	waiting.Do(func() {
		go gcWaiting()
	})

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
	log.Println(c.RemoteAddr(), "connected at", name)

	waiting.Lock()
	defer waiting.Unlock()
	if other, ok := waiting.names[name]; ok {
		delete(waiting.names, name)
		join(c, other)
	} else {
		waiting.names[name] = conn{Conn: c, t: time.Now()}
	}
}

func join(a, b net.Conn) {
	// buffer of 1 to prevent goroutine leak
	done := make(chan struct{}, 1)
	cp := func(from, to net.Conn) {
		_, err := io.Copy(to, from)
		if err != nil && !strings.Contains(err.Error(), "closed") {
			log.Println("io.Copy:", err)
		}
		done <- struct{}{}
	}
	go cp(a, b)
	go cp(b, a)
	<-done
	a.Close()
	b.Close()
	log.Println(a.RemoteAddr(), "<->", b.RemoteAddr(), "closed")
}
