package handler

import (
	"container/list"
	"github.com/Xiaomi-mimc/mimc-go-sdk/message"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
	"time"
)

type MsgHandler struct {
}

var logger *log.Logger = log.GetLogger(log.InfoLevel)

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{}
}

func (this MsgHandler) HandleMessage(packets *list.List) {
	for ele := packets.Front(); ele != nil; ele = ele.Next() {
		p2pmsg := ele.Value.(*msg.P2PMessage)
		logger.Info("[handle p2p msg]%v -> %v, pcktId: %v, timestamp: %v.", *(p2pmsg.FromAccount()), string(p2pmsg.Payload()), *(p2pmsg.PacketId()), *(p2pmsg.Timestamp()))
	}

}
func (this MsgHandler) HandleGroupMessage(packets *list.List) {
	for ele := packets.Front(); ele != nil; ele = ele.Next() {
		p2tmsg := ele.Value.(*msg.P2TMessage)
		logger.Info("[handle p2t msg]%v -> %v, pcktId: %v, timestamp: %v.", *(p2tmsg.FromAccount()), string(p2tmsg.Payload()), *(p2tmsg.PacketId()), *(p2tmsg.Timestamp()))
	}
}
func (this MsgHandler) HandleServerAck(packetId *string, sequence, timestamp *int64) {
	logger.Info("[handle server ack] packetId:%v, timestamp:%v.", *packetId, *timestamp)
}
func (this MsgHandler) HandleSendMessageTimeout(message *msg.P2PMessage) {
	logger.Info("[handle p2pmsg timeout] packetId:%v, msg:%v, time: %v.", *(message.PacketId()), string(message.Payload()), time.Now())
}
func (this MsgHandler) HandleSendGroupMessageTimeout(message *msg.P2TMessage) {
	logger.Info("[handle p2tmsg timeout] packetId:%v, msg:%v.", *(message.PacketId()), string(message.Payload()))
}
