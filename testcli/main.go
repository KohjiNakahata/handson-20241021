package main

import (
	"bufio"
	"fmt"
	"handson/testcli/handson"
	"handson/testcli/resonance"
	"os"
	"strconv"
	"strings"

	"github.com/Diarkis/diarkis/client/go/test/cli"
)

var (
	tcpResonance *resonance.Resonance
	udpResonance *resonance.Resonance
	tcpHandson   *handson.Handson
	udpHandson   *handson.Handson
)

func main() {
	cli.SetupBuiltInCommands()
	cli.RegisterCommands("test", []cli.Command{{CmdName: "resonate", Desc: "Resonate your message", CmdFunc: resonate}})
	cli.RegisterCommands("handson", []cli.Command{
		{CmdName: "attack", Desc: "Attack for hands-on", CmdFunc: attack},
	})
	cli.Connect()
	if cli.TCPClient != nil {
		tcpResonance = resonance.SetupAsTCP(cli.TCPClient)
		tcpHandson = handson.SetupHandsonAsTCP(cli.TCPClient)
	}
	if cli.UDPClient != nil {
		udpResonance = resonance.SetupAsUDP(cli.UDPClient)
		udpHandson = handson.SetupHandsonAsUDP(cli.UDPClient)
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

func attack() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Which client to join a room? [tcp/udp]")
	client, _ := reader.ReadString('\n')
	fmt.Println("Enter the attack type. [1: melee, 2: range]")
	attackTypeStr, _ := reader.ReadString('\n')
	attackTypeStr = strings.Trim(attackTypeStr, "\n")
	attackType, err := strconv.Atoi(attackTypeStr)
	if err != nil || attackType != 1 && attackType != 2 {
		fmt.Println("Invalid attack type.")
		return
	}

	switch client {
	case "tcp\n":
		if tcpHandson == nil {
			return
		}
		tcpHandson.Attack(uint8(attackType))
	case "udp\n":
		if udpHandson == nil {
			return
		}
		udpHandson.Attack(uint8(attackType))
	}
}
