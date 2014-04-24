package main

import (
	"flag"
	"log"

	"github.com/captaincronos/msgproxy"
)

var laddr = flag.String("a", ":8080", "local address used to accept incoming connections")

func main() {
	flag.Parse()
	if err := msgproxy.ListenAndServe(*laddr); err != nil {
		log.Println(err)
	}
}
