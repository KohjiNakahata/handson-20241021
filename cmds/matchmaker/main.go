package matchmakercmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
	customcmds "handson/cmds/custom"
)

const sampleTicketType0 uint8 = 0
const sampleTicketType1 uint8 = 1

var logger = log.New("matching")

func Setup() {
	matching.SetOnIssueTicket(sampleTicketType0, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{ProfileIDs: []string{"RankMatch"}, MaxMembers: 2, SearchInterval: 100, SearchTries: uint8(util.RandomInt(0, 300)), EmptySearches: 3, TicketDuration: 60, HowMany: 20, Tags: []string{""}, AddProperties: map[string]int{"rank": 1}, SearchProperties: map[string][]int{"rank": {1, 2, 3, 4, 5}}}
	})
	matching.SetOnTicketMatch(sampleTicketType0, func(t *matching.Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool {
		return false
	})
	matching.SetOnTicketComplete(sampleTicketType0, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType0, owner)
		list := make([]string, len(memberIDs)+1)
		list[0] = owner.ID
		index := 1
		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}
		return packet.StringListToBytes(list)
	})
	matching.SetOnTicketMemberLeaveAnnounce(sampleTicketType0, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		logger.Sys("Matched Member Leave Announce")
		return customcmds.CustomVer, customcmds.MatchedMemberLeaveCmdID, []byte(leftUser.ID)
	})
	matching.SetOnIssueTicket(sampleTicketType1, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{ProfileIDs: []string{"RankMatch20"}, MaxMembers: 4, SearchInterval: 100, SearchTries: uint8(util.RandomInt(0, 300)), EmptySearches: 3, TicketDuration: 60, HowMany: 20, Tags: []string{""}, AddProperties: map[string]int{"rank": 1}, SearchProperties: map[string][]int{"rank": {1, 2, 3, 4, 5}}}
	})
	matching.SetOnTicketMatch(sampleTicketType1, func(t *matching.Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool {
		return false
	})
	matching.SetOnTicketComplete(sampleTicketType1, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType1, owner)
		list := make([]string, len(memberIDs)+1)
		list[0] = owner.ID
		index := 1
		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}
		return packet.StringListToBytes(list)
	})
	matching.SetOnTicketMemberLeaveAnnounce(sampleTicketType1, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		logger.Sys("Ticket:1 Matched Member Leave Announce")
		return customcmds.CustomVer, customcmds.MatchedMemberLeaveCmdID, []byte(leftUser.ID)
	})
}
