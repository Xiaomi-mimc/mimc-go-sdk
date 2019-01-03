package msg

type P2PMessage struct {
	packetId    *string
	sequence    *int64
	timestamp   *int64
	fromAccount *string
	toAccount   *string
	bizType     *string
	payload     []byte
}

func NewP2pMsg(packetId, fromAccount, toAccount *string, sequence, timestamp *int64, bizType *string, payload []byte) *P2PMessage {
	p2pMsg := new(P2PMessage)
	p2pMsg.packetId = packetId
	p2pMsg.sequence = sequence
	p2pMsg.timestamp = timestamp
	p2pMsg.fromAccount = fromAccount
	p2pMsg.toAccount = toAccount
	p2pMsg.payload = payload
	if bizType != nil {
		p2pMsg.bizType = bizType
	}
	return p2pMsg
}

func (this *P2PMessage) PacketId() *string {
	return this.packetId
}

func (this *P2PMessage) Sequence() *int64 {
	return this.sequence
}

func (this *P2PMessage) Timestamp() *int64 {
	return this.timestamp
}

func (this *P2PMessage) FromAccount() *string {
	return this.fromAccount
}

func (this *P2PMessage) ToAccount() *string {
	return this.toAccount
}

func (this *P2PMessage) Payload() []byte {
	return this.payload
}
