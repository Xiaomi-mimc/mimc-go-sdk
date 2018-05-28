package packet

import (
	"bytes"
	"github.com/Xiaomi-mimc/mimc-go-sdk/cipher"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/ims"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/byte"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
	"github.com/golang/protobuf/proto"
)

var logger *log.Logger = log.GetLogger(log.InfoLevel)

type MIMCV6Packet struct {
	magic     uint16
	version   uint16
	packetLen int

	v6BodyBin []byte

	payloadType     uint16
	clientHeaderLen uint16
	payloadLen      uint32

	clientHeader *ims.ClientHeader
	payload      []byte
}

func NewV6Packet() *MIMCV6Packet {
	packet := new(MIMCV6Packet)
	return packet
}
func ParseBytesToPacket(headerBins, bodyBins, crcBins *[]byte, bodyKey *[]byte, secKey *string) *MIMCV6Packet {
	v6BinsBuffer := new(bytes.Buffer)
	v6BinsBuffer.Write(*headerBins)
	v6BinsBuffer.Write(*bodyBins)
	v6Bins := v6BinsBuffer.Bytes()

	crcfe := byteutil.GetIntFromBytes(crcBins, 0)
	crc := byteutil.Crc(v6Bins)
	if crcfe != crc {
		logger.Error("[ParseBytesToPacket] crc check fail.")
		return nil
	}
	v6Packet := NewV6Packet()
	v6Packet.magic = byteutil.GetUint16FromBytes(&v6Bins, cnst.V6_MAGIC_OFFSET)
	v6Packet.version = byteutil.GetUint16FromBytes(&v6Bins, cnst.V6_VERSION_OFFSET)
	if bodyKey != nil && len(*bodyKey) > 0 && bodyBins != nil && len(*bodyBins) > 0 {
		v6BodyBinsUnEn := cipher.Encrypt(*bodyKey, *bodyBins)
		bodyBins = &v6BodyBinsUnEn
	}
	if len(*bodyBins) == 0 {
		v6Packet.packetLen = 0
		return v6Packet
	} else {
		v6Packet.packetLen = len(*bodyBins)
	}
	v6Packet.payloadType = byteutil.GetUint16FromBytes(bodyBins, cnst.V6_PAYLOADTYPE_OFFSET)
	v6Packet.clientHeaderLen = byteutil.GetUint16FromBytes(bodyBins, cnst.V6_HEADERLEN_OFFSET)
	v6Packet.payloadLen = uint32(byteutil.GetIntFromBytes(bodyBins, cnst.V6_PAYLOADLEN_OFFSET))
	if v6Packet.payloadType != cnst.PAYLOAD_TYPE {
		return nil
	}
	headerBytes := byteutil.Copy(bodyBins, int(cnst.V6_BODY_HEADER_LENGTH), int(v6Packet.clientHeaderLen))
	payloadBytes := byteutil.Copy(bodyBins, int(cnst.V6_BODY_HEADER_LENGTH)+int(v6Packet.clientHeaderLen), int(v6Packet.payloadLen))

	clientHeader := new(ims.ClientHeader)
	err := proto.Unmarshal(headerBytes, clientHeader)

	if err != nil {
		logger.Error("[ParseBytesToPacket] deserial clientHeader fail: %v", err)
		return nil
	}

	if cnst.CMD_SECMSG == *(clientHeader.Cmd) {
		payloadKey := cipher.GenerateKeyForRC4(secKey, clientHeader.Id)
		payloadBytes = cipher.Encrypt(payloadKey, payloadBytes)
	}
	v6Packet.clientHeader = clientHeader
	v6Packet.payload = payloadBytes
	return v6Packet

}

func (this *MIMCV6Packet) PayloadType(payloadType uint16) {
	this.payloadType = payloadType
}

func (this *MIMCV6Packet) Payload(payload []byte) *MIMCV6Packet {
	this.payload = payload
	return this
}
func (this *MIMCV6Packet) GetPayload() []byte {
	return this.payload
}
func (this *MIMCV6Packet) ClientHeader(clientHeader *ims.ClientHeader) *MIMCV6Packet {
	this.clientHeader = clientHeader
	return this
}

func (this *MIMCV6Packet) Header(header *ims.ClientHeader) *MIMCV6Packet {
	this.clientHeader = header
	return this
}
func (this *MIMCV6Packet) GetHeader() *ims.ClientHeader {
	return this.clientHeader
}

func (this *MIMCV6Packet) Bytes(v6BodyKey []byte, payloadKey []byte) []byte {
	packetHead := make([]byte, cnst.V6_HEAD_LENGTH)
	var packet, body []byte

	if this.clientHeader == nil {
		initV6Head(&packetHead, 0)
		packet = packetHead
	} else {
		bodyHead := make([]byte, cnst.V6_BODY_HEADER_LENGTH)
		headBin, err := proto.Marshal(this.clientHeader)
		if err != nil {
			logger.Error("[bytes] marshaling error: %v", err)
			return nil
		}
		var headerBinLen, payloadLen int
		headerBinLen = len(headBin)
		if this.payload != nil {
			payloadLen = len(this.payload)
			if payloadLen != 0 && this.clientHeader.GetCmd() == cnst.CMD_SECMSG {
				this.payload = cipher.Encrypt(payloadKey, this.payload)
			}
			this.packetLen = int(cnst.V6_BODY_HEADER_LENGTH) + headerBinLen + payloadLen
			initBodyHead(&bodyHead, uint16(headerBinLen), payloadLen)
			body = byteutil.Integrate(bodyHead, headBin, this.payload)
		} else {
			payloadLen = 0
			this.packetLen = int(cnst.V6_BODY_HEADER_LENGTH) + headerBinLen + payloadLen
			initBodyHead(&bodyHead, uint16(headerBinLen), payloadLen)
			body = byteutil.Integrate(bodyHead, headBin)
		}
		if this.clientHeader.GetCmd() != cnst.CMD_CONN {
			body = cipher.Encrypt(v6BodyKey, body)
		}
		initV6Head(&packetHead, this.packetLen)
		packet = byteutil.Integrate(packetHead, body)
	}
	crc := byteutil.Crc(packet)
	buffer := new(bytes.Buffer)
	buffer.Write(packet)
	buffer.Write(byteutil.Bytes(crc))
	pktBytes := buffer.Bytes()
	return pktBytes
}

func (this *MIMCV6Packet) HeaderId() []byte {
	if this.clientHeader == nil {
		return nil
	}
	logger.Debug("v6header: %v", this.clientHeader)
	return []byte(*(this.clientHeader.Id))
}

func initV6Head(head *[]byte, packetLength int) {
	byteutil.TransferUint16(head, cnst.MAGIC, cnst.V6_MAGIC_OFFSET)
	byteutil.TransferUint16(head, cnst.V6_VERSION, cnst.V6_VERSION_OFFSET)
	byteutil.TransferInt(head, packetLength, cnst.V6_BODYLEN_OFFSET)
}
func initBodyHead(bodyHead *[]byte, headLen uint16, payloadLen int) {
	byteutil.TransferUint16(bodyHead, cnst.PAYLOAD_TYPE, cnst.V6_PAYLOADTYPE_OFFSET)
	byteutil.TransferUint16(bodyHead, headLen, cnst.V6_HEADERLEN_OFFSET)
	byteutil.TransferInt(bodyHead, payloadLen, cnst.V6_PAYLOADLEN_OFFSET)
}
