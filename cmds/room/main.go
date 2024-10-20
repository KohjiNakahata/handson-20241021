package roomcmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomsupport"
	"github.com/Diarkis/diarkis/user"
	"handson/lib/meshCmds"
)

const ver uint8 = 3
const onRoomOwnerChangeCmd uint16 = 100

var logger = log.New("room")

func Setup() {
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
	room.SetOnRoomOwnerChange(onRoomOwnerChange)
	room.AfterCreateRoomCmd(afterCreateRoom)
	roomsupport.AfterRandomRoomCmd(afterRandomJoin)
	meshCmds.Setup()
}
func onDiscardCustomMessage(roomID string, userID string, sid string) []byte {
	logger.Debug("OnDiscardCustomMessage roomID:%v userID:%v sid:%v", roomID, userID, sid)
	return []byte(userID)
}
func onRoomOwnerChange(roomID string, ownerID string) {
	syncRoomOwnerID(roomID, ownerID)
}
func afterCreateRoom(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		next(nil)
		return
	}
	setupOnJoinCallback(roomID)
	ownerID := room.GetRoomOwnerID(roomID)
	syncRoomOwnerID(roomID, ownerID)
}
func afterRandomJoin(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		next(nil)
		return
	}
	if payload[0] != 0x00 {
		next(nil)
		return
	}
	setupOnJoinCallback(roomID)
}
func setupOnJoinCallback(roomID string) {
	room.SetOnJoinCompleteByID(roomID, func(rid string, ud *user.User) {
		ownerID := room.GetRoomOwnerID(roomID)
		if ownerID == "" {
			return
		}
		syncRoomOwnerID(roomID, ownerID)
	})
}
func syncRoomOwnerID(roomID string, ownerID string) {
	logger.Debug("OnRoomOwnerChange => broadcast the new owner ID %s to room %s", ownerID, roomID)
	message := []byte(ownerID)
	room.Announce(roomID, room.GetMemberIDs(roomID), ver, onRoomOwnerChangeCmd, message, true)
}
