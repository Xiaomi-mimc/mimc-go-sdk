package packet

type PacketBytes struct {
	HeaderBins *[]byte
	BodyBins   *[]byte
	CrcBins    *[]byte
	BodyKey    *[]byte
	SecKey     *string
}

func NewPacketBytes(HeaderBins, BodyBins, CrcBins, BodyKey *[]byte, Seckey *string) *PacketBytes {
	packetBytes := new(PacketBytes)
	packetBytes.HeaderBins = HeaderBins
	packetBytes.BodyBins = BodyBins
	packetBytes.CrcBins = CrcBins
	packetBytes.BodyKey = BodyKey
	packetBytes.SecKey = Seckey
	return packetBytes
}
