package scenarios

import (
	"encoding/json"
	"fmt"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/util"
	bot_client "handson/bot/scenario/lib/client"
	"handson/bot/scenario/lib/report"
	"handson/bot/scenario/packets"
	"strconv"
	"time"
)

type TicketScenarioParams struct {
	ServerTypeMM   string `json:"serverTypeMM"`
	ServerTypeTurn string `json:"serverTypeTurn"`
	UID            string `json:"userID"`
	TicketType     uint8  `json:"ticketType"`
	BattleDuration int    `json:"battleDuration"`
}
type TicketScenario struct {
	gp                *GlobalParams
	params            *TicketScenarioParams
	client            *bot_client.UDPClient
	trnClient         *bot_client.UDPClient
	metrics           *report.CustomMetrics
	createRoomReq     *packets.CreateRoomReq
	matchingStartedAt time.Time
	isOwner           bool
	roomID            string
}

var _ Scenario = &TicketScenario{}

func NewTicketScenario() Scenario {
	return &TicketScenario{}
}
func (s *TicketScenario) GetUserID() string {
	return s.params.UID
}
func (s *TicketScenario) ParseParam(index int, params []byte) error {
	var ticketParams *TicketScenarioParams
	err := json.Unmarshal(params, &ticketParams)
	if err != nil {
		logger.Erroru(strconv.Itoa(index), "Failed to unmarshal scenario params.", err.Error())
		return err
	}
	s.params = ticketParams
	logger.Debugu(s.GetUserID(), "Scenario Params. %#v", ticketParams)
	s.createRoomReq = &packets.CreateRoomReq{}
	json.Unmarshal(params, s.createRoomReq)
	logger.Debugu(s.GetUserID(), "Params for create Room. %#v", s.createRoomReq)
	return nil
}
func (s *TicketScenario) Run(gp *GlobalParams) error {
	logger.Infou(s.GetUserID(), "Starting scenario for user %v", s.GetUserID())
	s.gp = gp
	s.metrics = report.NewCustomMetrics()
	_, udpClient, err := bot_client.NewAndConnect(gp.Host, s.GetUserID(), s.params.ServerTypeMM, nil, gp.ReceiveByteSize, gp.UDPSendInterval)
	if err != nil {
		return err
	}
	s.client = udpClient
	udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicket, []uint8{bot_client.ResponseBad, bot_client.ResponseError}, s.handleTicketIssueError)
	udpClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdMMTicketComplete, s.handleTicketComplete)
	udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicketLeave, []uint8{bot_client.ResponseOk}, s.handleTicketLeave)
	udpClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdMMTicketBroadcast, s.handleTicketBroadcast)
	s.issueTicket()
	return nil
}
func (s *TicketScenario) OnIdle() {
	s.leaveTicket()
	return
}
func (s *TicketScenario) OnScenarioEnd() error {
	s.metrics.Stop()
	isActive := report.IsActive(s.GetUserID())
	if !isActive {
		kind, ver, cmd := s.client.GetLastActivity()
		logger.Warnu(s.GetUserID(), "I did not have any activities more than %d seconds, last command was ver: %d cmd: %d type: %s", report.Interval, ver, cmd, kind)
	}
	logger.Infou(s.GetUserID(), "disconnecting client...")
	s.client.Disconnect()
	if s.params.ServerTypeMM != s.params.ServerTypeTurn && s.trnClient != nil {
		s.trnClient.Disconnect()
	}
	logger.Noticeu(s.GetUserID(), "result per client   === \\")
	s.metrics.Print()
	logger.Noticeu(s.GetUserID(), "result per client   === /")
	return nil
}
func (s *TicketScenario) handleTicketIssueError(payload []byte) {
	time.Sleep(5 * time.Second)
	s.issueTicket()
}
func (s *TicketScenario) handleTicketComplete(payload []byte) {
	duration := time.Since(s.matchingStartedAt)
	s.metrics.Add("MATCHING_DURATION", "", duration.Seconds())
	res := packet.BytesToBytesList(payload)
	s.connectTurnServer()
	s.isOwner = string(res[0]) == s.GetUserID()
	if s.isOwner {
		logger.Sysu(s.GetUserID(), "Room Owner")
		s.createRoom()
	}
}
func (s *TicketScenario) handleTicketLeave(payload []byte) {
	s.regenerateParams()
	s.issueTicket()
	logger.Infou(s.GetUserID(), "Restarting ticket...")
}
func (s *TicketScenario) handleTicketBroadcast(payload []byte) {
	if !s.isOwner {
		s.roomID = string(payload)
		req := payload
		req = append(req, []byte("hello")...)
		s.trnClient.RSend(util.CmdBuiltInVer, util.CmdJoinRoom, req)
		logger.Sysu(s.GetUserID(), "Joining room... roomID: %s", s.roomID)
	}
}
func (s *TicketScenario) issueTicket() {
	s.matchingStartedAt = time.Now()
	s.metrics.Increment("ISSUE_TICKET", fmt.Sprintf("TYPE%d", s.params.TicketType))
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicket, []byte{s.params.TicketType})
}
func (s *TicketScenario) leaveTicket() {
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicketLeave, []byte{s.params.TicketType})
}
func (s *TicketScenario) connectTurnServer() {
	if s.params.ServerTypeMM == s.params.ServerTypeTurn {
		s.trnClient = s.client
	} else {
		_, trnClient, err := bot_client.NewAndConnect(s.gp.Host, s.GetUserID(), s.params.ServerTypeTurn, nil, s.gp.ReceiveByteSize, s.gp.UDPSendInterval)
		if err != nil {
			logger.Erroru(s.GetUserID(), "Failed to get Turn server")
			return
		}
		s.trnClient = trnClient
	}
	s.trnClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdCreateRoom, []uint8{bot_client.ResponseOk}, s.onCreateRoom)
	s.trnClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdJoinRoom, s.battle)
	if s.params.ServerTypeMM != s.params.ServerTypeTurn {
		s.trnClient.Connect()
	}
}
func (s *TicketScenario) createRoom() {
	bytes := packets.CreateCreateRoomReq(s.createRoomReq)
	s.trnClient.RSend(util.CmdBuiltInVer, util.CmdCreateRoom, bytes)
}
func (s *TicketScenario) leaveRoom() {
	s.trnClient.RSend(util.CmdBuiltInVer, util.CmdLeaveRoom, []byte(s.roomID))
}
func (s *TicketScenario) onCreateRoom(payload []byte) {
	s.roomID = string(payload[4:])
	req := []byte{s.params.TicketType}
	req = append(req, []byte(s.roomID)...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicketBroadcast, req)
}
func (s *TicketScenario) battle(payload []byte) {
	for i := 0; i < s.params.BattleDuration; i++ {
		s.trnClient.Send(util.CmdBuiltInVer, util.CmdBroadcastRoom, []byte("some battle command"))
		time.Sleep(time.Second)
	}
	if s.params.ServerTypeMM != s.params.ServerTypeTurn {
		s.trnClient.Disconnect()
		s.trnClient = nil
	}
	s.leaveRoom()
	s.leaveTicket()
}
func (s *TicketScenario) regenerateParams() {
	bytes, _ := s.gp.GenerateParams(0)
	var ticketParams *TicketScenarioParams
	json.Unmarshal(bytes, &ticketParams)
	s.params.TicketType = ticketParams.TicketType
}
