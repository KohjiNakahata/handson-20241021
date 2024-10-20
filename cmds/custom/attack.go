// © 2019-2024 Diarkis Inc. All rights reserved.

package customcmds

import (
	"errors"
	pattack "handson/puffer/go/custom"

	"github.com/Diarkis/diarkis/derror"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

func attack(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		err := errors.New("not in the room")
		userData.ServerRespond(derror.ErrData(err.Error(), derror.NotAllowed(0)), ver, cmd, server.Bad, true)
		next(err)
		return
	}

	req := pattack.NewAttack()
	req.Unpack(payload)

	damage := 0
	switch req.Type {
	case 1: // melee attack
		damage = util.RandomInt(1, 20)
	case 2: // range attack
		damage = util.RandomInt(1, 12) + 3
	default:
		err := errors.New("invalid attack type")
		userData.ServerRespond(derror.ErrData(err.Error(), derror.InvalidParameter(0)), ver, cmd, server.Bad, true)
		next(err)
		return
	}

	updatedDamage, updated := room.IncrProperty(roomID, "DAMAGE", int64(damage))
	if !updated {
		err := errors.New("incr property failed")
		userData.ServerRespond(derror.ErrData(err.Error(), derror.Internal(0)), ver, cmd, server.Err, true)
		next(err)
		return
	}
	logger.Info("Room %s has been attacked by %s using %s attack by %d damage. Total damage: %d", roomID, userData.ID, req.Type, damage, updatedDamage)
	res := pattack.NewAttackResult()
	res.Type = req.Type
	res.Uid = userData.ID
	res.Damage = uint16(damage)
	res.TotalDamage = uint16(updatedDamage)

	userData.ServerRespond(res.Pack(), ver, cmd, server.Ok, true)
	room.Relay(roomID, userData, ver, cmd, res.Pack(), true)
	next(nil)
}
