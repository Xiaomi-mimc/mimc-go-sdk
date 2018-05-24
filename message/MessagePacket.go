package msg

import (
	"github.com/Xiaomi-mimc/mimc-go-sdk/packet"
)

type MsgPacket struct {
	msgType  string
	v6Packet *packet.MIMCV6Packet
}

func NewMsgPacket(msgType string, v6Packet *packet.MIMCV6Packet) *MsgPacket {
	return &MsgPacket{msgType, v6Packet}
}
func (this *MsgPacket) MsgType() string {
	return this.msgType
}

func (this *MsgPacket) Packet() *packet.MIMCV6Packet {
	return this.v6Packet
}
