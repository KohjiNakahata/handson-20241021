package httpcmds

import (
	"encoding/json"
	"github.com/Diarkis/diarkis/server/http"
	"handson/lib/onlinestatus"
	"strings"
)

func exposeOnlineStatus() {
	http.Get("/onlinestatus/uids/:uids", getOnlineStatusList)
}
func getOnlineStatusList(resp *http.Response, req *http.Request, params *http.Params, next func(error)) {
	src, err := params.GetAsString("uids")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	uids := strings.Split(src, ",")
	list, err := onlinestatus.GetUserStatusList(uids)
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	ret := make(map[string]interface{}, 0)
	for _, userStatusData := range list {
		uid := userStatusData.UID
		inRoom := userStatusData.InRoom
		sessionData := make(map[uint8]string)
		for _, sd := range userStatusData.SessionData {
			sessionData[sd.Type] = sd.ID
		}
		v := make(map[string]interface{})
		v["InRoom"] = inRoom
		v["SessionData"] = sessionData
		ret[uid] = v
	}
	encoded, err := json.Marshal(ret)
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	resp.Respond(string(encoded), http.Ok)
	next(nil)
}
