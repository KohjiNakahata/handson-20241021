package main

import (
	"bufio"
	"fmt"
	"github.com/Diarkis/diarkis/client/go/test/cli"
	"handson/testcli/resonance"
	"os"
)

var (
	tcpResonance *resonance.Resonance
	udpResonance *resonance.Resonance
)

func main() {
	cli.SetupBuiltInCommands()
	cli.RegisterCommands("test", []cli.Command{{CmdName: "resonate", Desc: "Resonate your message", CmdFunc: resonate}})
	cli.Connect()
	if cli.TCPClient != nil {
		tcpResonance = resonance.SetupAsTCP(cli.TCPClient)
	}
	if cli.UDPClient != nil {
		udpResonance = resonance.SetupAsUDP(cli.UDPClient)
	}
	cli.Run()
}
func resonate() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Which client to join a room? [tcp/udp]")
	client, _ := reader.ReadString('\n')
	fmt.Println("Enter the message you want to resonate.")
	message, _ := reader.ReadString('\n')
	switch client {
	case "tcp\n":
		if tcpResonance == nil {
			return
		}
		tcpResonance.Resonate(message)
	case "udp\n":
		if udpResonance == nil {
			return
		}
		udpResonance.Resonate(message)
	}
}
