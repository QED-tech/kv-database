package main

import (
	"bufio"
	"database/internal/network"
	"fmt"
	"github.com/docopt/docopt-go"
	"log"
	"os"
)

const version = "1.0.0"

func main() {
	usage := `
Usage:
  client [--address=<a>]
  client -h | --help
  client --version

Options:
  -h --help      Show this screen.
  -v --version   Show version.
  --address=<a>  Speed in knots [default: 127.0.0.1:8080].
`

	arguments, _ := docopt.ParseDoc(usage)

	showVersion, err := arguments.Bool("--version")
	if err != nil {
		log.Fatalf("failed to parse version: %v", err)
	}

	if showVersion {
		fmt.Printf("Current version: %s\n", version)

		return
	}

	address, err := arguments.String("--address")
	if err != nil {
		log.Fatalf("failed to parse address: %v", err)
	}

	c := network.NewTCPClient(address)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("[client] > ")

		in, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("[client] failed to read input: %v\n", err)
		}

		out, err := c.Send(in)
		if err != nil {
			log.Fatalf("[client] failed to send input: %v\n", err)
		}

		fmt.Println(out)
	}
}
