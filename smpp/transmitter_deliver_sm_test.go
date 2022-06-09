package smpp

import (
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/fiorix/go-smpp/smpp/smpptest"
	"github.com/fiorix/go-smpp/v2/types"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestSendDeliverSM(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			r := pdu.NewDeliverSMResp()
			r.Header().Seq = p.Header().Seq
			r.Fields().Set(pdufield.MessageID, "foobar")
			c.Write(r)
		default:
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:        s.Addr(),
		User:        smpptest.DefaultUser,
		Passwd:      smpptest.DefaultPasswd,
		RateLimiter: rate.NewLimiter(rate.Limit(10), 1),
	}
	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}

	smInternal := types.SmsInternalFormat{
		Src:           "a",
		Dst:           "b",
		Text:          []byte("kokoko"),
		SourceAddrTON: 1,
		SourceAddrNPI: 1,
		DestAddrTON:   1,
		DestAddrNPI:   1,
		PriorityFlag:  3,
		DataCoding:    0x08,
		Validity:      time.Minute,
		WhereIsText:   types.SmsTextIntoText,
	}

	sm := &ShortMessage{
		Src:           "a",
		Dst:           "b",
		Text:          pdutext.UCS2("kokoko"),
		SourceAddrTON: 1,
		SourceAddrNPI: 1,
		DestAddrTON:   1,
		DestAddrNPI:   1,
		PriorityFlag:  3,
		Validity:      time.Minute,
	}

	ndsm, err := NewDeliverSM(&smInternal)
	if err != nil {
		t.Fatalf("NewDeliverSM error: %v", err)
	}

	answer, err := tx.SendDeliverSM(sm, ndsm)

	msgid := answer.RespID()
	if msgid == "" {
		t.Fatalf("pdu does not contain msgid: %#v", sm.Resp())
	}
	if msgid != "foobar" {
		t.Fatalf("unexpected msgid: want foobar, have %q", msgid)
	}
}
