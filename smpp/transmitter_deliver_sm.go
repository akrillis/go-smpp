package smpp

import (
	"fmt"
	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdutext"
	"github.com/fiorix/go-smpp/v2/types"
	"time"
)

func NewDeliverSMMessage(sif *types.SmsInternalFormat) (pdu.Body, error) {
	message := pdu.NewDeliverSM()
	mf := message.Fields()

	if err := mf.Set(pdufield.DataCoding, uint8(sif.DataCoding)); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "DataCoding", err)
	}

	if err := mf.Set(pdufield.SourceAddrNPI, sif.SourceAddrNPI); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "SourceAddrNPI", err)
	}

	if err := mf.Set(pdufield.SourceAddrTON, sif.SourceAddrTON); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "SourceAddrTON", err)
	}

	if err := mf.Set(pdufield.SourceAddr, sif.Src); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "SourceAddr", err)
	}

	if err := mf.Set(pdufield.DestAddrNPI, sif.DestAddrNPI); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "DestAddrNPI", err)
	}

	if err := mf.Set(pdufield.DestAddrTON, sif.DestAddrTON); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "DestAddrTON", err)
	}

	if err := mf.Set(pdufield.DestinationAddr, sif.Dst); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "DestinationAddr", err)
	}

	if err := mf.Set(pdufield.ServiceType, sif.ServiceType); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ServiceType", err)
	}

	if err := mf.Set(pdufield.ESMClass, sif.ESMClass); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ESMClass", err)
	}

	if err := mf.Set(pdufield.PriorityFlag, sif.PriorityFlag); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "PriorityFlag", err)
	}

	if err := mf.Set(pdufield.ProtocolID, sif.ProtocolID); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ProtocolID", err)
	}

	if err := mf.Set(pdufield.RegisteredDelivery, sif.RegisteredDelivery); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "RegisteredDelivery", err)
	}

	if err := mf.Set(pdufield.ReplaceIfPresentFlag, sif.ReplaceIfPresentFlag); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ReplaceIfPresentFlag", err)
	}

	if err := mf.Set(pdufield.ScheduleDeliveryTime, sif.ScheduleDeliveryTime); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ScheduleDeliveryTime", err)
	}

	if err := mf.Set(pdufield.SMDefaultMsgID, sif.SMDefaultMsgID); err != nil {
		return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "SMDefaultMsgID", err)
	}

	if sif.Validity != time.Duration(0) {
		if err := mf.Set(pdufield.ValidityPeriod, convertValidity(sif.Validity)); err != nil {
			return message, fmt.Errorf("[%s] %s set error: %w", "NewDeliverSM", "ValidityPeriod", err)
		}
	}

	if sif.WhereIsText == types.SmsTextIntoText {
		switch pdutext.DataCoding(sif.DataCoding) {
		case pdutext.Latin1Type:
			if err := mf.Set(pdufield.ShortMessage, pdutext.Latin1(sif.Text)); err != nil {
				return message, fmt.Errorf("[%s] %s with Latin1 set error: %w", "NewDeliverSM", "ShortMessage", err)
			}
		case pdutext.ISO88595Type:
			if err := mf.Set(pdufield.ShortMessage, pdutext.ISO88595(sif.Text)); err != nil {
				return message, fmt.Errorf("[%s] %s with ISO88595 set error: %w", "NewDeliverSM", "ShortMessage", err)
			}
		case pdutext.UCS2Type:
			if err := mf.Set(pdufield.ShortMessage, pdutext.UCS2(sif.Text)); err != nil {
				return message, fmt.Errorf("[%s] %s with UCS2 set error: %w", "NewDeliverSM", "ShortMessage", err)
			}
		default:
			if err := mf.Set(pdufield.ShortMessage, pdutext.Raw(sif.Text)); err != nil {
				return message, fmt.Errorf("[%s] %s with Raw set error: %w", "NewDeliverSM", "ShortMessage", err)
			}
		}
	}

	return message, nil
}

func (t *Transmitter) SendDeliverSM(p pdu.Body) (pdu.Body, error) {
	resp, err := t.do(p)
	if err != nil {
		return nil, err
	}
	if resp.PDU == nil {
		return nil, fmt.Errorf("unexpected empty PDU")
	}
	if id := resp.PDU.Header().ID; id != pdu.DeliverSMRespID {
		return resp.PDU, fmt.Errorf("unexpected PDU ID: %s", id)
	}
	if s := resp.PDU.Header().Status; s != 0 {
		return resp.PDU, s
	}

	return resp.PDU, resp.Err
}
