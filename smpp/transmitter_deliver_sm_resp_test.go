package smpp

import (
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/smpptest"
	"github.com/fiorix/go-smpp/v2/types"
	"golang.org/x/time/rate"
	"testing"
)

func TestSendDeliverSMResp(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMRespID:
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

	dsrInternal := types.DeliverSMRespInternalFormat{
		Len:       6,
		Status:    0,
		Seq:       10,
		MessageID: "ololo",
	}

	ndsmr, err := NewDeliverSMResp(&dsrInternal)
	if err != nil {
		t.Fatalf("NewDeliverSMResp error: %v", err)
	}

	if err := tx.SendDeliverSMResp(ndsmr); err != nil {
		t.Fatalf("SendDeliverSMResp error: %v", err)
	}
}
