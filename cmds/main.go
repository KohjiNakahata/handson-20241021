package cmds

import (
	customcmds "handson/cmds/custom"
	dmcmds "handson/cmds/dm"
	httpcmds "handson/cmds/http"
	matchmakercmds "handson/cmds/matchmaker"
	roomcmds "handson/cmds/room"
	"handson/lib/onlinestatus"
)

func SetupUDP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	onlinestatus.Setup()
	customcmds.Expose()
}
func SetupTCP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	onlinestatus.Setup()
	customcmds.Expose()
}
func SetupHTTP() {
	httpcmds.Expose()
	onlinestatus.Setup()
}
