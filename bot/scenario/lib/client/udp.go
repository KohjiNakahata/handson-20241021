package bot_client

import (
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/util"
	"handson/bot/scenario/lib/report"
	"slices"
)

type UDPClient struct {
	*TransportClient
	client       *udp.Client
	lastActivity struct {
		kind string
		ver  uint8
		cmd  uint16
	}
}
type HandlerType int

const (
	HandlerOnPush HandlerType = iota
	HandlerOnResponse
)

type ResponseStatus uint8

const (
	ResponseOk    = server.Ok
	ResponseBad   = server.Bad
	ResponseError = server.Err
)

func NewUDPClient(userID string, endpoint string, credentials *Credentials, rcvMaxSize int, interval int64) (*UDPClient, error) {
	dUDPClient := udp.New(rcvMaxSize, interval)
	udpClient := &UDPClient{TransportClient: &TransportClient{userID: userID, credentials: credentials, endpoint: endpoint}, client: dUDPClient}
	dUDPClient.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		report.IncrementResponseMetrics(ver, cmd)
		udpClient.setLastActivity("Response", ver, cmd)
		if status == ResponseOk {
			logger.Sysu(userID, util.StrConcat("\x1b[38;5;27m", "Response received, ver: %d, cmd: %d, status: %v, payload: %s (0x%x)", "\x1b[0m"), ver, cmd, status, string(payload), payload)
		} else {
			logger.Erroru(userID, util.StrConcat("\x1b[38;5;27m", "Response received, ver: %d, cmd: %d, status: %v, payload: %s (0x%x)", "\x1b[0m"), ver, cmd, status, string(payload), payload)
		}
	})
	dUDPClient.OnPush(func(ver uint8, cmd uint16, payload []byte) {
		report.IncrementPushMetrics(ver, cmd)
		udpClient.setLastActivity("Push", ver, cmd)
		logger.Sysu(userID, util.StrConcat("\x1b[38;5;27m", "Push received,     ver: %d, cmd: %d, payload: %s (0x%x)", "\x1b[0m"), ver, cmd, string(payload), payload)
	})
	return udpClient, nil
}
func (c *UDPClient) Connect() {
	c.client.SetEncryptionKeys(c.credentials.SID, c.credentials.Key, c.credentials.Iv, c.credentials.MacKey)
	logger.Info("Connecting UDP... %s", c.endpoint)
	c.client.Connect(c.endpoint)
}
func (c *UDPClient) Disconnect() {
	c.client.Disconnect()
}
func (c *UDPClient) Send(ver uint8, cmd uint16, payload []byte) {
	logger.Sysu(c.userID, util.StrConcat("\x1b[38;5;219m", "Sending Command,   ver: %d, cmd: %d, payload: %s (0x%x)", "\x1b[0m"), ver, cmd, string(payload), payload)
	report.IncrementCallCommandMetrics(ver, cmd)
	report.TouchAsActiveUser(c.userID)
	c.client.Send(ver, cmd, payload)
	c.setLastActivity("Send", ver, cmd)
}
func (c *UDPClient) RSend(ver uint8, cmd uint16, payload []byte) {
	logger.Sysu(c.userID, util.StrConcat("\x1b[38;5;219m", "RSending Command,  ver: %d, cmd: %d, payload: %s (0x%x)", "\x1b[0m"), ver, cmd, string(payload), payload)
	report.IncrementCallCommandMetrics(ver, cmd)
	report.TouchAsActiveUser(c.userID)
	c.client.RSend(ver, cmd, payload)
	c.setLastActivity("Send", ver, cmd)
}
func (c *UDPClient) OnPush(callback func(uint8, uint16, []byte)) {
	c.client.OnPush(callback)
}
func (c *UDPClient) OnResponse(callback func(uint8, uint16, uint8, []byte)) {
	c.client.OnResponse(callback)
}
func (c *UDPClient) RegisterOnPush(ver uint8, cmd uint16, cb func([]byte)) {
	c.OnPush(func(ver_ uint8, cmd_ uint16, payload []byte) {
		if ver_ == ver && cmd_ == cmd {
			cb(payload)
		}
	})
}
func (c *UDPClient) RegisterOnResponse(ver uint8, cmd uint16, statuses []uint8, cb func([]byte)) {
	c.OnResponse(func(ver_ uint8, cmd_ uint16, status uint8, payload []byte) {
		if ver_ == ver && cmd_ == cmd && slices.Contains(statuses, status) {
			cb(payload)
		}
	})
}
func (c *UDPClient) RegisterHandler(handlerType HandlerType, ver uint8, cmd uint16, cb func([]byte)) {
	switch handlerType {
	case HandlerOnPush:
		c.OnPush(func(ver_ uint8, cmd_ uint16, payload []byte) {
			if ver_ == ver && cmd_ == cmd {
				cb(payload)
			}
		})
	case HandlerOnResponse:
		c.OnResponse(func(ver_ uint8, cmd_ uint16, status uint8, payload []byte) {
			if ver_ == ver && cmd_ == cmd {
				cb(payload)
			}
		})
	default:
	}
}
func (c *UDPClient) SendRequest(ver uint8, cmd uint16, payload []byte, cb func(uint8, []byte)) {
	onResponse := func(ver_ uint8, cmd_ uint16, status uint8, res []byte) {
		if ver_ == ver && cmd_ == cmd {
			cb(status, res)
		}
	}
	c.OnResponse(onResponse)
	c.RSend(ver, cmd, payload)
}
func (c *UDPClient) GetLowLevelClient() *udp.Client {
	return c.client
}
func (c *UDPClient) setLastActivity(kind string, ver uint8, cmd uint16) {
	c.lastActivity.ver = ver
	c.lastActivity.cmd = cmd
	c.lastActivity.kind = kind
}
func (c *UDPClient) GetLastActivity() (string, uint8, uint16) {
	return c.lastActivity.kind, c.lastActivity.ver, c.lastActivity.cmd
}
