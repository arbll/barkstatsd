package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/arbll/barkstatsd/cmd/barkstatsd/command"
)

func main() {
	sock, err := net.Listen("tcp", "0.0.0.0:8123")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go func() {
		http.Serve(sock, nil)
	}()

	fmt.Println("Woof!")
	command.Bark()
}
