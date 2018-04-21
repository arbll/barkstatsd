package main

import (
	"fmt"

	"github.com/arbll/barkstatsd/cmd/barkstatsd/command"
)

func main() {
	fmt.Println("Woof!")
	command.Bark()
}
