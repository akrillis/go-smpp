package smpp

import (
	"fmt"
	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/types"
)

func NewDeliverSMResp(dsr *types.DeliverSMRespInternalFormat) (pdu.Body, error) {
	message := pdu.NewDeliverSMResp()
	message.Header().Len = dsr.Len
	message.Header().Status = pdu.Status(dsr.Status)
	message.Header().Seq = dsr.Seq

	if err := message.Fields().Set(pdufield.MessageID, dsr.MessageID); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSMResp", "MessageID", err)
	}

	return message, nil
}

func (t *Transmitter) SendDeliverSMResp(p pdu.Body) error {
	_, err := t.do(p)
	return err
}
