package mimc

import (
	"container/list"
	"github.com/Xiaomi-mimc/mimc-go-sdk/message"
)

type Token interface {
	FetchToken() *string
}

type StatusDelegate interface {
	/**
	 * @param[status bool] 在线true；离线false
	 * @param[errType *string] 登录失败类型
	 * @param[errReason *string] 登录失败原因
	 * @param[errDec *string] 登录失败原因描述
	 */
	HandleChange(status bool, errType, errReason, errDec *string)
}

type MessageHandlerDelegate interface {
	HandleMessage(packets *list.List)
	HandleGroupMessage(packets *list.List)
	HandleServerAck(packetId *string, sequence, timestamp *int64, errMsg *string)
	HandleSendMessageTimeout(message *msg.P2PMessage)
	HandleSendGroupMessageTimeout(message *msg.P2TMessage)
}
