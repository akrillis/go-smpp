package smpp

import (
	"fmt"
	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/smpp/smpptest"
	"github.com/fiorix/go-smpp/v2/types"
	"golang.org/x/time/rate"
	"testing"
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
	transmitter := &Transmitter{
		Addr:        s.Addr(),
		User:        smpptest.DefaultUser,
		Passwd:      smpptest.DefaultPasswd,
		RateLimiter: rate.NewLimiter(rate.Limit(10), 1),
	}
	defer func() {
		if err := transmitter.Close(); err != nil {
			fmt.Printf("transmitter Close error: %v\n", err)
		}
	}()
	smppConnection := <-transmitter.Bind()
	switch smppConnection.Status() {
	case Connected:
	default:
		t.Fatal(smppConnection.Error())
	}

	smInternal := types.SmsInternalFormat{
		Src:           "79265033214",
		Dst:           "79265033277",
		Text:          []byte("Этот абонент пытался Вам позвонить"),
		SourceAddrTON: 1,
		SourceAddrNPI: 1,
		DestAddrTON:   1,
		DestAddrNPI:   1,
		PriorityFlag:  3,
		DataCoding:    0x08,
		WhereIsText:   types.SmsTextIntoText,
	}

	ndsm, err := NewDeliverSMMessage(&smInternal)
	if err != nil {
		t.Fatalf("NewDeliverSM error: %v", err)
	}

	fmt.Println(ndsm)

	answer, err := transmitter.SendDeliverSM(ndsm)
	if err != nil {
		t.Fatalf("SendDeliverSM error: %v", err)
	}

	msgid := answer.Fields()[pdufield.MessageID].String()
	if msgid != "foobar" {
		t.Fatalf("unexpected msgid: want foobar, have %q", msgid)
	}
}
