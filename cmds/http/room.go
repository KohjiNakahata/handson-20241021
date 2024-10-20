package httpcmds

import (
	"errors"
	"fmt"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/server/http"
	"handson/lib/meshCmds"
	"strings"
)

func exposeRoom() {
	http.Post("/room/create/:serverType/:maxMembers/:ttl/:interval", createRoom)
}
func createRoom(resp *http.Response, req *http.Request, params *http.Params, next func(error)) {
	serverType, err := params.GetAsString("serverType")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	serverType = strings.ToUpper(serverType)
	if serverType != server.UDPType && serverType != server.TCPType {
		resp.Respond(fmt.Sprintf("Invalid server type %s", serverType), http.Bad)
		next(errors.New(fmt.Sprintf("Invalid server type %s", serverType)))
		return
	}
	maxMembers, err := params.GetAsInt("maxMembers")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	ttl, err := params.GetAsInt64("ttl")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	interval, err := params.GetAsInt64("interval")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	meshCmds.CreateRemoteRoom(serverType, maxMembers, ttl, interval, func(err error, roomID string) {
		if err != nil {
			resp.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		resp.Respond(roomID, http.Ok)
		next(nil)
	})
}
