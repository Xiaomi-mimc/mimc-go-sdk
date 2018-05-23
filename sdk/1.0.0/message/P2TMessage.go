package msg

type P2TMessage struct {
	packetId     *string
	sequence     *int64
	timestamp    *int64
	fromAccount  *string
	fromResource *string
	groupId      *int64
	payload      []byte
}

func NewP2tMsg(packetId, fromAccount, fromResource *string, sequence, timestamp, groupId *int64, payload []byte) *P2TMessage {
	p2tMsg := new(P2TMessage)
	p2tMsg.packetId = packetId
	p2tMsg.sequence = sequence
	p2tMsg.timestamp = timestamp
	p2tMsg.fromAccount = fromAccount
	p2tMsg.fromResource = fromResource
	p2tMsg.groupId = groupId
	p2tMsg.payload = payload
	return p2tMsg
}

func (this *P2TMessage) PacketId() *string {
	return this.packetId
}

func (this *P2TMessage) Sequence() *int64 {
	return this.sequence
}

func (this *P2TMessage) Timestamp() *int64 {
	return this.timestamp
}

func (this *P2TMessage) FromAccount() *string {
	return this.fromAccount
}

func (this *P2TMessage) FromResource() *string {
	return this.fromResource
}

func (this *P2TMessage) Payload() []byte {
	return this.payload
}
func (this *P2TMessage) GroupId() *int64 {
	return this.groupId
}
