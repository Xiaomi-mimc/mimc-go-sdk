package packet

import (
	"mimc-go-sdk/protobuf/mimc"
)

type MIMCTimeoutPacket struct {
	timestamp int64
	packet    *mimc.MIMCPacket
}

func NewTimeoutPacket(timestamp int64, packet *mimc.MIMCPacket) *MIMCTimeoutPacket {
	return &MIMCTimeoutPacket{timestamp, packet}
}

func (this *MIMCTimeoutPacket) Timestamp() int64 {
	return this.timestamp
}

func (this *MIMCTimeoutPacket) Packet() *mimc.MIMCPacket {
	return this.packet
}
