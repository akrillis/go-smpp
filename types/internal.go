package types

import "time"

type SmsInternalWhereIsText uint8

const (
	SmsTextIntoText SmsInternalWhereIsText = iota
	SmsTextIntoTLVFieldsTagMessagePayload
)

type SmsInternalFormat struct {
	HeaderCommandLength  uint32
	HeaderSequenceNumber uint32
	Src                  string
	Dst                  string
	DstList              []string
	DLs                  []string
	Text                 []byte
	Validity             time.Duration
	Register             uint8
	TLVFields            map[uint16]interface{}
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	DestAddrTON          uint8
	DestAddrNPI          uint8
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	RegisteredDelivery   uint8
	ReplaceIfPresentFlag uint8
	SMDefaultMsgID       uint8
	NumberDests          uint8
	DataCoding           uint64
	WhereIsText          SmsInternalWhereIsText
	SMLength             uint8
}

type DeliverSMRespInternalFormat struct {
	Len       uint32
	ID        uint32
	Status    uint32
	Seq       uint32
	MessageID string
}
