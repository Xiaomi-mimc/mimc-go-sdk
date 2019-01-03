package handler

import (
	"container/list"
	"github.com/Xiaomi-mimc/mimc-go-sdk/message"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
	"time"
)

type MsgHandler struct {
	appAccount string
}

var logger *log.Logger

func NewMsgHandler(appAccount string) *MsgHandler {
	logger = log.GetLogger()
	return &MsgHandler{appAccount}
}

func (this MsgHandler) HandleMessage(packets *list.List) {
	for ele := packets.Front(); ele != nil; ele = ele.Next() {
		p2pmsg := ele.Value.(*msg.P2PMessage)
		logger.Info("[%v] [handle p2p msg]%v -> %v: %v, pcktId: %v, seqId: %v, timestamp: %v.", this.appAccount, *(p2pmsg.FromAccount()), *(p2pmsg.ToAccount()), string(p2pmsg.Payload()), *(p2pmsg.PacketId()), *(p2pmsg.Sequence()), *(p2pmsg.Timestamp()))
	}

}
func (this MsgHandler) HandleGroupMessage(packets *list.List) {
	for ele := packets.Front(); ele != nil; ele = ele.Next() {
		p2tmsg := ele.Value.(*msg.P2TMessage)
		logger.Info("[%v] [handle p2t msg]%v  -> %v: %v, pcktId: %v, timestamp: %v.", this.appAccount, *(p2tmsg.FromAccount()), *(p2tmsg.GroupId()), string(p2tmsg.Payload()), *(p2tmsg.PacketId()), *(p2tmsg.Timestamp()))
	}
}

func (this MsgHandler) HandleServerAck(packetId *string, sequence, timestamp *int64, errMsg *string) {
	logger.Info("[%v] [handle server ack] packetId:%v, seqId: %v, timestamp:%v.", this.appAccount, *packetId, *sequence, *timestamp)
}
func (this MsgHandler) HandleSendMessageTimeout(message *msg.P2PMessage) {
	logger.Info("[%v] [handle p2pmsg timeout] packetId:%v, msg:%v, time: %v.", this.appAccount, *(message.PacketId()), string(message.Payload()), time.Now())
}
func (this MsgHandler) HandleSendGroupMessageTimeout(message *msg.P2TMessage) {
	logger.Info("[%v] [handle p2tmsg timeout] packetId:%v, msg:%v.", this.appAccount, *(message.PacketId()), string(message.Payload()))
}
