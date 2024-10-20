package customcmds

import (
	"github.com/Diarkis/diarkis/field"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	custom "handson/puffer/go/custom"
)

func getFieldInfo(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Sys("Get Field Info Received from user: {}", userData.ID)
	fieldInfo := custom.NewGetFieldInfo()
	fieldInfo.NodeCount = int32(field.GetNodeNum())
	fieldInfo.FieldSize = field.GetFieldSize()
	fieldInfo.FieldOfVisionSize = field.GetFieldOfVisionSize()
	userData.ServerRespond(fieldInfo.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}
