package customcmds

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"handson/puffer/go/custom"
)

const CustomVer = 2
const helloCmdID = 10
const pushCmdID = 11
const resonanceCmdID = 13
const clientErrLog = 12
const matchmakerAdd = 100
const matchmakerRm = 101
const matchmakerSearch = 102
const matchmakerComplete = 103
const MatchedMemberLeaveCmdID = 1011
const p2pReportAddr = 110
const p2pInit = 111
const getUserStatusListCmdID = 500
const mmAddInterval = 40

var logger = log.New("CUSTOM")

func Expose() {
	diarkisexec.SetServerCommandHandler(CustomVer, helloCmdID, helloCmd)
	diarkisexec.SetServerCommandHandler(CustomVer, pushCmdID, pushCmd)
	diarkisexec.SetServerCommandHandler(custom.EchoVer, custom.EchoCmd, echoPufferCmd)
	diarkisexec.SetServerCommandHandler(CustomVer, matchmakerAdd, addToMatchMaker)
	diarkisexec.SetServerCommandHandler(CustomVer, matchmakerSearch, searchMatchMaker)
	diarkisexec.SetServerCommandHandler(CustomVer, p2pReportAddr, reportP2PAddr)
	diarkisexec.SetServerCommandHandler(CustomVer, p2pInit, initP2P)
	diarkisexec.SetServerCommandHandler(custom.GetFieldInfoVer, custom.GetFieldInfoCmd, getFieldInfo)
	diarkisexec.SetServerCommandHandler(CustomVer, getUserStatusListCmdID, getUserStatusList)
	diarkisexec.SetServerCommandHandler(CustomVer, resonanceCmdID, resonanceCmd)
}
func helloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	reliable := true
	userData.ServerRespond(payload, ver, cmd, server.Ok, reliable)
	next(nil)
}
func pushCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Push command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	reliable := true
	userData.ServerPush(ver, cmd, payload, reliable)
	next(nil)
}
func echoPufferCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello puffer command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	echoData := custom.NewEcho()
	err := echoData.Unpack(payload)
	if err != nil {
		logger.Error("Failed to unpack echo data: %v", err)
		userData.ServerRespond(nil, ver, cmd, server.Err, true)
		next(nil)
		return
	}
	logger.Debug("Unpacked echo data: %#v", echoData)
	userData.ServerRespond(echoData.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}
