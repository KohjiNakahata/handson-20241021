// © 2019-2024 Diarkis Inc. All rights reserved.

package handson

import (
	"fmt"

	pattack "handson/puffer/go/custom"

	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

type Handson struct {
	tcp *tcp.Client
	udp *udp.Client
}

func SetupHandsonAsTCP(c *tcp.Client) *Handson {
	h := &Handson{tcp: c}
	h.setup()
	return h
}

func SetupHandsonAsUDP(c *udp.Client) *Handson {
	h := &Handson{udp: c}
	h.setup()
	return h
}

func (h *Handson) setup() {
	if h.tcp != nil {
		h.tcp.OnResponse(h.onResponse)
		h.tcp.OnPush(h.onPush)
		return
	}
	if h.udp != nil {
		h.udp.OnResponse(h.onResponse)
		h.udp.OnPush(h.onPush)
	}
}

func (h *Handson) onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != pattack.AttackVer || cmd != pattack.AttackCmd {
		return
	}
	if status != uint8(1) {
		fmt.Printf("Attack failed: %v\n", string(payload))
		return
	}
	res := pattack.NewAttackResult()
	err := res.Unpack(payload)
	if err != nil {
		fmt.Printf("Failed to unpack attack response: %v\n", err)
		return
	}
	fmt.Printf("You dealt %d damage. Total damage: %d\n", res.Damage, res.TotalDamage)
}

func (h *Handson) onPush(ver uint8, cmd uint16, payload []byte) {
	// no push
	if ver != pattack.AttackVer || cmd != pattack.AttackCmd {
		return
	}
	res := pattack.NewAttackResult()
	err := res.Unpack(payload)
	if err != nil {
		fmt.Printf("Failed to unpack attack response: %v\n", err)
	}
	fmt.Printf("%s dealt %d damage. Total damage: %d\n", res.Uid, res.Damage, res.TotalDamage)
}

func (h *Handson) Attack(attackType uint8) {
	req := pattack.NewAttack()
	req.Type = attackType
	if h.tcp != nil {
		h.tcp.Send(pattack.AttackVer, pattack.AttackCmd, req.Pack())
		return
	}
	if h.udp != nil {
		h.udp.Send(pattack.AttackVer, pattack.AttackCmd, req.Pack())
	}
}
