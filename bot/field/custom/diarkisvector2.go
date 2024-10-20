// © 2019-2024 Diarkis Inc. All rights reserved.

// Code generated by Diarkis Puffer module: DO NOT EDIT.
//
// # Auto-generated by Diarkis Version 1.0.0
//
// - Maximum length of a string is 65535 bytes
// - Maximum length of a byte array is 65535 bytes
// - Maximum length of any array is 65535 elements
package custom

import "encoding/binary"
import "errors"
import "fmt"
import "math"
import "strings"
import "github.com/Diarkis/diarkis/util"

// DiarkisVector2Ver represents the ver of the protocol's command.
//
//	[NOTE] The value is optional and if ver is not given in the definition JSON, it will be 0.
const DiarkisVector2Ver uint8 = 0

// DiarkisVector2Cmd represents the command ID of the protocol's command ID.
//
//	[NOTE] The value is optional and if cmd is not given in the definition JSON, it will be 0.
const DiarkisVector2Cmd uint16 = 0

// DiarkisVector2 represents the command protocol data structure.
type DiarkisVector2 struct {
	// Command version of the protocol
	Ver uint8
	// Command ID of the protocol
	Cmd uint16
	X   float32
	Y   float32
}

// NewDiarkisVector2 creates a new instance of DiarkisVector2 struct.
func NewDiarkisVector2() *DiarkisVector2 {
	return &DiarkisVector2{Ver: 0, Cmd: 0, Y: 0, X: 0}
}

// Pack encodes DiarkisVector2 struct to a byte array to be delivered over the command.
func (proto *DiarkisVector2) Pack() []byte {
	bytes := make([]byte, 0)

	/* float32 */
	xBytes := make([]byte, 4)
	xBits := math.Float32bits(proto.X)
	xBytes[0] = byte(xBits >> 24)
	xBytes[1] = byte(xBits >> 16)
	xBytes[2] = byte(xBits >> 8)
	xBytes[3] = byte(xBits)
	xBytes = util.ReverseBytes(xBytes)
	bytes = append(bytes, xBytes...)

	/* float32 */
	yBytes := make([]byte, 4)
	yBits := math.Float32bits(proto.Y)
	yBytes[0] = byte(yBits >> 24)
	yBytes[1] = byte(yBits >> 16)
	yBytes[2] = byte(yBits >> 8)
	yBytes[3] = byte(yBits)
	yBytes = util.ReverseBytes(yBytes)
	bytes = append(bytes, yBytes...)

	// done
	return bytes
}

// Unpack decodes the command payload byte array to DiarkisVector2 struct.
func (proto *DiarkisVector2) Unpack(bytes []byte) error {
	if len(bytes) < 8 {
		return errors.New("DiarkisVector2UnpackError")
	}

	offset := 0

	/* float32 */
	xBytes := util.ReverseBytes(bytes[offset : offset+4])
	xBits := binary.BigEndian.Uint32(xBytes)
	offset += 4
	proto.X = math.Float32frombits(xBits)

	/* float32 */
	yBytes := util.ReverseBytes(bytes[offset : offset+4])
	yBits := binary.BigEndian.Uint32(yBytes)
	offset += 4
	proto.Y = math.Float32frombits(yBits)

	return nil
}

func (proto *DiarkisVector2) String() string {
	list := make([]string, 0)
	list = append(list, fmt.Sprint("X = ", proto.X))
	list = append(list, fmt.Sprint("Y = ", proto.Y))
	return strings.Join(list, " | ")
}

func (proto *DiarkisVector2) GetVer() uint8 {
	return proto.Ver
}
func (proto *DiarkisVector2) GetCmd() uint16 {
	return proto.Cmd
}
