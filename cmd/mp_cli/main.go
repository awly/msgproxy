package main

import (
	"fmt"
	"io"
	"os"

	"github.com/captaincronos/msgproxy"
)

var usage = `Usage: %s server_address conn_name`

func main() {
	if len(os.Args) != 3 {
		fmt.Printf(usage, os.Args[0])
		os.Exit(1)
	}
	c, err := msgproxy.Dial(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("failed connecting to server:", err)
	}
	defer c.Close()
	fmt.Println("connected to", os.Args[1], os.Args[2])
	go func() {
		_, err := io.Copy(c, os.Stdin)
		if err != nil {
			fmt.Println(err)
		}
		// interrupt the other io.Copy blocked on reading
		c.Close()
	}()
	_, err = io.Copy(os.Stdout, c)
	if err != nil {
		fmt.Println(err)
	}
}
