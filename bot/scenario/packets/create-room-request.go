// © 2019-2024 Diarkis Inc. All rights reserved.

// Code generated by "./generator ./config"; DO NOT EDIT.
package packets

import (
	"encoding/binary"
)

// Packet Format
//
//	 No.    | Name             | Type      | Size    |
//	--------|------------------|-----------|---------|
//	 0.     | maxMembers       | int       | 2 bytes |
//	 1.     | allowEmpty       | int       | 1 byte  |
//	 2.     | join             | int       | 1 byte  |
//	 3.     | ttl              | int       | 2 bytes |
//	 4.     | void             |           | 2 bytes |
//	 5.     | interval         | int       | 4 bytes |
//	--------|------------------|-----------|---------|
type CreateRoomReq struct {
	MaxMembers uint16 `json:"maxMembers"`
	AllowEmpty uint8  `json:"allowEmpty"`
	Join       uint8  `json:"join"`
	Ttl        uint16 `json:"ttl"`
	Void       uint16 `json:"void"`
	Interval   uint32 `json:"roomInterval"`
}

func CreateCreateRoomReq(values *CreateRoomReq) []byte {
	var bytes []byte

	// Append maxMembers
	bytes0 := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes0, uint64(values.MaxMembers))
	bytes = append(bytes, bytes0[6:]...)

	// Append allowEmpty
	bytes = append(bytes, values.AllowEmpty)

	// Append join
	bytes = append(bytes, values.Join)

	// Append ttl
	bytes3 := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes3, uint64(values.Ttl))
	bytes = append(bytes, bytes3[6:]...)

	// Append void
	bytes4 := make([]byte, 8)
	bytes = append(bytes, bytes4[6:]...)

	// Append interval
	bytes5 := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes5, uint64(values.Interval))
	bytes = append(bytes, bytes5[4:]...)

	return bytes
}
func ParseCreateRoomReq(payload []byte) *CreateRoomReq {
	var parsed CreateRoomReq

	// Parse maxMembers
	parsed.MaxMembers = binary.BigEndian.Uint16(payload[0:2])

	// Parse allowEmpty
	parsed.AllowEmpty = payload[2]

	// Parse join
	parsed.Join = payload[3]

	// Parse ttl
	parsed.Ttl = binary.BigEndian.Uint16(payload[4:6])

	// Parse interval
	parsed.Interval = binary.BigEndian.Uint32(payload[8:12])

	return &parsed
}
